// Package goflag provides a simple, type-safe command-line flag parsing library
// with support for arbitrarily nested subcommands, required-flag validation,
// per-flag validators, mutually exclusive flag groups, persistent (inherited)
// flags, and shell completion generation.
//
// # Quick start
//
//	cli := goflag.New("myapp", "Does something useful")
//	cli.String("config", "c", &cfg, "Path to config file").Required()
//	cli.Bool("verbose", "v", &verbose, "Verbose output")
//
//	subcmd, err := cli.Parse(os.Args)
//	if err != nil { log.Fatal(err) }
//	if subcmd != nil {
//	    if err := subcmd.Handler(nil); err != nil { log.Fatal(err) }
//	}
//
// Author: Dr. Abiira Nathan.
// Date: Sept 25. 2023
// License: MIT License
package goflag

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

//go:generate go tool stringer -type flagType

type flagType int

const (
	flagString flagType = iota
	flagInt
	flagInt64
	flagFloat32
	flagFloat64
	flagBool
	flagRune
	flagDuration
	flagStringSlice
	flagIntSlice
	flagTime
	flagIP
	flagMAC
	flagURL
	flagUUID
	flagHostPortPair
	flagEmail
	flagFilePath
	flagDirPath
)

// FlagValidator is a user-supplied validation function called after a flag's
// value has been parsed. It receives the concrete (dereferenced) value and
// must return (true, "") on success or (false, errorMessage) on failure.
type FlagValidator func(value any) (valid bool, errmsg string)

// Flag represents a single command-line flag with its type, names, default
// value pointer, usage string, and optional validators.
type Flag struct {
	flagType   flagType
	name       string          // Long flag name (--name).
	shortName  string          // Short flag name (-n); may be empty.
	value      any             // Pointer to the caller-supplied variable.
	usage      string          // Description shown in help output.
	required   bool            // Whether absence of this flag is an error.
	validators []FlagValidator // Run in order after parsing; first failure wins.
}

// Validate appends validators to the flag and returns the flag for chaining.
func (flag *Flag) Validate(validators ...FlagValidator) *Flag {
	flag.validators = append(flag.validators, validators...)
	return flag
}

// Required marks the flag as mandatory and returns it for chaining.
func (flag *Flag) Required() *Flag {
	flag.required = true
	return flag
}

// CLI is the root command-line context. It owns global flags, top-level
// subcommands, mutually exclusive flag groups, and the positional arguments
// that remain after parsing.
//
// Not safe for concurrent use.
type CLI struct {
	name            string     // Program name shown in help output.
	description     string     // Program description shown in help output.
	flags           []*Flag    // Global (root-level) flags.
	subcommands     []*SubCMD  // Top-level subcommands.
	args            []string   // Positional arguments captured at the root level.
	exclusiveGroups [][]string // Sets of mutually exclusive flag names.
}

// Args returns the positional arguments collected at the root level during
// the most recent Parse call.
func (c *CLI) Args() []string {
	return c.args
}

// ExclusiveFlags declares that at most one flag from names may be set in a
// single invocation. Returns c for chaining.
//
// Example:
//
//	cli.ExclusiveFlags("install", "uninstall")
func (c *CLI) ExclusiveFlags(names ...string) *CLI {
	if len(names) < 2 {
		panic("ExclusiveFlags requires at least two flag names")
	}
	group := make([]string, len(names))
	copy(group, names)
	c.exclusiveGroups = append(c.exclusiveGroups, group)
	return c
}

// completionCmd is the auto-registered completion subcommand. It is kept as a
// package-level variable so that Parse can skip required-flag checks for it.
var completionCmd *SubCMD

// New creates a new CLI with the given program name and description. A
// "completion" subcommand is automatically registered that can generate and
// install bash/zsh completion scripts.
func New(name, description string) *CLI {
	cli := &CLI{
		name:        name,
		description: description,
		flags: []*Flag{
			{name: "help", shortName: "h", flagType: flagBool, usage: "Print help message and exit"},
		},
	}

	// Completion subcommand — registered before the caller adds anything so
	// that it is always present. Variables are captured by closure.
	var shell string
	var install bool
	var uninstall bool

	completionCmd = cli.subCommand(name, description, "completion", "Generate shell completion scripts",
		func(userdata any) error {
			if install && uninstall {
				return fmt.Errorf("cannot use --install and --uninstall together")
			}

			if uninstall {
				if err := cli.UninstallCompletion(shell); err != nil {
					return fmt.Errorf("failed to uninstall completion: %w", err)
				}
				return nil
			}

			if install {
				if err := cli.InstallCompletion(shell); err != nil {
					return fmt.Errorf("failed to install completion: %w", err)
				}
				return nil
			}

			switch shell {
			case "bash":
				cli.GenBashCompletion(os.Stdout)
			case "zsh":
				cli.GenZshCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell: %s", shell)
			}
			return nil
		},
	)

	// Attach the completion subcommand's own flags via the low-level method
	// so chaining returns *SubCMD (Required/Validate apply to last flag).
	completionCmd.
		Flag(flagString, "shell", "s", &shell, "The shell to generate completions for [bash|zsh]").
		Required().Validate(Choices([]string{"zsh", "bash"})).
		Flag(flagBool, "install", "i", &install, "Install the completion script").
		Flag(flagBool, "uninstall", "u", &uninstall, "Uninstall the completion script")

	// The completion subcommand's own --install/--uninstall are mutually
	// exclusive; register that constraint on the subcommand itself.
	completionCmd.ExclusiveFlags("install", "uninstall")

	return cli
}

// subCommand is the internal constructor shared by CLI and SubCMD registration
// paths. It builds the child and appends it to cli.subcommands; for SubCMD
// nesting the caller sets child.parent afterward.
func (c *CLI) subCommand(
	_ /* programName */, _ /* programDesc */ string,
	name, description string,
	handler func(userdata any) error,
) *SubCMD {
	child := &SubCMD{
		name:        name,
		description: description,
		Handler:     handler,
		flags: []*Flag{
			{name: "help", shortName: "h", flagType: flagBool, usage: "Print help message and exit"},
		},
	}
	c.subcommands = append(c.subcommands, child)
	return child
}

// addFlag creates a Flag and appends it to the CLI's global flag list.
func (c *CLI) addFlag(ft flagType, name, shortName string, valuePtr any, usage string) *Flag {
	f := &Flag{
		flagType:  ft,
		name:      name,
		shortName: shortName,
		value:     valuePtr,
		usage:     usage,
	}
	validateFlag(f)
	c.flags = append(c.flags, f)
	return f
}

// SubCommand registers a top-level subcommand and returns it for flag
// configuration via method chaining.
//
// Panics if name or description is empty, or if handler is nil.
func (c *CLI) SubCommand(name, description string, handler func(userdata any) error) *SubCMD {
	if handler == nil {
		panic("subcommand cannot be registered with a nil handler")
	}
	if name == "" {
		panic("subcommand name cannot be empty")
	}
	if description == "" {
		panic("subcommand description cannot be empty")
	}
	return c.subCommand("", "", name, description, handler)
}

// Parse tokenises argv (pass os.Args), populates flag variables, and returns
// the deepest matching SubCMD in the subcommand tree. Positional arguments are
// stored on the matched SubCMD (or on the CLI itself when no subcommand
// matched). Returns nil, nil when no subcommand was found and no error
// occurred.
//
// The first element of argv is assumed to be the program name and is skipped.
func (c *CLI) Parse(argv []string) (*SubCMD, error) {
	if len(argv) >= 1 {
		argv = argv[1:] // drop program name
	}

	var positional []string
	processed := make(map[string]bool)

	subcmd, remaining, err := c.parseLevel(argv, c.flags, c.subcommands, processed, &positional)
	if err != nil {
		return nil, err
	}

	if subcmd == nil {
		// No subcommand matched — check global required flags and exclusive groups.
		for _, f := range c.flags {
			if f.required && !processed[f.name] {
				return nil, fmt.Errorf("missing required flag [-%s | --%s]", f.shortName, f.name)
			}
		}
		if err := checkExclusiveGroups(c.exclusiveGroups, processed); err != nil {
			return nil, err
		}
		c.args = positional
		return nil, nil
	}

	// Recurse into the matched subcommand tree with the remaining tokens.
	leaf, err := c.parseSubCommandTree(subcmd, remaining)
	if err != nil {
		return nil, err
	}
	return leaf, nil
}

// parseLevel scans tokens and dispatches them as global flags, subcommand
// names, or positional arguments. It returns the first subcommand token that
// matches, along with the tokens that follow it.
func (c *CLI) parseLevel(
	argv []string,
	flags []*Flag,
	subcommands []*SubCMD,
	processed map[string]bool,
	positional *[]string,
) (*SubCMD, []string, error) {

	for i := 0; i < len(argv); i++ {
		arg := argv[i]

		if strings.TrimSpace(arg) == "" {
			continue
		}

		// Expand --key=value into two tokens in-place.
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			arg = parts[0]
			// Insert the value token right after the current position.
			tail := make([]string, len(argv[i+1:]))
			copy(tail, argv[i+1:])
			argv = append(argv[:i+1], append([]string{parts[1]}, tail...)...)
		}

		if len(arg) > 1 && arg[0] == '-' {
			// Flag token.
			var name string
			if len(arg) > 2 && arg[1] == '-' {
				name = arg[2:] // --long
			} else {
				name = arg[1:] // -s
			}

			if isHelpFlag(name) {
				c.PrintUsage(os.Stdout)
				os.Exit(0)
			}

			f, err := parseFlags(&flags, name, i, argv)
			if err != nil {
				return nil, nil, err
			}
			if f != nil {
				processed[f.name] = true
			}
			continue
		}

		// Non-flag token: try to match a subcommand.
		for _, cmd := range subcommands {
			if cmd.name == arg {
				return cmd, argv[i+1:], nil
			}
		}

		// Not a subcommand — treat as positional argument.
		*positional = append(*positional, arg)
	}

	return nil, nil, nil
}

// parseSubCommandTree recursively parses flags and nested subcommands for cmd,
// returning the deepest matched subcommand (the "leaf").
func (c *CLI) parseSubCommandTree(cmd *SubCMD, argv []string) (*SubCMD, error) {
	processed := make(map[string]bool)
	var positional []string

	// Reset positional args from any previous parse.
	cmd.args = nil

	// Build the effective flag set: own flags + inherited persistent flags.
	effectiveFlags := buildEffectiveFlags(cmd)

	nested, remaining, err := c.parseSubLevel(cmd, effectiveFlags, argv, processed, &positional)
	if err != nil {
		return nil, err
	}

	// Skip required-flag validation for the built-in completion subcommand.
	if cmd != completionCmd {
		for _, f := range effectiveFlags {
			if f.required && !processed[f.name] {
				return nil, fmt.Errorf("missing required flag [-%s | --%s]", f.shortName, f.name)
			}
		}
		if err := checkExclusiveGroups(cmd.exclusiveGroups, processed); err != nil {
			return nil, err
		}
	}

	if nested == nil {
		// This is the leaf node.
		cmd.args = positional
		return cmd, nil
	}

	// Recurse into the nested subcommand.
	return c.parseSubCommandTree(nested, remaining)
}

// parseSubLevel is the per-SubCMD token scanner; mirrors parseLevel but
// operates on the provided effective flag set (own + inherited) and child
// subcommands.
func (c *CLI) parseSubLevel(
	cmd *SubCMD,
	effectiveFlags []*Flag,
	argv []string,
	processed map[string]bool,
	positional *[]string,
) (*SubCMD, []string, error) {

	for i := 0; i < len(argv); i++ {
		arg := argv[i]

		if strings.TrimSpace(arg) == "" {
			continue
		}

		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			arg = parts[0]
			tail := make([]string, len(argv[i+1:]))
			copy(tail, argv[i+1:])
			argv = append(argv[:i+1], append([]string{parts[1]}, tail...)...)
		}

		if len(arg) > 1 && arg[0] == '-' {
			var name string
			if len(arg) > 2 && arg[1] == '-' {
				name = arg[2:]
			} else {
				name = arg[1:]
			}

			if isHelpFlag(name) {
				cmd.PrintUsage(os.Stdout)
				os.Exit(0)
			}

			f, err := parseFlags(&effectiveFlags, name, i, argv)
			if err != nil {
				return nil, nil, err
			}

			if f != nil {
				processed[f.name] = true
				if f.flagType != flagBool {
					i++ // skip the value token already consumed by parseFlags
				}
			}
			continue
		}

		// Try to match a nested subcommand.
		for _, child := range cmd.subcommands {
			if child.name == arg {
				return child, argv[i+1:], nil
			}
		}

		// Positional argument.
		*positional = append(*positional, arg)
	}

	return nil, nil, nil
}

// ParseAndInvoke is a convenience wrapper around Parse that immediately calls
// the matched subcommand's Handler, passing userdata through to it. If a
// preInvokeCallback is provided it is called with the matched subcommand (which
// may be nil) and userdata before the handler runs.
//
// Returns any error from Parse or from the handler itself.
func (c *CLI) ParseAndInvoke(argv []string, userdata any, preInvokeCallback func(cmd *SubCMD, userdata any)) error {
	subcmd, err := c.Parse(argv)
	if err != nil {
		return err
	}

	if preInvokeCallback != nil {
		preInvokeCallback(subcmd, userdata)
	}

	if subcmd != nil && subcmd.Handler != nil {
		return subcmd.Handler(userdata)
	}
	return nil
}

// parseFlags locates the named flag in the provided slice, reads its value
// from argv[i+1] when required, runs validators, and returns the flag.
func parseFlags(flags *[]*Flag, name string, i int, argv []string) (*Flag, error) {
	f := findFlag(*flags, name)
	if f == nil {
		return nil, fmt.Errorf("unknown flag: %s", name)
	}

	valueIndex := i + 1

	// Boolean flags are true when present without an explicit value.
	if valueIndex >= len(argv) {
		if f.flagType == flagBool {
			*f.value.(*bool) = true
			return f, nil
		}
		return f, fmt.Errorf("missing value for flag [-%s | --%s]", f.shortName, f.name)
	}

	next := argv[valueIndex]
	if next == "" {
		return f, fmt.Errorf("empty value for flag [-%s | --%s]", f.shortName, f.name)
	}

	if next[0] == '-' {
		if f.flagType == flagBool {
			*f.value.(*bool) = true
			return f, nil
		}
		return f, fmt.Errorf("missing value for flag [-%s | --%s]", f.shortName, f.name)
	}

	if err := parseFlagValue(f, next); err != nil {
		return f, err
	}

	// Run validators against the dereferenced, typed value.
	for _, v := range f.validators {
		if v == nil {
			continue
		}
		concrete := reflect.ValueOf(f.value).Elem().Interface()
		if ok, msg := v(concrete); !ok {
			return f, fmt.Errorf("invalid value (%v) for flag [--%s]: %s", concrete, f.name, msg)
		}
	}

	return f, nil
}

// findFlag searches flags by long name or short name and returns the first match.
func findFlag(flags []*Flag, name string) *Flag {
	for _, f := range flags {
		if f.name == name || f.shortName == name {
			return f
		}
	}
	return nil
}

func isHelpFlag(name string) bool {
	return name == "help" || name == "h"
}

// buildEffectiveFlags returns the combined flag set for cmd: its own flags,
// its own persistent flags, followed by all persistent flags collected by
// walking up the parent chain.
// Flags with names already present in cmd.flags are skipped to avoid shadowing.
func buildEffectiveFlags(cmd *SubCMD) []*Flag {
	// Start with the subcommand's own flags.
	result := make([]*Flag, len(cmd.flags))
	copy(result, cmd.flags)

	// Track which names are already covered so ancestors cannot shadow them.
	seen := make(map[string]bool, len(cmd.flags))
	for _, f := range cmd.flags {
		seen[f.name] = true
	}

	// Walk up the chain starting from cmd itself to collect persistent flags.
	ancestor := cmd
	for ancestor != nil {
		for _, f := range ancestor.persistentFlags {
			if !seen[f.name] {
				result = append(result, f)
				seen[f.name] = true
			}
		}
		ancestor = ancestor.parent
	}

	return result
}

// checkExclusiveGroups returns an error if more than one flag within any
// exclusive group was set in processed.
func checkExclusiveGroups(groups [][]string, processed map[string]bool) error {
	for _, group := range groups {
		var set []string
		for _, name := range group {
			if processed[name] {
				set = append(set, "--"+name)
			}
		}
		if len(set) > 1 {
			return fmt.Errorf("flags %s are mutually exclusive; only one may be provided",
				strings.Join(set, " and "))
		}
	}
	return nil
}

// printFlag writes a single flag's help line to w, aligned using longestFlagName.
func printFlag(flag *Flag, w io.Writer, longestFlagName int, indent string) {
	fmt.Fprintf(w, "%s--%-*s ", indent, longestFlagName, flag.name)

	value := ""
	if v := reflect.ValueOf(flag.value); v.IsValid() {
		value = fmt.Sprintf("%v", v.Elem().Interface())
	}

	if flag.flagType == flagString {
		if flag.shortName != "" {
			fmt.Fprintf(w, "-%s: %s (default: %q)\n", flag.shortName, flag.usage, value)
		} else {
			fmt.Fprintf(w, "%s (default: %q)\n", flag.usage, value)
		}
	} else {
		if flag.shortName != "" {
			fmt.Fprintf(w, "-%s: %s (default: %v)\n", flag.shortName, flag.usage, value)
		} else {
			fmt.Fprintf(w, "%s (default: %v)\n", flag.usage, value)
		}
	}
}

// printSubCommand writes a subcommand's name, description, flags, inherited
// persistent flags, and nested subcommands to w.
func printSubCommand(cmd *SubCMD, w io.Writer) {
	fmt.Fprintf(w, "%s: %s\n", cmd.name, cmd.description)

	longest := 0
	for _, f := range cmd.flags {
		if f.name != "help" && len(f.name) > longest {
			longest = len(f.name)
		}
	}
	for _, f := range cmd.flags {
		if f.name == "help" {
			continue
		}
		printFlag(f, w, longest, "    ")
	}

	if len(cmd.persistentFlags) > 0 {
		fmt.Fprintf(w, "  Persistent flags:\n")
		longestP := 0
		for _, f := range cmd.persistentFlags {
			if len(f.name) > longestP {
				longestP = len(f.name)
			}
		}
		for _, f := range cmd.persistentFlags {
			printFlag(f, w, longestP, "    ")
		}
	}

	// Collect and display persistent flags inherited from ancestors.
	var inherited []*Flag
	seen := make(map[string]bool)
	for _, f := range cmd.flags {
		seen[f.name] = true
	}
	for _, f := range cmd.persistentFlags {
		seen[f.name] = true
	}

	ancestor := cmd.parent
	for ancestor != nil {
		for _, f := range ancestor.persistentFlags {
			if !seen[f.name] {
				inherited = append(inherited, f)
				seen[f.name] = true
			}
		}
		ancestor = ancestor.parent
	}

	if len(inherited) > 0 {
		fmt.Fprintf(w, "  Inherited flags:\n")
		longestI := 0
		for _, f := range inherited {
			if len(f.name) > longestI {
				longestI = len(f.name)
			}
		}
		for _, f := range inherited {
			printFlag(f, w, longestI, "    ")
		}
	}

	// Print nested subcommands if any.
	if len(cmd.subcommands) > 0 {
		fmt.Fprintf(w, "  Subcommands:\n")
		for _, child := range cmd.subcommands {
			fmt.Fprintf(w, "    %s: %s\n", child.name, child.description)
		}
	}

	fmt.Fprintln(w)
}

// PrintUsage writes the full usage message (global flags + subcommand list) to w.
func (c *CLI) PrintUsage(w io.Writer) {
	longest := 0
	for _, f := range c.flags {
		if len(f.name) > longest {
			longest = len(f.name)
		}
	}

	programName := c.name
	if programName == "" {
		programName = os.Args[0]
	}

	fmt.Fprintf(w, "Usage: %s [global flags] [subcommand] [subcommand flags] [args...]\n", programName)
	if c.description != "" {
		fmt.Fprintf(w, "%s\n", c.description)
	}

	fmt.Fprintf(w, "\nGlobal Flags:\n")
	for _, f := range c.flags {
		printFlag(f, w, longest, "  ")
	}

	fmt.Fprintf(w, "\nSubcommands:\n")
	for _, cmd := range c.subcommands {
		printSubCommand(cmd, w)
	}
}

// Ensure completionCmd handler compiles — log.Fatal kept for CLI entry points
// only; internal package callers use the returned error.
var _ = log.Fatal
