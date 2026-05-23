# goflag

A simple, type-safe command-line flag parsing library for Go with support for subcommands, custom flag validation, persistent flags, and mutually exclusive flags.

## Features

- **Type-Safe Flags**: Dedicated registration methods for strings, slices, numbers, networks, emails, UUIDs, files, and more.
- **Nested Subcommands**: Arbitrarily deep nested subcommand trees.
- **Improved Handler Signature**: Handlers match the `func(userdata any) error` signature, enabling robust error propagation and context passing.
- **Direct Dispatching**: `ParseAndInvoke` automatically executes matched subcommand handlers, allowing arbitrary configurations or services to be threaded through as user data.
- **Persistent (Inherited) Flags**: Declare persistent flags on subcommands that are automatically parsed and made available to all downstream subcommands.
- **Mutually Exclusive Groups**: Define sets of flags where at most one can be provided.
- **Rich Flag Validators**: Chain built-in validators or apply custom validation logic to any parsed flag value.
- **Shell Completions**: Automated generation of Bash and Zsh completion scripts that fully support persistent flags and subcommand hierarchies.

## Installation

```bash
go get github.com/abiiranathan/goflag
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abiiranathan/goflag"
)

type AppConfig struct {
	Verbose bool
}

func main() {
	var (
		config  string
		verbose bool
		timeout time.Duration
		port    int
	)

	// Create CLI with application name and description
	cli := goflag.New("myapp", "A simple command-line utility")

	// Define global flags
	cli.String("config", "c", &config, "Path to config file").Required()
	cli.Bool("verbose", "v", &verbose, "Enable verbose output")
	cli.Duration("timeout", "t", &timeout, "Request timeout")
	cli.Int("port", "p", &port, "Port to listen on")

	// Option A: Traditional manual execution
	subcmd, err := cli.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Config: %s\n", config)
	fmt.Printf("Verbose: %v\n", verbose)

	if subcmd != nil && subcmd.Handler != nil {
		if err := subcmd.Handler(nil); err != nil {
			log.Fatal(err)
		}
	}

	// Option B: Automated dispatching using ParseAndInvoke with threaded userdata
	appConfig := &AppConfig{Verbose: verbose}
	err = cli.ParseAndInvoke(os.Args, appConfig, func(cmd *goflag.SubCMD, userdata any) {
		// Optional pre-invoke callback
		if cmd != nil {
			fmt.Printf("Dispatched subcommand: %s\n", cmd.Name())
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

## Subcommands

Create arbitrarily nested subcommands with their own custom handler execution and argument configurations.

```go
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abiiranathan/goflag"
)

var (
	name     string = "World"
	greeting string = "Hello"
	upper    bool
)

func greetUser(userdata any) error {
	message := fmt.Sprintf("%s, %s!", greeting, name)
	if upper {
		message = strings.ToUpper(message)
	}
	fmt.Println(message)
	return nil
}

func main() {
	cli := goflag.New("greeter", "A CLI to greet people")

	// Register a subcommand
	cli.SubCommand("greet", "Greet a person", greetUser).
		String("name", "n", &name, "Name of the person to greet").Required().
		String("greeting", "g", &greeting, "Greeting to use").
		Bool("upper", "u", &upper, "Print in upper case")

	_, err := cli.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
```

```bash
$ greeter greet --name John --upper
HELLO, JOHN!

$ greeter greet -n Alice -g "Good morning"
Good morning, Alice!
```

## Persistent Flags

Persistent flags are defined on a subcommand and are automatically visible to, and parseable by, all descendants (children, grandchildren, etc.).

```go
var deployEnv string

deployCmd := cli.SubCommand("deploy", "Deploy the application", func(userdata any) error {
	fmt.Println("Please specify a target environment: staging or prod")
	return nil
})

// Registered once on deploy command, inherited by both staging and prod subcommands.
deployCmd.PersistentString("env", "e", &deployEnv, "Deployment target environment (e.g. staging, prod)").Required()

deployCmd.SubCommand("staging", "Deploy to staging", func(userdata any) error {
	fmt.Printf("Deploying staging with target environment configuration: %s\n", deployEnv)
	return nil
})

deployCmd.SubCommand("prod", "Deploy to production", func(userdata any) error {
	fmt.Printf("Deploying production with target environment configuration: %s\n", deployEnv)
	return nil
})
```

## Mutually Exclusive Flags

Declare constraints to ensure that at most one flag from a specific group is provided. 

```go
var (
	exportJSON bool
	exportYAML bool
)

cli.SubCommand("export", "Export dataset", handleExport).
	Bool("json", "j", &exportJSON, "Output as JSON").
	Bool("yaml", "y", &exportYAML, "Output as YAML").
	ExclusiveFlags("json", "yaml") // Parser errors if both flags are supplied
```

## Custom Flag Validation

Attach custom criteria or use built-in helpers to validate any flag's deserialized value.

```go
cli.Int("port", "p", &port, "Target server port").
	Required().
	Validate(goflag.Range(1024, 65535))

cli.String("shell", "s", &shell, "Target generation shell").
	Validate(goflag.Choices([]string{"bash", "zsh"}))
```

Available built-in validation helpers:
- `Choices([]T)`
- `MinStringLen(length)`
- `MaxStringLen(length)`
- `Min(minValue)`
- `Max(maxValue)`
- `Range(minValue, maxValue)`
- `NotEmpty`

## Supported Flag Types

### Basic Types
```go
cli.String("name", "n", &str, "String value")
cli.Int("count", "c", &num, "Integer value")
cli.Int64("big", "b", &big, "64-bit integer")
cli.Float32("ratio", "r", &f32, "32-bit float")
cli.Float64("pi", "p", &f64, "64-bit float")
cli.Bool("verbose", "v", &verbose, "Boolean flag")
cli.Rune("char", "c", &char, "Single character")
```

### Slices
```go
cli.StringSlice("origins", "o", &origins, "Allowed origins")
cli.IntSlice("ports", "p", &ports, "Port numbers")
```

### Time & Network
```go
cli.Duration("timeout", "t", &duration, "Duration (e.g. 5s, 2m)")
cli.Time("start", "s", &start, "Timestamp")
cli.IP("address", "a", &ip, "IP address")
cli.MAC("mac", "m", &mac, "MAC address")
cli.URL("endpoint", "e", &url, "URL")
cli.HostPortPair("listen", "l", &hostport, "Host:Port pair")
```

### System & Filesystem
```go
cli.Email("contact", "c", &email, "Email address")
cli.UUID("id", "i", &uuid, "UUID")
cli.FilePath("input", "f", &file, "Input file path (must exist)")
cli.DirPath("output", "d", &dir, "Output directory path (must exist)")
```

---

## API Reference

### CLI Struct Methods

- **`New(name, description string) *CLI`**
  Initializes a new CLI instance. Registers a default `completion` subcommand for auto-generating bash and zsh scripts.
- **`Parse(argv []string) (*SubCMD, error)`**
  Tokenizes arguments, processes validation/constraints, and returns the deepest matching subcommand leaf.
- **`ParseAndInvoke(argv []string, userdata any, preInvokeCallback func(cmd *SubCMD, userdata any)) error`**
  Parses command arguments and immediately runs the matched handler, forwarding the user data pointer and invoking the optional pre-execute callback.
- **`SubCommand(name, description string, handler func(any) error) *SubCMD`**
  Registers a top-level subcommand.
- **`ExclusiveFlags(names ...string) *CLI`**
  Defines mutually exclusive flag sets at the global/root level.
- **`PrintUsage(w io.Writer)`**
  Outputs the usage layout including global flags and all top-level subcommands.
- **`Args() []string`**
  Retrieves root positional arguments.

### SubCMD Struct Methods

- **`SubCommand(name, description string, handler func(any) error) *SubCMD`**
  Registers a nested subcommand under this node.
- **`Flag(ft flagType, name, shortName string, valuePtr any, usage string) *SubCMD`**
  Attaches a typed subcommand flag.
- **`PersistentFlag(ft flagType, name, shortName string, valuePtr any, usage string) *SubCMD`**
  Attaches an inherited typed subcommand flag.
- **`Required() *SubCMD`**
  Marks the most recently appended flag as required.
- **`Validate(validators ...FlagValidator) *SubCMD`**
  Attaches validation logic to the last added flag.
- **`ExclusiveFlags(names ...string) *SubCMD`**
  Declares mutually exclusive flags on this subcommand.
- **`Args() []string`**
  Retrieves positional arguments gathered at this subcommand's scope.
- **`Name() string`**
  Returns the subcommand name.
- **`PrintUsage(w io.Writer)`**
  Outputs subcommand details, flags, persistent flags, and inherited options.

## License

MIT License