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
	durationValue time.Duration = 5 * time.Second

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
	cli := goflag.New()

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

	cli.SubCommand("greet", "Greet a person", greetUser).
		String("name", "n", &name, "Name of the person to greet").Required().
		String("greeting", "g", &greeting, "Greeting to use").
		Bool("upper", "u", &upperValue, "Print in upper case")

	cli.SubCommand("version", "Print version", printVersion).
		Bool("verbose", "v", &verbose, "Enable verbose output").
		Bool("short", "s", &short, "Print short version")

	cli.SubCommand("sleep", "Sleep for a while", handleSleep).
		Duration("time", "t", &durationValue, "Time to sleep in seconds").Required()

	cli.SubCommand("cors", "Enable CORS", handleCors).
		StringSlice("origins", "o", &origins, "Allowed origins").Required().
		StringSlice("methods", "m", &methods, "Allowed methods").Required().
		StringSlice("headers", "d", &headers, "Allowed headers").Required().
		Bool("credentials", "c", &credentials, "Allow credentials")

	// Parse the command line arguments and return the matching subcommand
	subcmd, err := cli.Parse(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

	if subcmd != nil {
		subcmd.Handler()
		os.Exit(0)
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
