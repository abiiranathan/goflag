package goflag

import (
	"net"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Helper methods on top-level CLI for defining flags with specific types.
// These methods provide a convenient way to add typed flags without explicitly
// specifying the FlagType constant.

// String adds a string flag to the CLI.
// Parameters:
//   - name: The long name of the flag (e.g., "output")
//   - shortName: The short name of the flag (e.g., "o"), can be empty
//   - valuePtr: Pointer to a string variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration (e.g., setting required status).
func (c *CLI) String(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagString, name, shortName, valuePtr, usage)
}

// Int adds an integer flag to the CLI.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to an int variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Int(name, shortName string, valuePtr *int, usage string) *Flag {
	return c.addFlag(flagInt, name, shortName, valuePtr, usage)
}

// Int64 adds a 64-bit integer flag to the CLI.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to an int64 variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Int64(name, shortName string, valuePtr *int64, usage string) *Flag {
	return c.addFlag(flagInt64, name, shortName, valuePtr, usage)
}

// Float32 adds a 32-bit floating point flag to the CLI.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a float32 variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Float32(name, shortName string, valuePtr *float32, usage string) *Flag {
	return c.addFlag(flagFloat32, name, shortName, valuePtr, usage)
}

// Float64 adds a 64-bit floating point flag to the CLI.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a float64 variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Float64(name, shortName string, valuePtr *float64, usage string) *Flag {
	return c.addFlag(flagFloat64, name, shortName, valuePtr, usage)
}

// Bool adds a boolean flag to the CLI.
// Boolean flags don't require a value; their presence sets them to true.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a bool variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Bool(name, shortName string, valuePtr *bool, usage string) *Flag {
	return c.addFlag(flagBool, name, shortName, valuePtr, usage)
}

// Rune adds a single Unicode character flag to the CLI.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a rune variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Rune(name, shortName string, valuePtr *rune, usage string) *Flag {
	return c.addFlag(flagRune, name, shortName, valuePtr, usage)
}

// Duration adds a time.Duration flag to the CLI.
// Accepts duration strings like "5s", "2m", "1h30m", etc.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a time.Duration variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Duration(name, shortName string, valuePtr *time.Duration, usage string) *Flag {
	return c.addFlag(flagDuration, name, shortName, valuePtr, usage)
}

// StringSlice adds a string slice flag to the CLI.
// Can be specified multiple times to build a list of values.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a []string variable where the parsed values will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) StringSlice(name, shortName string, valuePtr *[]string, usage string) *Flag {
	return c.addFlag(flagStringSlice, name, shortName, valuePtr, usage)
}

// IntSlice adds an integer slice flag to the CLI.
// Can be specified multiple times to build a list of integer values.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a []int variable where the parsed values will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) IntSlice(name, shortName string, valuePtr *[]int, usage string) *Flag {
	return c.addFlag(flagIntSlice, name, shortName, valuePtr, usage)
}

// Time adds a time.Time flag to the CLI.
// Accepts various time formats for parsing.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a time.Time variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Time(name, shortName string, valuePtr *time.Time, usage string) *Flag {
	return c.addFlag(flagTime, name, shortName, valuePtr, usage)
}

// IP adds an IP address flag to the CLI.
// Accepts both IPv4 and IPv6 addresses.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a net.IP variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) IP(name, shortName string, valuePtr *net.IP, usage string) *Flag {
	return c.addFlag(flagIP, name, shortName, valuePtr, usage)
}

// MAC adds a MAC address flag to the CLI.
// Accepts MAC addresses in standard formats (e.g., "01:23:45:67:89:ab").
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a net.HardwareAddr variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) MAC(name, shortName string, valuePtr *net.HardwareAddr, usage string) *Flag {
	return c.addFlag(flagMAC, name, shortName, valuePtr, usage)
}

// URL adds a URL flag to the CLI.
// Validates and parses URLs according to RFC 3986.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a *url.URL variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) URL(name, shortName string, valuePtr *url.URL, usage string) *Flag {
	return c.addFlag(flagURL, name, shortName, valuePtr, usage)
}

// UUID adds a UUID flag to the CLI.
// Accepts UUIDs in standard format (e.g., "550e8400-e29b-41d4-a716-446655440000").
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a string or UUID type where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) UUID(name, shortName string, valuePtr *uuid.UUID, usage string) *Flag {
	return c.addFlag(flagUUID, name, shortName, valuePtr, usage)
}

// HostPortPair adds a host:port pair flag to the CLI.
// Accepts values like "localhost:8080" or "192.168.1.1:443".
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a string variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) HostPortPair(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagHostPortPair, name, shortName, valuePtr, usage)
}

// Email adds an email address flag to the CLI.
// Validates basic email format.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a string variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) Email(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagEmail, name, shortName, valuePtr, usage)
}

// FilePath adds a file path flag to the CLI.
// Can optionally validate that the file exists.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a string variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) FilePath(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagFilePath, name, shortName, valuePtr, usage)
}

// DirPath adds a directory path flag to the CLI.
// Can optionally validate that the directory exists.
// Parameters:
//   - name: The long name of the flag
//   - shortName: The short name of the flag, can be empty
//   - valuePtr: Pointer to a string variable where the parsed value will be stored
//   - usage: Description of the flag shown in help text
//
// Returns the created Flag for further configuration.
func (c *CLI) DirPath(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagDirPath, name, shortName, valuePtr, usage)
}

// Helper methods on subcommand for defining flags with specific types.
// These mirror the CLI-level helpers but return *subcommand for method chaining.

// String adds a string flag to the subcommand.
// See CLI.String for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) String(name, shortName string, valuePtr *string, usage string) *subcommand {
	return cmd.Flag(flagString, name, shortName, valuePtr, usage)
}

// Int adds an integer flag to the subcommand.
// See CLI.Int for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Int(name, shortName string, valuePtr *int, usage string) *subcommand {
	return cmd.Flag(flagInt, name, shortName, valuePtr, usage)
}

// Int64 adds a 64-bit integer flag to the subcommand.
// See CLI.Int64 for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Int64(name, shortName string, valuePtr *int64, usage string) *subcommand {
	return cmd.Flag(flagInt64, name, shortName, valuePtr, usage)
}

// Float32 adds a 32-bit floating point flag to the subcommand.
// See CLI.Float32 for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Float32(name, shortName string, valuePtr *float32, usage string) *subcommand {
	return cmd.Flag(flagFloat32, name, shortName, valuePtr, usage)
}

// Float64 adds a 64-bit floating point flag to the subcommand.
// See CLI.Float64 for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Float64(name, shortName string, valuePtr *float64, usage string) *subcommand {
	return cmd.Flag(flagFloat64, name, shortName, valuePtr, usage)
}

// Bool adds a boolean flag to the subcommand.
// See CLI.Bool for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Bool(name, shortName string, valuePtr *bool, usage string) *subcommand {
	return cmd.Flag(flagBool, name, shortName, valuePtr, usage)
}

// Rune adds a single Unicode character flag to the subcommand.
// See CLI.Rune for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Rune(name, shortName string, valuePtr *rune, usage string) *subcommand {
	return cmd.Flag(flagRune, name, shortName, valuePtr, usage)
}

// Duration adds a time.Duration flag to the subcommand.
// See CLI.Duration for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Duration(name, shortName string, valuePtr *time.Duration, usage string) *subcommand {
	return cmd.Flag(flagDuration, name, shortName, valuePtr, usage)
}

// StringSlice adds a string slice flag to the subcommand.
// See CLI.StringSlice for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) StringSlice(name, shortName string, valuePtr *[]string, usage string) *subcommand {
	return cmd.Flag(flagStringSlice, name, shortName, valuePtr, usage)
}

// IntSlice adds an integer slice flag to the subcommand.
// See CLI.IntSlice for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) IntSlice(name, shortName string, valuePtr *[]int, usage string) *subcommand {
	return cmd.Flag(flagIntSlice, name, shortName, valuePtr, usage)
}

// Time adds a time.Time flag to the subcommand.
// See CLI.Time for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Time(name, shortName string, valuePtr *time.Time, usage string) *subcommand {
	return cmd.Flag(flagTime, name, shortName, valuePtr, usage)
}

// IP adds an IP address flag to the subcommand.
// See CLI.IP for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) IP(name, shortName string, valuePtr *net.IP, usage string) *subcommand {
	return cmd.Flag(flagFloat32, name, shortName, valuePtr, usage) // BUG: Should be FlagIP
}

// MAC adds a MAC address flag to the subcommand.
// See CLI.MAC for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) MAC(name, shortName string, valuePtr *net.HardwareAddr, usage string) *subcommand {
	return cmd.Flag(flagMAC, name, shortName, valuePtr, usage)
}

// URL adds a URL flag to the subcommand.
// See CLI.URL for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) URL(name, shortName string, valuePtr *url.URL, usage string) *subcommand {
	return cmd.Flag(flagURL, name, shortName, valuePtr, usage)
}

// UUID adds a UUID flag to the subcommand.
// See CLI.UUID for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) UUID(name, shortName string, valuePtr *uuid.UUID, usage string) *subcommand {
	return cmd.Flag(flagUUID, name, shortName, valuePtr, usage)
}

// HostPortPair adds a host:port pair flag to the subcommand.
// See CLI.HostPortPair for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) HostPortPair(name, shortName string, valuePtr *string, usage string) *subcommand {
	return cmd.Flag(flagHostPortPair, name, shortName, valuePtr, usage)
}

// Email adds an email address flag to the subcommand.
// See CLI.Email for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) Email(name, shortName string, valuePtr *string, usage string) *subcommand {
	return cmd.Flag(flagEmail, name, shortName, valuePtr, usage)
}

// FilePath adds a file path flag to the subcommand.
// See CLI.FilePath for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) FilePath(name, shortName string, valuePtr *string, usage string) *subcommand {
	return cmd.Flag(flagFilePath, name, shortName, valuePtr, usage)
}

// DirPath adds a directory path flag to the subcommand.
// See CLI.DirPath for parameter details.
// Returns the subcommand for method chaining.
func (cmd *subcommand) DirPath(name, shortName string, valuePtr *string, usage string) *subcommand {
	return cmd.Flag(flagDirPath, name, shortName, valuePtr, usage)
}
