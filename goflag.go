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
	"os"
	"strings"
)

//go:generate stringer -type flagType
//go:generate go run cmd/genc.go -c "github.com/abiiranathan/goflag" constructors.go
//go:generate go run cmd/genc.go -g "github.com/abiiranathan/goflag" getters.go

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

// Subcommand callback handler.
// Will be invoked by user if it matches the subcommand returned by ctx.Parse.
type Handler func(ctx Getter, cmd Getter)

// A Flag as parsed from the command line.
type gflag struct {
	name      string
	shortName string
	value     any // The value of the flag. Value provided initially is the default.
	flagType  flagType
	usage     string
	required  bool
	validator func(any) (bool, string)
}

// Add validator to last flag in the subcommand chain. If no flag exists, it panics.
func (flag *gflag) Validate(validator func(any) (bool, string)) *gflag {
	flag.validator = validator
	return flag
}

// Global flag context. Stores global flags and subcommands.
type Context struct {
	flags       []*gflag
	subcommands []*subcommand
}

// Create a new flag context.
func NewContext() *Context {
	return &Context{
		flags: []*gflag{
			{name: "help", shortName: "h", flagType: flagBool, usage: "Print help message and exit"},
		},
		subcommands: []*subcommand{
			{
				name:        "completions",
				description: "Generate shell completions[zsh, bash]",
				Handler:     GenerateCompletions,
				flags: []*gflag{
					{name: "zsh", flagType: flagBool, value: false, usage: "Generate Zsh Completion"},
					{name: "bash", flagType: flagBool, value: false, usage: "Generate Zsh Completion"},
					{name: "out", flagType: flagString, value: "", usage: "Output file. Default [stdout]"},
				},
			},
		},
	}
}

// Add a flag to the context.
func (ctx *Context) AddFlag(flag *gflag) *gflag {
	if flag.name == "" {
		panic("flag name can't be empty")
	}
	ctx.flags = append(ctx.flags, flag)
	return flag
}

func (ctx *Context) Get(name string) any {
	for _, flag := range ctx.flags {
		if flag.name == name {
			return flag.value
		}
		if flag.shortName != "" && flag.shortName == name {
			return flag.value
		}
	}

	panic(fmt.Sprintf("flag %q not found", name))
}

// Add a subcommand to the context.
func (ctx *Context) AddSubCommand(cmd *subcommand) *subcommand {
	if cmd.Handler == nil {
		panic("subcommand can not be registered with nil handler")
	}

	if cmd.name == "" {
		panic("subcommand name can't be empty")
	}

	ctx.subcommands = append(ctx.subcommands, cmd)

	// add the help flag to the subcommand.
	cmd.AddFlag(&gflag{
		name:      "help",
		shortName: "h",
		flagType:  flagBool,
		usage:     "Print help message and exit",
	})

	return cmd
}

// Parse the flags and subcommands. args should be os.Args.
// The first argument is ignored as it is the program name.
//
// Populates the values of the flags and also finds the matching subcommand.
// Returns the matching subcommand.
func (ctx *Context) Parse(argv []string) (*subcommand, error) {
	if (argv == nil) || (len(argv) <= 1) {
		return nil, fmt.Errorf("can not call Parse() without at least 2 arguments")
	}

	// skip the first argument which is the program name.
	argv = argv[1:]
	var subcmd *subcommand = nil
	subCommandIndex := -1

	// store processed flags.
	processedGlobalFlags := make(map[string]bool)
	processedSubCommandFlags := make(map[string]bool)

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
			if len(ctx.subcommands) == 0 {
				continue
			}

			// subcommand or flag value.
			for _, cmd := range ctx.subcommands {
				if cmd.name == arg {
					subcmd = cmd
					subCommandIndex = i
					break outerloop
				}
			}
			continue // value will be consumed by looking at the next arg.
		}

		if isHelpFlag(name) {
			ctx.PrintUsage(os.Stdout)
			os.Exit(0)
		}

		flag, err := parseFlags(&ctx.flags, name, i, argv)
		if err != nil {
			return nil, err
		}

		if flag != nil {
			// Store the processed flag.
			// This is used to check if all required global flags are present.
			processedGlobalFlags[flag.name] = true
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

	// check if all required global flags are present.
	// Done after parsing the subcommand flags so that the subcommand help can be printed.
	// if the global flags are missing.
	for _, flag := range ctx.flags {
		if flag.required && !processedGlobalFlags[flag.name] {
			return nil, fmt.Errorf("missing required flag %q", flag.name)
		}
	}

	// check if all required subcommand flags are present.
	for _, flag := range subcmd.flags {
		if flag.required && !processedSubCommandFlags[flag.name] {
			return nil, fmt.Errorf("missing required flag %q", flag.name)
		}
	}
	return subcmd, nil
}

// Helper to Parse the flags.
// flags: The flags to parse.
// name: The name of the flag, may be the short name.
// i: The index of the flag in the argv.
// argv: The arguments.
func parseFlags(flags *[]*gflag, name string, i int, argv []string) (*gflag, error) {
	flag := findFlag(*flags, name)
	if flag == nil {
		return nil, fmt.Errorf("unknown flag : %s", name)
	}

	// look at the next arg for the value.
	valueIndex := i + 1
	if (valueIndex) >= len(argv) {
		if flag.flagType == flagBool { // bool falg may have no value associated. e.g. --verbose
			flag.value = true
			return flag, nil
		}
		return flag, fmt.Errorf("missing value for flag %q", flag.name)
	}

	if argv[valueIndex] == "" {
		// empty string, accessing argv[valueIndex][0] will panic.
		return flag, fmt.Errorf("empty value for flag %q", flag.name)
	}

	if argv[valueIndex][0] == '-' {
		if flag.flagType == flagBool { // bool falg may have no value.
			flag.value = true
			return flag, nil
		}
		return flag, fmt.Errorf("missing value for flag %q", flag.name)
	}

	var err error
	value := argv[valueIndex]
	flag.value, err = parseFlagValue(flag, value)
	if err != nil {
		return flag, err
	}

	// validate the flag.
	if flag.validator != nil {
		if valid, errMsg := flag.validator(flag.value); !valid {
			return flag, fmt.Errorf("invalid value(%v) for flag %q: %s", flag.value, flag.name, errMsg)
		}
	}
	return flag, nil
}

// Print a flag to the writer.
// Called by PrintUsage for each flag.
func printFlag(flag *gflag, w io.Writer, longestFlagName int, indent string) {
	fmt.Fprintf(w, "%s--%-*s ", indent, longestFlagName, flag.name)

	if flag.shortName != "" {
		fmt.Fprintf(w, "-%s: %s\n", flag.shortName, flag.usage)
	} else {
		fmt.Fprintf(w, "%s\n", flag.usage)
	}
}

func findFlag(flags []*gflag, name string) *gflag {
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
func (ctx *Context) PrintUsage(w io.Writer) {
	longestFlagName := 0
	for _, flag := range ctx.flags {
		if len(flag.name) > longestFlagName {
			longestFlagName = len(flag.name)
		}
	}

	// find the longest subcommand name.
	longestSubCommandName := 0
	for _, cmd := range ctx.subcommands {
		if len(cmd.name) > longestSubCommandName {
			longestSubCommandName = len(cmd.name)
		}
	}

	fmt.Fprintf(w, "Usage: %s [global flags] [subcommand] [subcommand flags]\n", os.Args[0])
	// print the global flags.
	fmt.Fprintf(w, "Global Flags:\n")
	for _, flag := range ctx.flags {
		printFlag(flag, w, longestFlagName, "  ")
	}

	fmt.Fprintln(w)

	// print the subcommands.
	fmt.Fprintf(w, "Subcommands:\n")
	for _, cmd := range ctx.subcommands {
		printSubCommand(cmd, w)
	}
}
