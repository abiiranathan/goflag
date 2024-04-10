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
	"reflect"
	"strings"
)

//go:generate stringer -type FlagType

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
	FlagType  FlagType
	Name      string
	ShortName string
	Value     any // pointer to default value. Will be populated by Parse.
	Usage     string
	Required  bool
	Validator func(any) (bool, string)
}

// Add validator to last flag in the subcommand chain. If no flag exists, it panics.
func (flag *Flag) Validate(validator func(any) (bool, string)) *Flag {
	flag.Validator = validator
	return flag
}

// Global flag context. Stores global flags and subcommands.
type Context struct {
	flags       []*Flag
	subcommands []*Subcommand
}

// Create a new flag context.
func NewContext() *Context {
	return &Context{
		flags: []*Flag{
			{Name: "help", ShortName: "h", FlagType: FlagBool, Usage: "Print help message and exit"},
		},
	}
}

// Add a flag to the context.
func (ctx *Context) AddFlag(flagType FlagType, name, shortName string, valuePtr any, usage string, required bool, validator ...func(any) (bool, string)) *Flag {
	flag := &Flag{
		FlagType:  flagType,
		Name:      name,
		ShortName: shortName,
		Value:     valuePtr,
		Usage:     usage,
		Required:  required,
	}

	if len(validator) > 0 {
		flag.Validator = validator[0]
	}

	validateFlag(flag)
	ctx.flags = append(ctx.flags, flag)
	return flag
}

// Add a flag to a subcommand.
func (ctx *Context) AddFlagPtr(flag *Flag) *Flag {
	validateFlag(flag)
	ctx.flags = append(ctx.flags, flag)
	return flag
}

// Add a subcommand to the context.
func (ctx *Context) AddSubCommand(name, description string, handler func()) *Subcommand {
	if handler == nil {
		panic("subcommand can not be registered with nil handler")
	}

	if name == "" {
		panic("subcommand name can't be empty")
	}
	if description == "" {
		panic("subcommand description can't be empty")
	}

	cmd := &Subcommand{
		name:        name,
		description: description,
		Handler:     handler,
		flags: []*Flag{
			{Name: "help", ShortName: "h", FlagType: FlagString, Usage: "Print help message and exit"},
		},
	}
	ctx.subcommands = append(ctx.subcommands, cmd)

	// add the help flag to the subcommand.
	return cmd
}

// Parse the flags and subcommands. args should be os.Args.
// The first argument is ignored as it is the program name.
//
// Populates the values of the flags and also finds the matching subcommand.
// Returns the matching subcommand.
func (ctx *Context) Parse(argv []string) (*Subcommand, error) {
	var subcmd *Subcommand = nil
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
				processedGlobalFlags[flag.Name] = true
			}
		}
	}

	// check if all required global flags are present.
	// Done after parsing the subcommand flags so that the subcommand help can be printed.
	// if the global flags are missing.
	for _, flag := range ctx.flags {
		if _, found := processedGlobalFlags[flag.Name]; !found && flag.Required {
			return nil, fmt.Errorf("missing required flag %q", flag.Name)
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
			processedSubCommandFlags[flag.Name] = true
		}
	}

	// check if all required subcommand flags are present.
	for _, flag := range subcmd.flags {
		if _, found := processedSubCommandFlags[flag.Name]; !found && flag.Required {
			return nil, fmt.Errorf("missing required flag %q", flag.Name)
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
		if flag.FlagType == FlagBool { // bool falg may have no value associated. e.g. --verbose
			*flag.Value.(*bool) = true
			return flag, nil
		}
		return flag, fmt.Errorf("missing value for flag %q", flag.Name)
	}

	if argv[valueIndex] == "" {
		// empty string, accessing argv[valueIndex][0] will panic.
		return flag, fmt.Errorf("empty value for flag %q", flag.Name)
	}

	if argv[valueIndex][0] == '-' {
		if flag.FlagType == FlagBool { // bool falg may have no value.
			*flag.Value.(*bool) = true
			return flag, nil
		}
		return flag, fmt.Errorf("missing value for flag %q", flag.Name)
	}

	var err error
	value := argv[valueIndex]
	err = parseFlagValue(flag, value)
	if err != nil {
		return flag, err
	}

	// validate the flag.
	if flag.Validator != nil {
		// dereference the pointer to get the value.
		value := reflect.ValueOf(flag.Value).Elem().Interface()
		if valid, errMsg := flag.Validator(value); !valid {
			return flag, fmt.Errorf("invalid value(%v) for flag %q: %s", flag.Value, flag.Name, errMsg)
		}
	}
	return flag, nil
}

// Print a flag to the writer.
// Called by PrintUsage for each flag.
func printFlag(flag *Flag, w io.Writer, longestFlagName int, indent string) {
	fmt.Fprintf(w, "%s--%-*s ", indent, longestFlagName, flag.Name)
	valid := reflect.ValueOf(flag.Value).IsValid()
	value := ""
	if valid {
		value = fmt.Sprintf("%v", reflect.ValueOf(flag.Value).Elem().Interface())
	}

	if flag.ShortName != "" {
		fmt.Fprintf(w, "-%s: %s (default: %s)\n", flag.ShortName, flag.Usage, value)
	} else {
		fmt.Fprintf(w, "%s (default: %s)\n", flag.Usage, value)
	}
}

// Parse the flag value and set the flag value.
func findFlag(flags []*Flag, name string) *Flag {
	for index := range flags {
		flag := flags[index]
		if flag.Name == name || flag.ShortName == name {
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
func printSubCommand(cmd *Subcommand, w io.Writer) {
	fmt.Fprintf(w, "%s: %s", cmd.name, cmd.description)
	fmt.Fprintln(w)

	longestFlagName := 0
	for _, flag := range cmd.flags {
		if flag.Name == "help" {
			continue
		}
		if len(flag.Name) > longestFlagName {
			longestFlagName = len(flag.Name)
		}
	}

	// print the subcommand flags.
	for _, flag := range cmd.flags {
		if flag.Name == "help" {
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
		if len(flag.Name) > longestFlagName {
			longestFlagName = len(flag.Name)
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
