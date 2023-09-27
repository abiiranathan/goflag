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
	"os"

	"github.com/abiiranathan/goflag"
)

func greetUser(ctx goflag.Getter, cmd goflag.Getter) {
	name := cmd.GetString("name")
	greeting := cmd.GetString("greeting")
	fmt.Println(greeting, name)

	// you have access to global flags
	ctx.GetBool("verbose") // etc

}

func printVersion(ctx *goflag.Context, cmd *goflag.Subcmd) {
	if cmd.Get("short").(bool) {
		fmt.Println("1.0.0")
	} else {
		fmt.Println("1.0.0")
		fmt.Println("Build Date: 2021-01-01")
		fmt.Println("Commit: 1234567890")
	}
}

func main() {
	ctx := goflag.NewContext()

	ctx.AddFlag(goflag.String("config", "c", "config.json", "Path to config file", true))
	ctx.AddFlag(goflag.Bool("verbose", "v", false, "Enable verbose output", false))

	ctx.AddSubCommand(goflag.SubCommand("greet", "Greet a person", greetUser)).
		AddFlag(goflag.String("name", "n", "World", "Name of the person to greet", true)).
		AddFlag(goflag.String("greeting", "g", "Hello", "Greeting to use", false)).
		Validate(func(a any) (bool, string) {
			if a.(string) == "World" {
				return false, "ERR: Name cannot be World"
			}
			return true, ""
		}).AddFlag(goflag.Bool("upper", "u", false, "Print in upper case", false))

	ctx.AddSubCommand(goflag.SubCommand("version", "Print version", printVersion)).
		AddFlag(goflag.Bool("verbose", "v", false, "Enable verbose output", false)).
		AddFlag(goflag.Bool("short", "s", false, "Print short version", false))

	// Parse the command line arguments and return the matching subcommand
	subcmd, err := ctx.Parse(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

	if subcmd != nil {
		subcmd.Handler(ctx, subcmd)
	}

	fmt.Println(ctx.GetString("config"))
	fmt.Println(ctx.GetBool("verbose"))

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
Usage: ./cli [flags] [subcommand] [flags]

Global Flags:
  --help     -h: Print help message and exit
  --config   -c: Path to config file
  --verbose  -v: Enable verbose output
  --timeout  -t: Timeout for the request
  --port     -p: Port to listen on
  --start    -s: Start time
  --url      -u: URL to fetch
  --uuid     -i: UUID to use
  --ip       -i: IP to use
  --mac      -m: MAC address to use
  --email    -e: Email address to use
  --hostport -h: Host:Port pair to use
  --file     -f: File path to use
  --dir      -d: Directory path to use

Subcommands:
  greet   : Greet a person
         --help     -h: Print help message and exit
         --name     -n: Name of the person to greet
         --greeting -g: Greeting to use
         --upper    -u: Print in upper case

  version : Print version
         --help    -h: Print help message and exit
         --verbose -v: Enable verbose output
         --short   -s: Print short version

  sleep   : Sleep for a while
         --help -h: Print help message and exit
         --time -t: Time to sleep in seconds

  cors    : Enable CORS
         --help        -h: Print help message and exit
         --origins     -o: Allowed origins
         --methods     -m: Allowed methods
         --headers     -d: Allowed headers
         --credentials -c: Allow credentials


```

## Contributing

Contributions are welcome. Please open an issue to discuss your ideas before opening a PR.

### License

MIT

## TODO

- [ ] Add support for subcommands of subcommands.
- [ ] Implement more tests
