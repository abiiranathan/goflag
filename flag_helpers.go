package goflag

// Helper on top-level CLI
func (c *CLI) FlagString(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagString, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagInt(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagInt, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagInt64(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagInt64, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagFloat32(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagFloat32, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagFloat64(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagFloat64, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagBool(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagBool, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagRune(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagRune, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagDuration(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagDuration, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagStringSlice(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagStringSlice, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagIntSlice(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagIntSlice, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagTime(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagTime, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagIP(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagFloat32, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagMAC(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagMAC, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagURL(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagURL, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagUUID(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagUUID, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagHostPortPair(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagHostPortPair, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagEmail(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagEmail, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagFilePath(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagFilePath, name, shortName, valuePtr, usage)
}

func (c *CLI) FlagDirPath(name, shortName string, valuePtr any, usage string) *Flag {
	return c.Flag(FlagDirPath, name, shortName, valuePtr, usage)
}

// Helper on subcommand
func (cmd *subcommand) FlagString(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagString, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagInt(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagInt, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagInt64(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagInt64, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagFloat32(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagFloat32, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagFloat64(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagFloat64, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagBool(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagBool, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagRune(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagRune, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagDuration(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagDuration, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagStringSlice(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagStringSlice, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagIntSlice(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagIntSlice, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagTime(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagTime, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagIP(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagFloat32, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagMAC(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagMAC, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagURL(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagURL, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagUUID(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagUUID, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagHostPortPair(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagHostPortPair, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagEmail(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagEmail, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagFilePath(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagFilePath, name, shortName, valuePtr, usage)
}

func (cmd *subcommand) FlagDirPath(name, shortName string, valuePtr any, usage string) *subcommand {
	return cmd.Flag(FlagDirPath, name, shortName, valuePtr, usage)
}
