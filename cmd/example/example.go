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

	ctx.AddFlag(goflag.FlagString, "config", "c", &config, "Path to config file", true)
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
