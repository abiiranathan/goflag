package goflag

import (
	"fmt"
	"io"
	"reflect"
)

// SubCMD represents a subcommand in the CLI application. Subcommands may be
// nested arbitrarily deep; each SubCMD holds a pointer to its parent so that
// help output and completion generation can walk the full ancestry chain.
//
// Not safe for concurrent use during registration; safe for concurrent reads
// after Parse has returned.
type SubCMD struct {
	name        string             // Subcommand name used as the dispatch key.
	description string             // Human-readable description shown in help output.
	Handler     func(userdata any) // Invoked by ParseAndInvoke when this subcommand is matched.
	flags       []*Flag            // Flags that belong to this subcommand.
	subcommands []*SubCMD          // Nested subcommands registered under this one.
	parent      *SubCMD            // Parent subcommand; nil for top-level subcommands.
	args        []string           // Positional arguments collected after all flags are consumed.
}

// Args returns the positional arguments captured for this subcommand during
// the most recent Parse call. Positional arguments are any tokens that are not
// flag names, flag values, or recognised nested subcommand names.
func (cmd *SubCMD) Args() []string {
	return cmd.args
}

// Name returns the name of the subcommand.
func (cmd *SubCMD) Name() string {
	return cmd.name
}

// SubCommand registers a nested subcommand under cmd and returns the new child
// SubCMD so that its flags can be configured via method chaining.
//
// Panics if name or description is empty, or if handler is nil.
func (cmd *SubCMD) SubCommand(name, description string, handler func(userdata any)) *SubCMD {
	if handler == nil {
		panic("subcommand cannot be registered with a nil handler")
	}
	if name == "" {
		panic("subcommand name cannot be empty")
	}
	if description == "" {
		panic("subcommand description cannot be empty")
	}

	child := &SubCMD{
		name:        name,
		description: description,
		Handler:     handler,
		parent:      cmd,
		flags: []*Flag{
			{name: "help", shortName: "h", flagType: flagBool, usage: "Print help message and exit"},
		},
	}
	cmd.subcommands = append(cmd.subcommands, child)
	return child
}

// Flag adds a typed flag to the subcommand and returns the subcommand for
// continued method chaining. The final flag added is the one affected by
// subsequent Required() and Validate() calls.
func (cmd *SubCMD) Flag(ft flagType, name, shortName string, valuePtr any, usage string) *SubCMD {
	f := &Flag{
		flagType:   ft,
		name:       name,
		shortName:  shortName,
		value:      valuePtr,
		usage:      usage,
		validators: make([]FlagValidator, 0),
	}
	validateFlag(f)
	cmd.flags = append(cmd.flags, f)
	return cmd
}

// Required marks the most recently added flag as required.
func (cmd *SubCMD) Required() *SubCMD {
	if len(cmd.flags) > 0 {
		cmd.flags[len(cmd.flags)-1].required = true
	}
	return cmd
}

// Validate appends validators to the most recently added flag.
func (cmd *SubCMD) Validate(validators ...FlagValidator) *SubCMD {
	if len(cmd.flags) > 0 {
		last := cmd.flags[len(cmd.flags)-1]
		last.validators = append(last.validators, validators...)
	}
	return cmd
}

// PrintUsage writes usage information for this subcommand to w.
func (cmd *SubCMD) PrintUsage(w io.Writer) {
	printSubCommand(cmd, w)
}

// validateFlag checks that a Flag is well-formed. Panics on any violation
// so that programming errors surface immediately at registration time.
func validateFlag(flag *Flag) {
	if flag == nil {
		panic("flag cannot be nil")
	}
	if flag.name == "" {
		panic("flag name cannot be empty")
	}
	if flag.value == nil {
		panic("flag value pointer cannot be nil")
	}

	vt := reflect.TypeOf(flag.value)
	if vt.Kind() != reflect.Pointer {
		panic(fmt.Errorf("flag value for %q must be a pointer, got %s", flag.name, vt.Kind()))
	}
}
