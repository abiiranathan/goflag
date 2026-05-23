package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/abiiranathan/goflag"
	"github.com/google/uuid"
)

// AppConfig is an example of arbitrary userdata threaded through ParseAndInvoke.
type AppConfig struct {
	DBConnStr string
	Debug     bool
}

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

	origins     []string = []string{"*"}
	methods     []string = []string{"GET", "POST"}
	headers     []string = []string{"Content-Type"}
	credentials bool
	verbose     bool
	config      string        = "config.json"
	port        int           = 8080
	start       time.Time     = time.Now()
	timeout     time.Duration = 5 * time.Second

	durationValue time.Duration = 5 * time.Second
	upperValue    bool

	// Nested subcommand flags — "greet user" and "greet bot".
	userName string
	userAge  int
	botName  string
	botDelay time.Duration

	// Persistent flag shared across all "deploy" sub-subcommands.
	deployEnv string

	// Mutually exclusive export flags.
	exportJSON bool
	exportYAML bool
)

func greetUser(userdata any) error {
	cfg, _ := userdata.(*AppConfig)
	if cfg != nil && cfg.Debug {
		fmt.Println("[debug] greet user handler running")
	}
	fmt.Println(greeting, name)
	return nil
}

func greetUserSub(userdata any) error {
	cfg, _ := userdata.(*AppConfig)
	if cfg != nil && cfg.Debug {
		fmt.Printf("[debug] DB: %s\n", cfg.DBConnStr)
	}
	fmt.Printf("Greeting user: %s (age %d)\n", userName, userAge)
	return nil
}

func greetBotSub(userdata any) error {
	fmt.Printf("Greeting bot: %s (delay: %v)\n", botName, botDelay)
	return nil
}

func printVersion(userdata any) error {
	if short {
		fmt.Println("1.0.0")
	} else {
		fmt.Println("1.0.0")
		fmt.Println("Build Date: 2021-01-01")
		fmt.Println("Commit: 1234567890")
	}
	return nil
}

func handleSleep(userdata any) error {
	fmt.Printf("Sleeping for %v...\n", durationValue)
	time.Sleep(durationValue)
	return nil
}

func handleCors(userdata any) error {
	fmt.Println("Origins:", origins)
	fmt.Println("Methods:", methods)
	fmt.Println("Headers:", headers)
	fmt.Println("Credentials:", credentials)
	return nil
}

// handleDeployStaging demonstrates reading a persistent flag (deployEnv) that
// was declared on the parent "deploy" subcommand.
func handleDeployStaging(userdata any) error {
	if deployEnv == "" {
		// Feature 2 in action: return a real error instead of calling log.Fatal.
		return errors.New("--env is required for staging deployments")
	}
	fmt.Printf("Deploying to staging (env=%s)\n", deployEnv)
	return nil
}

func handleDeployProd(userdata any) error {
	fmt.Printf("Deploying to production (env=%s)\n", deployEnv)
	return nil
}

func handleExport(userdata any) error {
	// Feature 3 makes the CLI reject --json + --yaml together; by the time
	// this handler runs, exactly one (or neither) is set.
	switch {
	case exportJSON:
		fmt.Println("Exporting as JSON")
	case exportYAML:
		fmt.Println("Exporting as YAML")
	default:
		fmt.Println("Exporting in default format")
	}
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	cli := goflag.New("MyApp", "An example application using goflag")

	// Global flags.
	cli.String("config", "c", &config, "Path to config file")
	cli.Bool("verbose", "v", &verbose, "Enable verbose output")
	cli.Duration("timeout", "t", &timeout, "Timeout for the request")
	cli.Int("port", "p", &port, "Port to listen on")
	cli.HostPortPair("hostport", "h", &hpVal, "Host:Port to listen on")
	cli.Time("start", "s", &start, "Start time")
	cli.URL("url", "u", &urlValue, "URL to fetch")
	cli.UUID("uuid", "i", &uuidVal, "UUID to use")
	cli.IP("ip", "i", &ipVal, "IP to use")
	cli.MAC("mac", "m", &macVal, "MAC address to use")
	cli.Email("email", "e", &emailVal, "Email address to use")
	cli.FilePath("file", "f", &fileVal, "File path to use")
	cli.DirPath("dir", "d", &dirVal, "Directory path to use")

	// "greet" subcommand with nested "user" and "bot" sub-subcommands.
	greetCmd := cli.SubCommand("greet", "Greet a person", greetUser).
		String("name", "n", &name, "Name of the person to greet").Required().
		String("greeting", "g", &greeting, "Greeting to use").
		Bool("upper", "u", &upperValue, "Print in upper case")

	// Nested: greet user --name Alice --age 30
	greetCmd.SubCommand("user", "Greet a specific user with age", greetUserSub).
		String("name", "n", &userName, "User name").Required().
		Int("age", "a", &userAge, "User age")

	// Nested: greet bot --name R2D2 --delay 500ms
	greetCmd.SubCommand("bot", "Greet a bot", greetBotSub).
		String("name", "n", &botName, "Bot name").Required().
		Duration("delay", "d", &botDelay, "Response delay")

	cli.SubCommand("version", "Print version", printVersion).
		Bool("verbose", "v", &verbose, "Enable verbose output").
		Bool("short", "s", &short, "Print short version")

	cli.SubCommand("sleep", "Sleep for a while", handleSleep).
		Duration("time", "t", &durationValue, "Time to sleep").Required()

	cli.SubCommand("cors", "Enable CORS", handleCors).
		StringSlice("origins", "o", &origins, "Allowed origins").Required().
		StringSlice("methods", "m", &methods, "Allowed methods").Required().
		StringSlice("headers", "d", &headers, "Allowed headers").Required().
		Bool("credentials", "c", &credentials, "Allow credentials")

	// Persistent flags: --env is declared once on "deploy" and
	// automatically available to both "staging" and "prod" sub-subcommands.
	deployCmd := cli.SubCommand("deploy", "Deploy the application", func(userdata any) error {
		fmt.Println("specify a sub-subcommand: staging or prod")
		return nil
	})
	deployCmd.PersistentString("env", "e", &deployEnv, "Target environment (e.g. production)")

	deployCmd.SubCommand("staging", "Deploy to staging", handleDeployStaging)
	deployCmd.SubCommand("prod", "Deploy to production", handleDeployProd)

	// Mutually exclusive flags: --json and --yaml cannot both be
	// provided. The parser enforces this before the handler is called.
	cli.SubCommand("export", "Export data", handleExport).
		Bool("json", "j", &exportJSON, "Output as JSON").
		Bool("yaml", "y", &exportYAML, "Output as YAML").
		ExclusiveFlags("json", "yaml")

	// Thread application config as userdata.
	appCfg := &AppConfig{
		DBConnStr: "postgres://localhost/mydb",
		Debug:     verbose,
	}

	// Feature 2 — ParseAndInvoke now propagates handler errors.
	err := cli.ParseAndInvoke(os.Args, appCfg, func(cmd *goflag.SubCMD, userdata any) {
		if cmd == nil {
			return
		}
		if args := cmd.Args(); len(args) > 0 {
			fmt.Println("Positional args:", args)
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Reached only when no subcommand was matched.
	if args := cli.Args(); len(args) > 0 {
		fmt.Println("Root positional args:", args)
	}

	fmt.Println("Config:", config)
	fmt.Println("Verbose:", verbose)
	fmt.Println("Timeout:", timeout)
	fmt.Println("Port:", port)
	fmt.Println("Start:", start)
	fmt.Println("URL:", urlValue)
	fmt.Println("UUID:", uuidVal)
	fmt.Println("IP:", ipVal)
	fmt.Println("MAC:", macVal)
	fmt.Println("Email:", emailVal)
	fmt.Println("HostPort:", hpVal)
	fmt.Println("File:", fileVal)
	fmt.Println("Dir:", dirVal)
	fmt.Println("Origins:", origins)
	fmt.Println("Methods:", methods)
	fmt.Println("Headers:", headers)
	fmt.Println("Credentials:", credentials)
}
