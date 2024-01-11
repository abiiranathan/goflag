# goflag

[![GoDoc](https://godoc.org/github.com/bradleyjkemp/goflag?status.svg)](https://godoc.org/github.com/bradleyjkemp/goflag)

A simple library for parsing command line arguments. It is designed to be a drop-in replacement for the standard library's `flag` package.

The main difference is that `goflag`

- supports both short and long flags
- requires no global state
- has a beatiful API based on the builder pattern that allows you to define your subcommands and flags in a single expression
- supports subcommands. Subcommands can have their own flags. However, subcommands cannot have subcommands of their own (yet).
- Beautifully formatted help text, with support for subcommand help text too.
- Supports custom validation of flags and arguments.
- Supports custom flag types with -, -- and = syntax.

## Installation

```bash
go get -u github.com/abiiranathan/goflag
```

## Usage

```go

package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/abiiranathan/goflag"
	"github.com/google/uuid"
)

var (
	name     string = "World"
	greeting string = "Hello"
	short    bool

	urlValue url.URL
	uuidVal  uuid.UUID
	ipVal    net.IP
	macVal   net.HardwareAddr
	emailVal string
	hpVal    string
	fileVal  string
	dirVal   string

	origins       []string = []string{"*"}
	methods       []string = []string{"GET", "POST"}
	headers       []string = []string{"Content-Type"}
	credentials   bool
	verbose       bool
	config        string        = "config.json"
	port          int           = 8080
	start         time.Time     = time.Now()
	timeout       time.Duration = 5 * time.Second
	durationValue               = 5

	upperValue bool
)

func greetUser() {
	fmt.Println(greeting, name)
}

func printVersion() {
	if short {
		fmt.Println("1.0.0")
	} else {
		fmt.Println("1.0.0")
		fmt.Println("Build Date: 2021-01-01")
		fmt.Println("Commit: 1234567890")
	}
}

func handleSleep() {
	time.Sleep(time.Duration(durationValue) * time.Second)
}

func handleCors() {
	fmt.Println("Origins: ", origins)
	fmt.Println("Methods: ", methods)
	fmt.Println("Headers: ", headers)
	fmt.Println("Credentials: ", credentials)
}

func main() {
	log.SetFlags(log.Lshortfile)
	ctx := goflag.NewContext()

	ctx.AddFlag(goflag.FlagString, "config", "c", &config, "Path to config file", false)
	ctx.AddFlag(goflag.FlagBool, "verbose", "v", &verbose, "Enable verbose output", false)
	ctx.AddFlag(goflag.FlagDuration, "timeout", "t", &timeout, "Timeout for the request", false)
	ctx.AddFlag(goflag.FlagInt, "port", "p", &port, "Port to listen on", false)
	ctx.AddFlag(goflag.FlagString, "hostport", "h", &hpVal, "Host:Port to listen on", false)
	ctx.AddFlag(goflag.FlagTime, "start", "s", &start, "Start time", false)
	ctx.AddFlag(goflag.FlagURL, "url", "u", &urlValue, "URL to fetch", false)
	ctx.AddFlag(goflag.FlagUUID, "uuid", "i", &uuidVal, "UUID to use", false)
	ctx.AddFlag(goflag.FlagIP, "ip", "i", &ipVal, "IP to use", false)
	ctx.AddFlag(goflag.FlagMAC, "mac", "m", &macVal, "MAC address to use", false)
	ctx.AddFlag(goflag.FlagEmail, "email", "e", &emailVal, "Email address to use", false)
	ctx.AddFlag(goflag.FlagFilePath, "file", "f", &fileVal, "File path to use", false)
	ctx.AddFlag(goflag.FlagDirPath, "dir", "d", &dirVal, "Directory path to use", false)

	ctx.AddSubCommand("greet", "Greet a person", greetUser).
		AddFlag(goflag.FlagString, "name", "n", &name, "Name of the person to greet", true).
		AddFlag(goflag.FlagString, "greeting", "g", &greeting, "Greeting to use", false).
		AddFlag(goflag.FlagBool, "upper", "u", &upperValue, "Print in upper case", false)

	ctx.AddSubCommand("version", "Print version", printVersion).
		AddFlag(goflag.FlagBool, "verbose", "v", &verbose, "Enable verbose output", false).
		AddFlag(goflag.FlagBool, "short", "s", &short, "Print short version", false)

	ctx.AddSubCommand("sleep", "Sleep for a while", handleSleep).
		AddFlag(goflag.FlagInt, "time", "t", &durationValue, "Time to sleep in seconds", true)

	ctx.AddSubCommand("cors", "Enable CORS", handleCors).
		AddFlag(goflag.FlagStringSlice, "origins", "o", &origins, "Allowed origins", true).
		AddFlag(goflag.FlagStringSlice, "methods", "m", &methods, "Allowed methods", true).
		AddFlag(goflag.FlagStringSlice, "headers", "d", &headers, "Allowed headers", true).
		AddFlag(goflag.FlagBool, "credentials", "c", &credentials, "Allow credentials", false)

	// Parse the command line arguments and return the matching subcommand
	subcmd, err := ctx.Parse(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

	if subcmd != nil {
		subcmd.Handler()
	}

	// Print the values
	fmt.Println("Config: ", config)
	fmt.Println("Verbose: ", verbose)
	fmt.Println("Timeout: ", timeout)
	fmt.Println("Port: ", port)
	fmt.Println("Start: ", start)

	fmt.Println("URL: ", urlValue)
	fmt.Println("UUID: ", uuidVal)
	fmt.Println("IP: ", ipVal)
	fmt.Println("MAC: ", macVal)
	fmt.Println("Email: ", emailVal)
	fmt.Println("HostPort: ", hpVal)
	fmt.Println("File: ", fileVal)
	fmt.Println("Dir: ", dirVal)

	fmt.Println("Origins: ", origins)
	fmt.Println("Methods: ", methods)
	fmt.Println("Headers: ", headers)
	fmt.Println("Credentials: ", credentials)

	fmt.Println("Name: ", name)
	fmt.Println("Greeting: ", greeting)
	fmt.Println("Short: ", short)
	fmt.Println("Duration: ", durationValue)

}


```

---

## Accessing flags values.

Note that ctx.Parse() returns the matching subcommand. If no subcommand is matched, it returns nil.
The subcommand handler should be called with the context and the subcommand as arguments.

The handler can then access the flags and arguments using the `Get` or `GetString`, `GetBool` etc methods on either the context or the subcommand.

> The context carries the global flags and the subcommand carries the subcommand specific flags.

Supported flag types are

- [x] `string`
- [x] `bool`
- [x] `int`
- [x] `int64`
- [x] `float64`
- [x] `float32`
- [x] `time.Duration` with format `1h2m3s`, `1h`, `1h2m`, `1m2s`, `1m`, `1s` as supported by the standard library's `time.ParseDuration` function.
- [x] `time.Time` with format `2006-01-02T15:04 MST`
- [x] `rune`
- [x] `[]string` (comma separated list of strings) with format `a,b,c,d`
- [x] `[]int` (comma separated list of ints) with format `1,2,3,4`
- [x] `ip` (IP address) with format `xxx.xxx.xxx.xxx`
- [x] `mac` (MAC address) with format `xx:xx:xx:xx:xx:xx`
- [x] `hostport` (IP address with port) pair with format `host:port`
- [x] `path` (file path) with format. Will be converted to be an absolute path and will be validated to exist.
- [x] `url` with format `scheme://host:port/path?query#fragment`
- [x] `email` with format `local@domain`. Validated using the standard library's `mail.ParseAddress` function.
- [x] `uuid` with format `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` (UUID v4 with [Google's UUID package](http://github.com/google/uuid))

See [Example](./cmd/examples/example.go) for more details.
Run the example with `./test.sh` to see the output.

### Preview of the help text

```bash
Usage: ./cli --help
```

## Contributing

Contributions are welcome. Please open an issue to discuss your ideas before opening a PR.

### License

MIT

## TODO

- [ ] Add support for subcommands of subcommands.
- [ ] Implement more tests
