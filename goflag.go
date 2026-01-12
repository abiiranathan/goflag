// A simple flag package for go.
// Support for --flag value and -flag values.
// Built in subcommand support and flag validation.
// Author: Dr. Abiira Nathan.
// Date: Sept 25. 2023
// License: MIT License
package goflag

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

//go:generate go tool stringer -type FlagType

type FlagType int

const (
	FlagString FlagType = iota
	FlagInt
	FlagInt64
	FlagFloat32
	FlagFloat64
	FlagBool
	FlagRune
	FlagDuration
	FlagStringSlice
	FlagIntSlice
	FlagTime
	FlagIP
	FlagMAC
	FlagURL
	FlagUUID
	FlagHostPortPair
	FlagEmail
	FlagFilePath
	FlagDirPath
)

// A Flag as parsed from the command line.
type Flag struct {
	flagType  FlagType
	name      string
	shortName string
	value     any // pointer to default value. Will be populated by Parse.
	usage     string
	required  bool
	validator func(any) (bool, string)
}

// Add validator to last flag in the subcommand chain. If no flag exists, it panics.
func (flag *Flag) Validate(validator func(any) (bool, string)) *Flag {
	flag.validator = validator
	return flag
}

func (flag *Flag) Required() *Flag {
	flag.required = true
	return flag
}

// Global flag context. Stores global flags and subcommands.
type CLI struct {
	flags       []*Flag
	subcommands []*subcommand
}

// The completion subcommand.
var completionCmd *subcommand

// Create a new command-line interface.
func New() *CLI {
	cli := &CLI{
		flags: []*Flag{
			{name: "help", shortName: "h", flagType: FlagBool, usage: "Print help message and exit"},
		},
	}

	// Add completion subcommand for all CLIs
	var shell string
	var install bool
	var uninstall bool

	completionCmd = cli.SubCommand("completion", "Generate shell completion scripts", func() {
		binName := filepath.Base(os.Args[0])

		if shell == "" {
			log.Fatal("Error: --shell flag is required\n")
		}

		// Check for conflicting flags
		if install && uninstall {
			log.Fatal("Error: cannot use --install and --uninstall together\n")
		}

		if uninstall {
			// Uninstall the completion script
			if err := uninstallCompletion(shell, binName); err != nil {
				log.Fatalf("Failed to uninstall completion: %v\n", err)
			}
		} else if install {
			// Install the completion script
			var generateFunc func(io.Writer)
			switch shell {
			case "bash":
				generateFunc = cli.GenBashCompletion
			case "zsh":
				generateFunc = cli.GenZshCompletion
			default:
				log.Fatalf("Unsupported shell: %s\n", shell)
			}

			if err := installCompletion(shell, binName, generateFunc); err != nil {
				log.Fatalf("Failed to install completion: %v\n", err)
			}
		} else {
			// Just print to stdout
			switch shell {
			case "bash":
				cli.GenBashCompletion(os.Stdout)
			case "zsh":
				cli.GenZshCompletion(os.Stdout)
			default:
				log.Fatalf("Unsupported shell: %s\n", shell)
			}
		}
	}).
		FlagString("shell", "s", &shell, "The shell to generate completions for [bash|zsh]").Required().
		FlagBool("install", "i", &install, "Install the completion script to the appropriate location").
		FlagBool("uninstall", "u", &uninstall, "Uninstall the completion script")

	return cli
}

// Add a flag to the context.
func (c *CLI) Flag(flagType FlagType, name, shortName string, valuePtr any, usage string) *Flag {
	flag := &Flag{
		flagType:  flagType,
		name:      name,
		shortName: shortName,
		value:     valuePtr,
		usage:     usage,
	}

	validateFlag(flag)
	c.flags = append(c.flags, flag)
	return flag
}

// Add a subcommand to the command-line context.
func (c *CLI) SubCommand(name, description string, handler func()) *subcommand {
	if handler == nil {
		panic("subcommand can not be registered with nil handler")
	}

	if name == "" {
		panic("subcommand name can't be empty")
	}
	if description == "" {
		panic("subcommand description can't be empty")
	}

	cmd := &subcommand{
		name:        name,
		description: description,
		Handler:     handler,
		flags: []*Flag{
			{name: "help", shortName: "h", flagType: FlagString, usage: "Print help message and exit"},
		},
	}
	c.subcommands = append(c.subcommands, cmd)

	// add the help flag to the subcommand.
	return cmd
}

// Parse the flags and subcommands. args should be os.Args.
// The first argument is ignored as it is the program name.
//
// Populates the values of the flags and also finds the matching subcommand.
// Returns the matching subcommand.
func (c *CLI) Parse(argv []string) (*subcommand, error) {
	var subcmd *subcommand = nil
	subCommandIndex := -1

	// store processed flags.
	processedGlobalFlags := make(map[string]bool)
	processedSubCommandFlags := make(map[string]bool)

	if len(argv) >= 2 {
		// skip the first argument which is the program name.
		argv = argv[1:]

		// First pass, consume global flags.
	outerloop:
		for i := 0; i < len(argv); i++ {
			arg := argv[i]

			if strings.TrimSpace(arg) == "" {
				continue
			}

			// Check for = in the arg. If present, split the arg into two.
			// The first part is the flag name and the second part is the value.
			// e.g. --name=John
			if strings.Contains(arg, "=") {
				parts := strings.Split(arg, "=")       // split the arg into two.
				arg = parts[0]                         // the first part is the flag name.
				argv = append(argv[:i+1], argv[i:]...) // insert the second part into the argv.
				argv[i+1] = parts[1]                   // set the second part as the next arg.

			}

			var name string

			if arg[0] == '-' && arg[1] == '-' {
				// long flag
				name = arg[2:]
			} else if arg[0] == '-' {
				// short flag
				name = arg[1:]
			} else {
				if len(c.subcommands) == 0 {
					continue
				}

				// subcommand or flag value.
				for _, cmd := range c.subcommands {
					if cmd.name == arg {
						subcmd = cmd
						subCommandIndex = i
						break outerloop
					}
				}
				continue // value will be consumed by looking at the next arg.
			}

			if isHelpFlag(name) {
				c.PrintUsage(os.Stdout)
				os.Exit(0)
			}

			flag, err := parseFlags(&c.flags, name, i, argv)
			if err != nil {
				return nil, err
			}

			if flag != nil {
				// Store the processed flag.
				// This is used to check if all required global flags are present.
				processedGlobalFlags[flag.name] = true
			}
		}
	}

	// check if all required global flags are present.
	// Done after parsing the subcommand flags so that the subcommand help can be printed.
	// if the global flags are missing.
	if subcmd != completionCmd {
		for _, flag := range c.flags {
			if _, found := processedGlobalFlags[flag.name]; !found && flag.required {
				return nil, fmt.Errorf("missing required flag [-%s | --%s]", flag.shortName, flag.name)
			}
		}
	}

	// Second pass, consume subcommand flags.
	if subcmd == nil {
		return nil, nil
	}

	// remove the subcommand from the argv.
	subCommandIndex++

	// parse the subcommand flags.
	for i := subCommandIndex; i < len(argv); i++ {
		arg := argv[i]
		if strings.TrimSpace(arg) == "" {
			continue
		}

		// Check for = in the arg. If present, split the arg into two.
		// The first part is the flag name and the second part is the value.
		// e.g. --name=John
		if strings.Contains(arg, "=") {
			parts := strings.Split(arg, "=")       // split the arg into two.
			arg = parts[0]                         // the first part is the flag name.
			argv = append(argv[:i+1], argv[i:]...) // insert the second part into the argv.
			argv[i+1] = parts[1]                   // set the second part as the next arg.
		}

		var name string
		if arg[0] == '-' && arg[1] == '-' {
			// long flag
			name = arg[2:]
		} else if arg[0] == '-' {
			// short flag
			name = arg[1:]
		} else {
			continue
			// flag value. will be consumed by looking at the next arg.
		}

		if isHelpFlag(name) {
			subcmd.PrintUsage(os.Stdout)
			os.Exit(0)
		}

		flag, err := parseFlags(&subcmd.flags, name, i, argv)
		if err != nil {
			return nil, err
		}

		if flag != nil {
			// Store the processed flag.
			// This is used to check if all required subcommand flags are present.
			processedSubCommandFlags[flag.name] = true
		}
	}

	// check if all required subcommand flags are present.
	for _, flag := range subcmd.flags {
		if _, found := processedSubCommandFlags[flag.name]; !found && flag.required {
			return nil, fmt.Errorf("missing required flag [-%s | --%s]", flag.shortName, flag.name)
		}
	}

	return subcmd, nil
}

// Helper to Parse the flags.
// flags: The flags to parse.
// name: The name of the flag, may be the short name.
// i: The index of the flag in the argv.
// argv: The arguments.
func parseFlags(flags *[]*Flag, name string, i int, argv []string) (*Flag, error) {
	flag := findFlag(*flags, name)
	if flag == nil {
		return nil, fmt.Errorf("unknown flag : %s", name)
	}

	// look at the next arg for the value.
	valueIndex := i + 1
	if (valueIndex) >= len(argv) {
		if flag.flagType == FlagBool { // bool falg may have no value associated. e.g. --verbose
			*flag.value.(*bool) = true
			return flag, nil
		}
		return flag, fmt.Errorf("missing value for flag [-%s | --%s]", flag.shortName, flag.name)
	}

	if argv[valueIndex] == "" {
		// empty string, accessing argv[valueIndex][0] will panic.
		return flag, fmt.Errorf("empty value for flag [-%s | --%s]", flag.shortName, flag.name)
	}

	if argv[valueIndex][0] == '-' {
		if flag.flagType == FlagBool { // bool falg may have no value.
			*flag.value.(*bool) = true
			return flag, nil
		}
		return flag, fmt.Errorf("missing value for flag [-%s | --%s]", flag.shortName, flag.name)
	}

	var err error
	value := argv[valueIndex]
	err = parseFlagValue(flag, value)
	if err != nil {
		return flag, err
	}

	// validate the flag.
	if flag.validator != nil {
		// dereference the pointer to get the value.
		value := reflect.ValueOf(flag.value).Elem().Interface()
		if valid, errMsg := flag.validator(value); !valid {
			return flag, fmt.Errorf("invalid value (%v) for flag [-%s | --%s]: %v", flag.value, flag.shortName, flag.name, errMsg)
		}
	}
	return flag, nil
}

// Print a flag to the writer.
// Called by PrintUsage for each flag.
func printFlag(flag *Flag, w io.Writer, longestFlagName int, indent string) {
	fmt.Fprintf(w, "%s--%-*s ", indent, longestFlagName, flag.name)
	valid := reflect.ValueOf(flag.value).IsValid()
	value := ""
	if valid {
		value = fmt.Sprintf("%v", reflect.ValueOf(flag.value).Elem().Interface())
	}

	if flag.flagType == FlagString {
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

// Parse the flag value and set the flag value.
func findFlag(flags []*Flag, name string) *Flag {
	for index := range flags {
		flag := flags[index]
		if flag.name == name || flag.shortName == name {
			return flag
		}
	}
	return nil
}

func isHelpFlag(name string) bool {
	return name == "help" || name == "h"
}

// Print a subcommand to the writer.
// Called by PrintUsage for each subcommand.
func printSubCommand(cmd *subcommand, w io.Writer) {
	fmt.Fprintf(w, "%s: %s", cmd.name, cmd.description)
	fmt.Fprintln(w)

	longestFlagName := 0
	for _, flag := range cmd.flags {
		if flag.name == "help" {
			continue
		}
		if len(flag.name) > longestFlagName {
			longestFlagName = len(flag.name)
		}
	}

	// print the subcommand flags.
	for _, flag := range cmd.flags {
		if flag.name == "help" {
			continue
		}
		printFlag(flag, w, longestFlagName, "    ")
	}

	fmt.Fprintln(w)
}

// Print the usage to the writer.
// Called by Parse if the help flag is present.
// help flag is automatically added to the context.
// May be called as help, --help, -h, --h
//
// Help for a given subcommand can be printed by passing the subcommand name as the
// glag --subcommand or -c. e.g. --help -c greet
func (c *CLI) PrintUsage(w io.Writer) {
	longestFlagName := 0
	for _, flag := range c.flags {
		if len(flag.name) > longestFlagName {
			longestFlagName = len(flag.name)
		}
	}

	// find the longest subcommand name.
	longestSubCommandName := 0
	for _, cmd := range c.subcommands {
		if len(cmd.name) > longestSubCommandName {
			longestSubCommandName = len(cmd.name)
		}
	}

	fmt.Fprintf(w, "Usage: %s [global flags] [subcommand] [subcommand flags]\n", os.Args[0])
	// print the global flags.
	fmt.Fprintf(w, "Global Flags:\n")
	for _, flag := range c.flags {
		printFlag(flag, w, longestFlagName, "  ")
	}

	fmt.Fprintln(w)

	// print the subcommands.
	fmt.Fprintf(w, "Subcommands:\n")
	for _, cmd := range c.subcommands {
		printSubCommand(cmd, w)
	}
}
