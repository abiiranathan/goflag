package goflag

import (
	"fmt"
	"io"
	"reflect"
)

// A subcommand. It can have its own flags.
type subcommand struct {
	name        string  // Subcommand name. used as a key to find the subcommand.
	description string  // Description of what this subcommand does.
	Handler     func()  // Subcommand callback handler. Will be invoked by user if it matches.
	flags       []*Flag // subcommand flags.
}

// Add a flag to a subcommand.
func (cmd *subcommand) AddFlag(flagType FlagType, name, shortName string, valuePtr any, usage string, required bool, validator ...func(any) (bool, string)) *subcommand {
	if name == "" {
		panic("flag name can't be empty")
	}

	// check the flag value is a valid pointer.
	if valuePtr == nil {
		panic("flag value can't be nil")
	}

	valueType := reflect.TypeOf(valuePtr)
	if valueType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("flag value for %s must be a pointer, got %s", name, valueType.Kind()))
	}

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
	cmd.flags = append(cmd.flags, flag)
	return cmd
}

func (cmd *subcommand) PrintUsage(w io.Writer) {
	printSubCommand(cmd, w)
}
