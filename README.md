# goflag

A simple, type-safe command-line flag parsing library for Go with support for subcommands.

## Features

- Type-safe flag parsing with dedicated methods for common types
- Support for subcommands with their own flags
- Required flags validation and custom validators for each flag.
- Method chaining for cleaner API
- Rich set of built-in types (strings, numbers, durations, URLs, IPs, emails, etc.)
- Automatic help generation
- Bash and zsh completions

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

func main() {
    var (
        config  string
        verbose bool
        timeout time.Duration
        port    int
    )
    
    cli := goflag.New()
    
    // Define flags using helper methods
    cli.String("config", "c", &config, "Path to config file").Required()
    cli.Bool("verbose", "v", &verbose, "Enable verbose output")
    cli.Duration("timeout", "t", &timeout, "Request timeout")
    cli.Int("port", "p", &port, "Port to listen on")
    
    // Parse arguments
    subcmd, err := cli.Parse(os.Args)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Config: %s\n", config)
    fmt.Printf("Verbose: %v\n", verbose)
    fmt.Printf("Timeout: %v\n", timeout)
    fmt.Printf("Port: %d\n", port)

	if subcmd != nil{
		subcmd.Handler()
	}
}
```

## Subcommands

Create subcommands with their own flags:
```go
var (
    name     string = "World"
    greeting string = "Hello"
    upper    bool
)

func greetUser() {
    message := fmt.Sprintf("%s, %s!", greeting, name)
    if upper {
        message = strings.ToUpper(message)
    }
    fmt.Println(message)
}

func main() {
    cli := goflag.New()
    
    // Define a subcommand with flags
    cli.SubCommand("greet", "Greet a person", greetUser).
        String("name", "n", &name, "Name of the person to greet").Required().
        String("greeting", "g", &greeting, "Greeting to use").
        Bool("upper", "u", &upper, "Print in upper case")
    
    subcmd, err := cli.Parse(os.Args)
    if err != nil {
        log.Fatal(err)
    }
    
    if subcmd != nil {
        subcmd.Handler()
        os.Exit(0)
    }
}
```

Usage:
```bash
$ myapp greet --name John --upper
HELLO, JOHN!

$ myapp greet -n Alice -g "Good morning"
Good morning, Alice!
```

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

### Time Types
```go
cli.Duration("timeout", "t", &duration, "Duration (e.g., 5s, 2m)")
cli.Time("start", "s", &start, "Timestamp")
```

### Slice Types
```go
cli.StringSlice("origins", "o", &origins, "Allowed origins")
cli.IntSlice("ports", "p", &ports, "Port numbers")
```

### Network Types
```go
cli.IP("address", "a", &ip, "IP address")
cli.MAC("mac", "m", &mac, "MAC address")
cli.URL("endpoint", "e", &url, "URL")
cli.HostPortPair("listen", "l", &hostport, "Host:port pair")
```

### Special Types
```go
cli.Email("contact", "c", &email, "Email address")
cli.UUID("id", "i", &uuid, "UUID")
cli.FilePath("input", "i", &file, "Input file path")
cli.DirPath("output", "o", &dir, "Output directory")
```

## Required Flags

Mark flags as required using the `.Required()` method:
```go
cli.String("config", "c", &config, "Config file").Required()
cli.Int("port", "p", &port, "Server port").Required()
```

## Complete Example
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
    config    string
    verbose   bool
    port      int
    timeout   time.Duration
    start     time.Time
    urlValue  url.URL
    uuidVal   uuid.UUID
    ipVal     net.IP
    macVal    net.HardwareAddr
    emailVal  string
    fileVal   string
    dirVal    string
    origins   []string
    methods   []string
)

func main() {
    log.SetFlags(log.Lshortfile)
    cli := goflag.New()
    
    // Global flags
    cli.String("config", "c", &config, "Path to config file").Required()
    cli.Bool("verbose", "v", &verbose, "Enable verbose output")
    cli.Duration("timeout", "t", &timeout, "Request timeout")
    cli.Int("port", "p", &port, "Port to listen on")
    cli.Time("start", "s", &start, "Start time")
    cli.URL("url", "u", &urlValue, "URL to fetch")
    cli.UUID("uuid", "i", &uuidVal, "UUID to use")
    cli.IP("ip", "", &ipVal, "IP address")
    cli.MAC("mac", "m", &macVal, "MAC address")
    cli.Email("email", "e", &emailVal, "Email address")
    cli.FilePath("file", "f", &fileVal, "File path")
    cli.DirPath("dir", "d", &dirVal, "Directory path")
    cli.StringSlice("origins", "o", &origins, "Allowed origins")
    cli.StringSlice("methods", "", &methods, "HTTP methods")
    
    // Subcommands
    cli.SubCommand("serve", "Start the server", startServer).
        Int("workers", "w", &workers, "Number of workers").Required()
    
    cli.SubCommand("version", "Print version", printVersion).
        Bool("short", "s", &short, "Short version format")
    
    // Parse
    subcmd, err := cli.Parse(os.Args)
    if err != nil {
        log.Fatal(err)
    }
    
    if subcmd != nil {
        subcmd.Handler()
        os.Exit(0)
    }
    
    // Main program logic when no subcommand is provided
    fmt.Printf("Config: %s\n", config)
    fmt.Printf("Verbose: %v\n", verbose)
    fmt.Printf("Port: %d\n", port)
}
```

## Method Chaining

Both global flags and subcommand flags support method chaining:
```go
// Global flags
cli.String("config", "c", &config, "Config file").Required()

// Subcommand flags
cli.SubCommand("deploy", "Deploy application", deployApp).
    String("env", "e", &env, "Environment").Required().
    Bool("dry-run", "d", &dryRun, "Dry run mode").
    StringSlice("tags", "t", &tags, "Deployment tags")
```

## Usage Examples
```bash
# Using long flags
$ myapp --config app.json --verbose --port 8080

# Using short flags
$ myapp -c app.json -v -p 8080

# Using subcommands
$ myapp greet --name Alice

# slice flags
$ myapp --origins http://localhost:3000,http://localhost:8080

# Help
$ myapp --help
$ myapp greet --help
```

## API Reference

### CLI Methods

- `New() *CLI` - Create a new CLI instance
- `Parse(args []string) (*Subcommand, error)` - Parse command-line arguments
- `SubCommand(name, description string, handler func()) *Subcommand` - Add a subcommand

### Flag Definition Methods

All methods follow the pattern: `Type(name, shortName string, valuePtr any, usage string) *Flag`

Available on both `*CLI` and `*Subcommand`:

- `String()` - String flag
- `Int()` - Integer flag  
- `Int64()` - 64-bit integer flag
- `Float32()` - 32-bit float flag
- `Float64()` - 64-bit float flag
- `Bool()` - Boolean flag
- `Rune()` - Single character flag
- `Duration()` - time.Duration flag
- `Time()` - time.Time flag
- `StringSlice()` - String slice flag
- `IntSlice()` - Integer slice flag
- `IP()` - IP address flag
- `MAC()` - MAC address flag
- `URL()` - URL flag
- `UUID()` - UUID flag
- `HostPortPair()` - Host:port pair flag
- `Email()` - Email address flag
- `FilePath()` - File path flag
- `DirPath()` - Directory path flag

### Flag Methods

- `Required()` - Mark flag as required

### Subcommand Methods

- `Handler()` - Execute the subcommand handler

## License

MIT License