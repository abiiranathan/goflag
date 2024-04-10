package goflag

import (
	"fmt"
	"io"
	"reflect"
)

// A Subcommand. It can have its own flags.
type Subcommand struct {
	name        string  // Subcommand name. used as a key to find the subcommand.
	description string  // Description of what this subcommand does.
	Handler     func()  // Subcommand callback handler. Will be invoked by user if it matches.
	flags       []*Flag // subcommand flags.
}

// Add a flag to a subcommand.
func (cmd *Subcommand) AddFlag(flagType FlagType, name, shortName string, valuePtr any, usage string, required bool, validator ...func(any) (bool, string)) *Subcommand {
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
	cmd.flags = append(cmd.flags, flag)
	return cmd
}

func (cmd *Subcommand) PrintUsage(w io.Writer) {
	printSubCommand(cmd, w)
}

// Create a standalone flag that can be shared across multiple subcommands.
// Call AddFlagPtr to add the flag to a subcommand.
func NewFlag(flagType FlagType, name, shortName string, valuePtr any, usage string, required bool, validator ...func(any) (bool, string)) *Flag {
	flag := &Flag{
		Name:      name,
		ShortName: shortName,
		Value:     valuePtr,
		Usage:     usage,
		Required:  required,
	}

	if len(validator) > 0 {
		flag.Validator = validator[0]
	}
	return flag
}

// Add a flag to a subcommand.
func (cmd *Subcommand) AddFlagPtr(flag *Flag) *Subcommand {
	validateFlag(flag)
	cmd.flags = append(cmd.flags, flag)
	return cmd
}

func validateFlag(flag *Flag) {
	if flag == nil {
		panic("flag can't be nil")
	}

	if flag.Name == "" {
		panic("flag name can't be empty")
	}

	// check the flag value is a valid pointer.
	if flag.Value == nil {
		panic("flag value can't be nil")
	}

	valueType := reflect.TypeOf(flag.Value)
	if valueType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("flag value for %s must be a pointer, got %s", flag.Name, valueType.Kind()))
	}
}
