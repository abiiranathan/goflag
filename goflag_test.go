package goflag

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cli := New()
	var (
		name   string = "World"
		age    int    = 20
		height int    = 0
	)

	cli.subcommands = []*subcommand{
		{
			name:        "test",
			description: "Test command",
			flags: []*Flag{
				{
					name:      "name",
					flagType:  flagString,
					value:     &name,
					shortName: "n",
					required:  true,
					usage:     "Your name",
					validators: []FlagValidator{
						func(value any) (valid bool, errmsg string) {
							return value != "", "Name cannot be empty"
						},
					},
				},
				{
					name:      "age",
					flagType:  flagInt,
					value:     &age,
					shortName: "a",
					required:  true,
					usage:     "Your age",
				},
				{
					name:      "height",
					flagType:  flagInt,
					value:     &height,
					shortName: "l",
					required:  true,
					usage:     "Your height",
					validators: []FlagValidator{
						func(a any) (bool, string) {
							return a.(int) > 0, "Height must be greater than 0"
						},
					},
				},
			},
			Handler: func() {
				fmt.Printf("Hello %s, you are %d years old and %d cm tall\n", name, age, height)
			},
		},
		{
			name:        "test2",
			description: "Test command 2",
			Handler: func() {
				fmt.Printf("Hello %s, you are %d years old and %d cm tall\n", name, age, height)
			},
			flags: []*Flag{
				{
					name:      "age",
					flagType:  flagInt,
					value:     &age,
					shortName: "a",
					required:  true,
					usage:     "Your age",
				},
				{
					name:      "height",
					flagType:  flagInt,
					value:     &height,
					shortName: "l",
					required:  true,
					usage:     "Your height",
					validators: []FlagValidator{
						func(a any) (bool, string) {
							return a.(int) > 0, "Height must be greater than 0"
						},
					},
				},
			},
		},
	}

	// create a test argv
	argv1 := []string{
		"myapp",
		"test",
		"--name", "John",
		"--age", "30",
		"--height", "100",
	}

	s1, err := cli.Parse(argv1)
	if err != nil {
		t.Fatal("Expected no error: ", err)
	}

	fmt.Printf("subcommand: %v\n", s1)

	if s1 == nil {
		t.Fatalf("Expected subcommand not to be nil, but got nil")
	}

	if s1.name != "test" {
		t.Errorf("Expected subcommand name to be 'test', but got '%v'", s1.name)
	}

	// Height is expected to fail and panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected height flag to panic, but it didn't")
		}
	}()

	argv2 := []string{
		"myapp",
		"test2",
		"--age", "30",
		"--height", "0",
	}

	s2, err := cli.Parse(argv2)
	if err == nil {
		t.Fatalf("Expected error, but got nil\n")
	}

	if s2 != nil {
		t.Fatalf("Expected subcommand to be not nil, but got nil")
	}

	// check the flag values
	if cli.subcommands[0].flags[0].value.(string) != "John" {
		t.Errorf("Expected name flag value to be 'John', but got '%v'", cli.subcommands[0].flags[0].value)
	}

	if cli.subcommands[1].flags[0].value.(int) != 30 {
		t.Errorf("Expected age flag value to be 30, but got '%v'", cli.subcommands[1].flags[0].value)
	}

}

func TestAddCommand(t *testing.T) {
	cli := New()
	var name string
	cli.SubCommand("test", "Test command", func() {
		fmt.Println("Test command")
	}).Flag(flagString, "name", "n", &name, "Your name").Required()

	// cli.subCommands[0] is completions automatically inserted with NewContext.
	// subcommand[0].flags[0] is the --help flag auto inserted by .AddFlag.
	if len(cli.subcommands) != 2 {
		t.Errorf("Expected 2 subcommand, but got %v", len(cli.subcommands))
	}

	if cli.subcommands[1].name != "test" {
		t.Errorf("Expected subcommand name to be 'test', but got '%v'", cli.subcommands[0].name)
	}

	if cli.subcommands[1].description != "Test command" {
		t.Errorf("Expected subcommand usage to be 'Test command', but got '%v'", cli.subcommands[0].description)
	}

	if len(cli.subcommands[1].flags) != 2 {
		t.Errorf("Expected 2 flags, but got %v", len(cli.subcommands[0].flags))
	}

	if cli.subcommands[1].flags[1].name != "name" {
		t.Errorf("Expected flag name to be 'name', but got '%v'", cli.subcommands[1].flags[1].name)
	}
}

func TestPrintUsage(t *testing.T) {
	// Set up test data
	cli := &CLI{
		flags: []*Flag{
			{
				name:  "help",
				usage: "Show help message",
			},
			{
				name:      "name",
				shortName: "n",
				usage:     "Name of the person",
			},
			{
				name:      "age",
				shortName: "a",
				usage:     "Age of the person",
			},
		},
		subcommands: []*subcommand{
			{
				name:        "add",
				description: "Add a new person",
				flags: []*Flag{
					{
						name:  "help",
						usage: "Show help message",
					},
					{
						name:      "name",
						shortName: "n",
						usage:     "Name of the person",
					},
					{
						name:      "age",
						shortName: "a",
						usage:     "Age of the person",
					},
				},
			},
			{
				name:        "delete",
				description: "Delete an existing person",
				flags: []*Flag{
					{
						name:  "help",
						usage: "Show help message",
					},
					{
						name:      "name",
						shortName: "n",
						usage:     "Name of the person",
					},
				},
			},
		},
	}

	// Capture output
	var buf bytes.Buffer
	cli.PrintUsage(&buf)

	// Check output
	expectedInOutput := []string{
		"--help",
		"-n",
		"--name",
		"-a",
		"--age",
		"delete",
		"add",
		"Global Flags:",
		"Subcommands:",
		"Usage:",
		"Flags:",
		"Name of the person",
		"Age of the person",
		"Show help message",
		"Delete an existing person",
		"Add a new person",
	}

	for _, expected := range expectedInOutput {
		if !strings.Contains(buf.String(), expected) {
			t.Errorf("Expected output to contain '%v'", expected)
		}
	}

}

func TestAddFlag(t *testing.T) {
	cli := New()
	var name string

	cli.addFlag(flagString, "name", "n", &name, "Your name").Required()

	if len(cli.flags) != 2 {
		t.Errorf("Expected 2 flags, but got %v", len(cli.flags))
	}

	if cli.flags[1].name != "name" {
		t.Errorf("Expected flag name to be 'name', but got '%v'", cli.flags[1].name)
	}

	if cli.flags[1].shortName != "n" {
		t.Errorf("Expected flag short name to be 'n', but got '%v'", cli.flags[1].shortName)
	}
}

func TestGlobalRequiredFlags(t *testing.T) {
	// test global required flag
	cli := New()

	var verbose bool
	var port int

	cli.flags = []*Flag{
		{
			name:      "verbose",
			flagType:  flagBool,
			value:     &verbose,
			shortName: "v",
			usage:     "Enable verbose output",
			required:  true,
		},
		{
			name:      "port",
			flagType:  flagInt,
			value:     &port,
			shortName: "p",
			usage:     "Port to listen on",
			required:  true,
		},
	}

	// test cases
	argv1 := []string{
		"myapp",
		"--verbose",
		"--port", "8080",
	}

	if _, err := cli.Parse(argv1); err != nil {
		t.Fatalf("Expected no error, but got '%v'", err)
	}

}
