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

// Run the following commands to see the output
// go run cmd/examples/example.go
// go run cmd/examples/example.go -v -c config.json greet -n Abiira -g "Hello there"

func greetUser(ctx goflag.Getter, cmd goflag.Getter) {
	name := cmd.GetString("name")
	greeting := cmd.GetString("greeting")
	fmt.Println(greeting, name)

}

func printVersion(ctx goflag.Getter, cmd goflag.Getter) {
	if cmd.Get("short").(bool) {
		fmt.Println("1.0.0")
	} else {
		fmt.Println("1.0.0")
		fmt.Println("Build Date: 2021-01-01")
		fmt.Println("Commit: 1234567890")
	}
}

func handleSleep(ctx goflag.Getter, cmd goflag.Getter) {
	time.Sleep(time.Duration(cmd.GetInt("time")) * time.Second)
}

func handleCors(ctx goflag.Getter, cmd goflag.Getter) {
	fmt.Println("Origins: ", cmd.GetStringSlice("origins"))
	fmt.Println("Methods: ", cmd.GetStringSlice("methods"))
	fmt.Println("Headers: ", cmd.GetStringSlice("headers"))
	fmt.Println("Credentials: ", cmd.GetBool("credentials"))
}

func main() {
	ctx := goflag.NewContext()

	ctx.AddFlag(goflag.String("config", "c", "config.json", "Path to config file", true))
	ctx.AddFlag(goflag.Bool("verbose", "v", false, "Enable verbose output", false))
	ctx.AddFlag(goflag.Duration("timeout", "t", 5*time.Second, "Timeout for the request", true))
	ctx.AddFlag(goflag.Int("port", "p", 8080, "Port to listen on", false))
	ctx.AddFlag(goflag.Time("start", "s", time.Now(), "Start time", false))

	ctx.AddFlag(goflag.URL("url", "u", &url.URL{}, "URL to fetch", true))
	ctx.AddFlag(goflag.UUID("uuid", "i", uuid.UUID{}, "UUID to use", true))
	ctx.AddFlag(goflag.IP("ip", "i", net.IP{}, "IP to use", true))
	ctx.AddFlag(goflag.MAC("mac", "m", net.HardwareAddr{}, "MAC address to use", true))
	ctx.AddFlag(goflag.Email("email", "e", "", "Email address to use", true))
	ctx.AddFlag(goflag.HostPortPair("hostport", "h", "", "Host:Port pair to use", true))
	ctx.AddFlag(goflag.FilePath("file", "f", "", "File path to use", true))
	ctx.AddFlag(goflag.DirPath("dir", "d", "", "Directory path to use", true))

	ctx.AddSubCommand(goflag.SubCommand("greet", "Greet a person", greetUser)).
		AddFlag(goflag.String("name", "n", "World", "Name of the person to greet", true)).
		AddFlag(goflag.String("greeting", "g", "Hello", "Greeting to use", false)).
		Validate(func(a any) (bool, string) {
			if a.(string) == "World" {
				return false, "Name cannot be World"
			}
			return true, ""
		}).AddFlag(goflag.Bool("upper", "u", false, "Print in upper case", false))

	ctx.AddSubCommand(goflag.SubCommand("version", "Print version", printVersion)).
		AddFlag(goflag.Bool("verbose", "v", false, "Enable verbose output", false)).
		AddFlag(goflag.Bool("short", "s", false, "Print short version", false))

	ctx.AddSubCommand(goflag.SubCommand("sleep", "Sleep for a while", handleSleep)).
		AddFlag(goflag.Int("time", "t", 5, "Time to sleep in seconds", true))

	ctx.AddSubCommand(goflag.SubCommand("cors", "Enable CORS", handleCors)).
		AddFlag(goflag.StringSlice("origins", "o", []string{"*"}, "Allowed origins", true)).
		AddFlag(goflag.StringSlice("methods", "m", []string{"GET", "POST"}, "Allowed methods", true)).
		AddFlag(goflag.StringSlice("headers", "d", []string{"Content-Type"}, "Allowed headers", true)).
		AddFlag(goflag.Bool("credentials", "c", false, "Allow credentials", false))

	// time

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
	fmt.Println(ctx.GetDuration("timeout"))
	fmt.Println(ctx.GetInt("port"))
	fmt.Println(ctx.GetTime("start"))
	// url
	fmt.Printf("%+v\n", ctx.GetURL("url"))
	// uuid
	fmt.Println(ctx.GetUUID("uuid"))
	// ip
	fmt.Println(ctx.GetIP("ip"))
	// mac
	fmt.Println(ctx.GetMAC("mac"))
	// email
	fmt.Println(ctx.GetEmail("email"))
	// hostport
	fmt.Println(ctx.GetHostPortPair("hostport"))
	// file
	fmt.Println(ctx.GetFilePath("file"))
	// dir
	fmt.Println(ctx.GetDirPath("dir"))

}
