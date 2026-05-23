package goflag

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// SubCMD represents a subcommand in the CLI application. Subcommands may be
// nested arbitrarily deep; each SubCMD holds a pointer to its parent so that
// help output, persistent-flag inheritance, and completion generation can walk
// the full ancestry chain.
//
// Not safe for concurrent use during registration; safe for concurrent reads
// after Parse has returned.
type SubCMD struct {
	name                  string                   // Subcommand name used as the dispatch key.
	description           string                   // Human-readable description shown in help output.
	Handler               func(userdata any) error // Invoked by ParseAndInvoke when this subcommand is matched.
	flags                 []*Flag                  // Flags that belong to this subcommand.
	persistentFlags       []*Flag                  // Flags inherited by all nested subcommands.
	subcommands           []*SubCMD                // Nested subcommands registered under this one.
	parent                *SubCMD                  // Parent subcommand; nil for top-level subcommands.
	args                  []string                 // Positional arguments collected after all flags are consumed.
	exclusiveGroups       [][]string               // Sets of mutually exclusive flag names for this subcommand.
	lastAddedIsPersistent bool                     // Tracks whether the last added flag was persistent.
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
func (cmd *SubCMD) SubCommand(name, description string, handler func(userdata any) error) *SubCMD {
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
	cmd.lastAddedIsPersistent = false
	return cmd
}

// PersistentFlag adds a typed flag that is automatically visible to all nested
// subcommands (children, grandchildren, etc.) registered under cmd. It returns
// cmd for method chaining; Required() and Validate() target the last added
// persistent flag.
func (cmd *SubCMD) PersistentFlag(ft flagType, name, shortName string, valuePtr any, usage string) *SubCMD {
	f := &Flag{
		flagType:   ft,
		name:       name,
		shortName:  shortName,
		value:      valuePtr,
		usage:      usage,
		validators: make([]FlagValidator, 0),
	}
	validateFlag(f)
	cmd.persistentFlags = append(cmd.persistentFlags, f)
	cmd.lastAddedIsPersistent = true
	return cmd
}

// Required marks the most recently added flag (own or persistent) as required.
func (cmd *SubCMD) Required() *SubCMD {
	last := cmd.lastAddedFlag()
	if last != nil {
		last.required = true
	}
	return cmd
}

// Validate appends validators to the most recently added flag (own or persistent).
func (cmd *SubCMD) Validate(validators ...FlagValidator) *SubCMD {
	last := cmd.lastAddedFlag()
	if last != nil {
		last.validators = append(last.validators, validators...)
	}
	return cmd
}

// ExclusiveFlags declares that at most one flag from names may be provided in
// a single invocation of this subcommand. Returns cmd for chaining.
func (cmd *SubCMD) ExclusiveFlags(names ...string) *SubCMD {
	if len(names) < 2 {
		panic("ExclusiveFlags requires at least two flag names")
	}
	group := make([]string, len(names))
	copy(group, names)
	cmd.exclusiveGroups = append(cmd.exclusiveGroups, group)
	return cmd
}

// lastAddedFlag returns a pointer to the last flag appended to either
// cmd.flags or cmd.persistentFlags, whichever was most recently modified.
func (cmd *SubCMD) lastAddedFlag() *Flag {
	if cmd.lastAddedIsPersistent {
		nPersist := len(cmd.persistentFlags)
		if nPersist > 0 {
			return cmd.persistentFlags[nPersist-1]
		}
	} else {
		nOwn := len(cmd.flags)
		if nOwn > 0 {
			return cmd.flags[nOwn-1]
		}
	}

	// Fallback in case of inconsistency
	nOwn := len(cmd.flags)
	nPersist := len(cmd.persistentFlags)
	if nOwn > 0 {
		return cmd.flags[nOwn-1]
	}
	if nPersist > 0 {
		return cmd.persistentFlags[nPersist-1]
	}
	return nil
}

// PrintUsage writes usage information for this subcommand to w.
func (cmd *SubCMD) PrintUsage(w io.Writer) {
	printSubCommand(cmd, w)
}

// validateFlag checks that a Flag is well-formed. Panics on any violation.
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

// ---------------------------------------------------------------------------
// Typed PersistentFlag helpers on SubCMD
// ---------------------------------------------------------------------------

// PersistentString adds a persistent string flag to the subcommand.
func (cmd *SubCMD) PersistentString(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.PersistentFlag(flagString, name, shortName, valuePtr, usage)
}

// PersistentInt adds a persistent integer flag to the subcommand.
func (cmd *SubCMD) PersistentInt(name, shortName string, valuePtr *int, usage string) *SubCMD {
	return cmd.PersistentFlag(flagInt, name, shortName, valuePtr, usage)
}

// PersistentInt64 adds a persistent 64-bit integer flag to the subcommand.
func (cmd *SubCMD) PersistentInt64(name, shortName string, valuePtr *int64, usage string) *SubCMD {
	return cmd.PersistentFlag(flagInt64, name, shortName, valuePtr, usage)
}

// PersistentFloat32 adds a persistent 32-bit float flag to the subcommand.
func (cmd *SubCMD) PersistentFloat32(name, shortName string, valuePtr *float32, usage string) *SubCMD {
	return cmd.PersistentFlag(flagFloat32, name, shortName, valuePtr, usage)
}

// PersistentFloat64 adds a persistent 64-bit float flag to the subcommand.
func (cmd *SubCMD) PersistentFloat64(name, shortName string, valuePtr *float64, usage string) *SubCMD {
	return cmd.PersistentFlag(flagFloat64, name, shortName, valuePtr, usage)
}

// PersistentBool adds a persistent boolean flag to the subcommand.
func (cmd *SubCMD) PersistentBool(name, shortName string, valuePtr *bool, usage string) *SubCMD {
	return cmd.PersistentFlag(flagBool, name, shortName, valuePtr, usage)
}

// PersistentRune adds a persistent rune flag to the subcommand.
func (cmd *SubCMD) PersistentRune(name, shortName string, valuePtr *rune, usage string) *SubCMD {
	return cmd.PersistentFlag(flagRune, name, shortName, valuePtr, usage)
}

// PersistentDuration adds a persistent time.Duration flag to the subcommand.
func (cmd *SubCMD) PersistentDuration(name, shortName string, valuePtr *time.Duration, usage string) *SubCMD {
	return cmd.PersistentFlag(flagDuration, name, shortName, valuePtr, usage)
}

// PersistentStringSlice adds a persistent comma-separated string slice flag to the subcommand.
func (cmd *SubCMD) PersistentStringSlice(name, shortName string, valuePtr *[]string, usage string) *SubCMD {
	return cmd.PersistentFlag(flagStringSlice, name, shortName, valuePtr, usage)
}

// PersistentIntSlice adds a persistent comma-separated integer slice flag to the subcommand.
func (cmd *SubCMD) PersistentIntSlice(name, shortName string, valuePtr *[]int, usage string) *SubCMD {
	return cmd.PersistentFlag(flagIntSlice, name, shortName, valuePtr, usage)
}

// PersistentTime adds a persistent time.Time flag to the subcommand.
func (cmd *SubCMD) PersistentTime(name, shortName string, valuePtr *time.Time, usage string) *SubCMD {
	return cmd.PersistentFlag(flagTime, name, shortName, valuePtr, usage)
}

// PersistentIP adds a persistent IP address flag to the subcommand.
func (cmd *SubCMD) PersistentIP(name, shortName string, valuePtr *net.IP, usage string) *SubCMD {
	return cmd.PersistentFlag(flagIP, name, shortName, valuePtr, usage)
}

// PersistentMAC adds a persistent hardware (MAC) address flag to the subcommand.
func (cmd *SubCMD) PersistentMAC(name, shortName string, valuePtr *net.HardwareAddr, usage string) *SubCMD {
	return cmd.PersistentFlag(flagMAC, name, shortName, valuePtr, usage)
}

// PersistentURL adds a persistent URL flag to the subcommand.
func (cmd *SubCMD) PersistentURL(name, shortName string, valuePtr *url.URL, usage string) *SubCMD {
	return cmd.PersistentFlag(flagURL, name, shortName, valuePtr, usage)
}

// PersistentUUID adds a persistent UUID flag to the subcommand.
func (cmd *SubCMD) PersistentUUID(name, shortName string, valuePtr *uuid.UUID, usage string) *SubCMD {
	return cmd.PersistentFlag(flagUUID, name, shortName, valuePtr, usage)
}

// PersistentHostPortPair adds a persistent host:port flag to the subcommand.
func (cmd *SubCMD) PersistentHostPortPair(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.PersistentFlag(flagHostPortPair, name, shortName, valuePtr, usage)
}

// PersistentEmail adds a persistent email address flag to the subcommand.
func (cmd *SubCMD) PersistentEmail(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.PersistentFlag(flagEmail, name, shortName, valuePtr, usage)
}

// PersistentFilePath adds a persistent file path flag to the subcommand.
func (cmd *SubCMD) PersistentFilePath(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.PersistentFlag(flagFilePath, name, shortName, valuePtr, usage)
}

// PersistentDirPath adds a persistent directory path flag to the subcommand.
func (cmd *SubCMD) PersistentDirPath(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.PersistentFlag(flagDirPath, name, shortName, valuePtr, usage)
}
