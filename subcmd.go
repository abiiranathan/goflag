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

// Add validator to last flag in the subcommand chain.
func (cmd *subcommand) Validate(validators ...FlagValidator) *subcommand {
	if len(cmd.flags) > 0 {
		cmd.flags[len(cmd.flags)-1].validators = append(cmd.flags[len(cmd.flags)-1].validators, validators...)
	}
	return cmd
}

func (cmd *subcommand) Required() *subcommand {
	if len(cmd.flags) > 0 {
		cmd.flags[len(cmd.flags)-1].required = true
	}
	return cmd
}

// Add a flag to a subcommand.
func (cmd *subcommand) Flag(flagType flagType, name, shortName string, valuePtr any, usage string) *subcommand {
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

func (cmd *subcommand) PrintUsage(w io.Writer) {
	printSubCommand(cmd, w)
}

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
	if valueType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("flag value for %s must be a pointer, got %s", flag.name, valueType.Kind()))
	}
}
