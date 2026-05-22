package goflag

import (
	"fmt"
	"io"
	"reflect"
)

// SubCMD represents a subcommand in the CLI application.
// It contains the subcommand name, description, handler function, and associated flags.
type SubCMD struct {
	name        string  // Subcommand name. used as a key to find the subcommand.
	description string  // Description of what this subcommand does.
	Handler     func()  // Subcommand callback handler. Will be invoked by user if it matches.
	flags       []*Flag // subcommand flags.
}

// Validate adds validators to the latest added flag of the subcommand.
func (cmd *SubCMD) Validate(validators ...FlagValidator) *SubCMD {
	if len(cmd.flags) > 0 {
		cmd.flags[len(cmd.flags)-1].validators = append(cmd.flags[len(cmd.flags)-1].validators, validators...)
	}
	return cmd
}

// Required sets the latest added flag as required.
func (cmd *SubCMD) Required() *SubCMD {
	if len(cmd.flags) > 0 {
		cmd.flags[len(cmd.flags)-1].required = true
	}
	return cmd
}

// Flag adds a flag to the subcommand.
func (cmd *SubCMD) Flag(flagType flagType, name, shortName string, valuePtr any, usage string) *SubCMD {
	flag := &Flag{
		flagType:   flagType,
		name:       name,
		shortName:  shortName,
		value:      valuePtr,
		usage:      usage,
		validators: make([]FlagValidator, 0),
	}

	validateFlag(flag)
	cmd.flags = append(cmd.flags, flag)
	return cmd
}

// PrintUsage prints the usage information for the subcommand to the provided writer.
func (cmd *SubCMD) PrintUsage(w io.Writer) {
	printSubCommand(cmd, w)
}

// Name returns the name of the subcommand.
func (cmd *SubCMD) Name() string {
	return cmd.name
}

// validateFlag checks if the flag is valid. It panics if the flag is invalid.
func validateFlag(flag *Flag) {
	if flag == nil {
		panic("flag can't be nil")
	}

	if flag.name == "" {
		panic("flag name can't be empty")
	}

	// check the flag value is a valid pointer.
	if flag.value == nil {
		panic("flag value can't be nil")
	}

	valueType := reflect.TypeOf(flag.value)
	if valueType.Kind() != reflect.Pointer {
		panic(fmt.Errorf("flag value for %s must be a pointer, got %s", flag.name, valueType.Kind()))
	}
}
