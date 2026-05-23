package goflag

import (
	"net"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Helper methods on CLI for defining typed global flags.
// Each method wraps addFlag with the appropriate flagType constant and returns
// the created Flag so the caller can chain Required() or Validate().

// String adds a string flag to the CLI.
func (c *CLI) String(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagString, name, shortName, valuePtr, usage)
}

// Int adds an integer flag to the CLI.
func (c *CLI) Int(name, shortName string, valuePtr *int, usage string) *Flag {
	return c.addFlag(flagInt, name, shortName, valuePtr, usage)
}

// Int64 adds a 64-bit integer flag to the CLI.
func (c *CLI) Int64(name, shortName string, valuePtr *int64, usage string) *Flag {
	return c.addFlag(flagInt64, name, shortName, valuePtr, usage)
}

// Float32 adds a 32-bit floating-point flag to the CLI.
func (c *CLI) Float32(name, shortName string, valuePtr *float32, usage string) *Flag {
	return c.addFlag(flagFloat32, name, shortName, valuePtr, usage)
}

// Float64 adds a 64-bit floating-point flag to the CLI.
func (c *CLI) Float64(name, shortName string, valuePtr *float64, usage string) *Flag {
	return c.addFlag(flagFloat64, name, shortName, valuePtr, usage)
}

// Bool adds a boolean flag to the CLI. Boolean flags are set to true by their
// presence alone; no explicit value is required.
func (c *CLI) Bool(name, shortName string, valuePtr *bool, usage string) *Flag {
	return c.addFlag(flagBool, name, shortName, valuePtr, usage)
}

// Rune adds a single Unicode character flag to the CLI.
func (c *CLI) Rune(name, shortName string, valuePtr *rune, usage string) *Flag {
	return c.addFlag(flagRune, name, shortName, valuePtr, usage)
}

// Duration adds a time.Duration flag to the CLI.
// Accepts duration strings such as "5s", "2m30s", "1h".
func (c *CLI) Duration(name, shortName string, valuePtr *time.Duration, usage string) *Flag {
	return c.addFlag(flagDuration, name, shortName, valuePtr, usage)
}

// StringSlice adds a comma-separated string slice flag to the CLI.
func (c *CLI) StringSlice(name, shortName string, valuePtr *[]string, usage string) *Flag {
	return c.addFlag(flagStringSlice, name, shortName, valuePtr, usage)
}

// IntSlice adds a comma-separated integer slice flag to the CLI.
func (c *CLI) IntSlice(name, shortName string, valuePtr *[]int, usage string) *Flag {
	return c.addFlag(flagIntSlice, name, shortName, valuePtr, usage)
}

// Time adds a time.Time flag to the CLI.
func (c *CLI) Time(name, shortName string, valuePtr *time.Time, usage string) *Flag {
	return c.addFlag(flagTime, name, shortName, valuePtr, usage)
}

// IP adds an IP address flag to the CLI. Accepts both IPv4 and IPv6.
func (c *CLI) IP(name, shortName string, valuePtr *net.IP, usage string) *Flag {
	return c.addFlag(flagIP, name, shortName, valuePtr, usage)
}

// MAC adds a hardware (MAC) address flag to the CLI.
func (c *CLI) MAC(name, shortName string, valuePtr *net.HardwareAddr, usage string) *Flag {
	return c.addFlag(flagMAC, name, shortName, valuePtr, usage)
}

// URL adds a URL flag to the CLI, validated via url.ParseRequestURI.
func (c *CLI) URL(name, shortName string, valuePtr *url.URL, usage string) *Flag {
	return c.addFlag(flagURL, name, shortName, valuePtr, usage)
}

// UUID adds a UUID flag to the CLI.
func (c *CLI) UUID(name, shortName string, valuePtr *uuid.UUID, usage string) *Flag {
	return c.addFlag(flagUUID, name, shortName, valuePtr, usage)
}

// HostPortPair adds a host:port flag to the CLI (e.g. "localhost:8080").
func (c *CLI) HostPortPair(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagHostPortPair, name, shortName, valuePtr, usage)
}

// Email adds an email address flag to the CLI.
func (c *CLI) Email(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagEmail, name, shortName, valuePtr, usage)
}

// FilePath adds a file path flag to the CLI. The file must exist at parse time.
func (c *CLI) FilePath(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagFilePath, name, shortName, valuePtr, usage)
}

// DirPath adds a directory path flag to the CLI. The directory must exist at parse time.
func (c *CLI) DirPath(name, shortName string, valuePtr *string, usage string) *Flag {
	return c.addFlag(flagDirPath, name, shortName, valuePtr, usage)
}

// ---------------------------------------------------------------------------
// Helper methods on SubCMD — mirror the CLI helpers but return *SubCMD so
// that nested flag chains remain on the subcommand.
// ---------------------------------------------------------------------------

// String adds a string flag to the subcommand.
func (cmd *SubCMD) String(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.Flag(flagString, name, shortName, valuePtr, usage)
}

// Int adds an integer flag to the subcommand.
func (cmd *SubCMD) Int(name, shortName string, valuePtr *int, usage string) *SubCMD {
	return cmd.Flag(flagInt, name, shortName, valuePtr, usage)
}

// Int64 adds a 64-bit integer flag to the subcommand.
func (cmd *SubCMD) Int64(name, shortName string, valuePtr *int64, usage string) *SubCMD {
	return cmd.Flag(flagInt64, name, shortName, valuePtr, usage)
}

// Float32 adds a 32-bit floating-point flag to the subcommand.
func (cmd *SubCMD) Float32(name, shortName string, valuePtr *float32, usage string) *SubCMD {
	return cmd.Flag(flagFloat32, name, shortName, valuePtr, usage)
}

// Float64 adds a 64-bit floating-point flag to the subcommand.
func (cmd *SubCMD) Float64(name, shortName string, valuePtr *float64, usage string) *SubCMD {
	return cmd.Flag(flagFloat64, name, shortName, valuePtr, usage)
}

// Bool adds a boolean flag to the subcommand.
func (cmd *SubCMD) Bool(name, shortName string, valuePtr *bool, usage string) *SubCMD {
	return cmd.Flag(flagBool, name, shortName, valuePtr, usage)
}

// Rune adds a single Unicode character flag to the subcommand.
func (cmd *SubCMD) Rune(name, shortName string, valuePtr *rune, usage string) *SubCMD {
	return cmd.Flag(flagRune, name, shortName, valuePtr, usage)
}

// Duration adds a time.Duration flag to the subcommand.
func (cmd *SubCMD) Duration(name, shortName string, valuePtr *time.Duration, usage string) *SubCMD {
	return cmd.Flag(flagDuration, name, shortName, valuePtr, usage)
}

// StringSlice adds a comma-separated string slice flag to the subcommand.
func (cmd *SubCMD) StringSlice(name, shortName string, valuePtr *[]string, usage string) *SubCMD {
	return cmd.Flag(flagStringSlice, name, shortName, valuePtr, usage)
}

// IntSlice adds a comma-separated integer slice flag to the subcommand.
func (cmd *SubCMD) IntSlice(name, shortName string, valuePtr *[]int, usage string) *SubCMD {
	return cmd.Flag(flagIntSlice, name, shortName, valuePtr, usage)
}

// Time adds a time.Time flag to the subcommand.
func (cmd *SubCMD) Time(name, shortName string, valuePtr *time.Time, usage string) *SubCMD {
	return cmd.Flag(flagTime, name, shortName, valuePtr, usage)
}

// IP adds an IP address flag to the subcommand.
// Previously registered with flagFloat32 by mistake — fixed.
func (cmd *SubCMD) IP(name, shortName string, valuePtr *net.IP, usage string) *SubCMD {
	return cmd.Flag(flagIP, name, shortName, valuePtr, usage)
}

// MAC adds a hardware (MAC) address flag to the subcommand.
func (cmd *SubCMD) MAC(name, shortName string, valuePtr *net.HardwareAddr, usage string) *SubCMD {
	return cmd.Flag(flagMAC, name, shortName, valuePtr, usage)
}

// URL adds a URL flag to the subcommand.
func (cmd *SubCMD) URL(name, shortName string, valuePtr *url.URL, usage string) *SubCMD {
	return cmd.Flag(flagURL, name, shortName, valuePtr, usage)
}

// UUID adds a UUID flag to the subcommand.
func (cmd *SubCMD) UUID(name, shortName string, valuePtr *uuid.UUID, usage string) *SubCMD {
	return cmd.Flag(flagUUID, name, shortName, valuePtr, usage)
}

// HostPortPair adds a host:port flag to the subcommand.
func (cmd *SubCMD) HostPortPair(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.Flag(flagHostPortPair, name, shortName, valuePtr, usage)
}

// Email adds an email address flag to the subcommand.
func (cmd *SubCMD) Email(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.Flag(flagEmail, name, shortName, valuePtr, usage)
}

// FilePath adds a file path flag to the subcommand.
func (cmd *SubCMD) FilePath(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.Flag(flagFilePath, name, shortName, valuePtr, usage)
}

// DirPath adds a directory path flag to the subcommand.
func (cmd *SubCMD) DirPath(name, shortName string, valuePtr *string, usage string) *SubCMD {
	return cmd.Flag(flagDirPath, name, shortName, valuePtr, usage)
}
