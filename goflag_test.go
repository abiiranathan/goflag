package goflag

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	ctx := NewContext()
	var (
		name   string = "World"
		age    int    = 20
		height int    = 0
	)

	ctx.subcommands = []*subcommand{
		{
			name:        "test",
			description: "Test command",
			flags: []*Flag{
				{
					Name:      "name",
					FlagType:  FlagString,
					Value:     &name,
					ShortName: "n",
					Required:  true,
					Usage:     "Your name",
					Validator: func(value any) (bool, string) {
						return value != "", "Name cannot be empty"
					},
				},
				{
					Name:      "age",
					FlagType:  FlagInt,
					Value:     &age,
					ShortName: "a",
					Required:  true,
					Usage:     "Your age",
				},
				{
					Name:      "height",
					FlagType:  FlagInt,
					Value:     &height,
					ShortName: "l",
					Required:  true,
					Usage:     "Your height",
					Validator: func(a any) (bool, string) {
						return a.(int) > 0, "Height must be greater than 0"
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
					Name:      "age",
					FlagType:  FlagInt,
					Value:     &age,
					ShortName: "a",
					Required:  true,
					Usage:     "Your age",
				},
				{
					Name:      "height",
					FlagType:  FlagInt,
					Value:     &height,
					ShortName: "l",
					Required:  true,
					Usage:     "Your height",
					Validator: func(a any) (bool, string) {
						return a.(int) > 0, "Height must be greater than 0"
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

	s1, err := ctx.Parse(argv1)
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

	s2, err := ctx.Parse(argv2)
	if err == nil {
		t.Fatalf("Expected error, but got '%v'", err)
	}

	if s2 != nil {
		t.Fatalf("Expected subcommand to be not nil, but got nil")
	}

	if s2.name != "test2" {
		t.Errorf("Expected subcommand name to be 'test2', but got '%v'", s2.name)
	}

	// check the flag values
	if ctx.subcommands[0].flags[0].Value != "John" {
		t.Errorf("Expected name flag value to be 'John', but got '%v'", ctx.subcommands[0].flags[0].Value)
	}

	if ctx.subcommands[1].flags[0].Value != 30 {
		t.Errorf("Expected age flag value to be 30, but got '%v'", ctx.subcommands[1].flags[0].Value)
	}

}

func TestAddCommand(t *testing.T) {
	ctx := NewContext()
	var name string
	ctx.AddSubCommand("test", "Test command", func() {
		fmt.Println("Test command")
	}).AddFlag(FlagString, "name", "n", &name, "Your name", true)

	// ctx.subCommands[0] is completions automatically inserted with NewContext.
	// subcommand[0].flags[0] is the --help flag auto inserted by .AddFlag.
	if len(ctx.subcommands) != 1 {
		t.Errorf("Expected 2 subcommand, but got %v", len(ctx.subcommands))
	}

	if ctx.subcommands[0].name != "test" {
		t.Errorf("Expected subcommand name to be 'test', but got '%v'", ctx.subcommands[0].name)
	}

	if ctx.subcommands[0].description != "Test command" {
		t.Errorf("Expected subcommand usage to be 'Test command', but got '%v'", ctx.subcommands[0].description)
	}

	if len(ctx.subcommands[0].flags) != 2 {
		t.Errorf("Expected 2 flags, but got %v", len(ctx.subcommands[0].flags))
	}

	if ctx.subcommands[0].flags[1].Name != "name" {
		t.Errorf("Expected flag name to be 'name', but got '%v'", ctx.subcommands[1].flags[1].Name)
	}
}

func TestPrintUsage(t *testing.T) {
	// Set up test data
	ctx := &Context{
		flags: []*Flag{
			{
				Name:  "help",
				Usage: "Show help message",
			},
			{
				Name:      "name",
				ShortName: "n",
				Usage:     "Name of the person",
			},
			{
				Name:      "age",
				ShortName: "a",
				Usage:     "Age of the person",
			},
		},
		subcommands: []*subcommand{
			{
				name:        "add",
				description: "Add a new person",
				flags: []*Flag{
					{
						Name:  "help",
						Usage: "Show help message",
					},
					{
						Name:      "name",
						ShortName: "n",
						Usage:     "Name of the person",
					},
					{
						Name:      "age",
						ShortName: "a",
						Usage:     "Age of the person",
					},
				},
			},
			{
				name:        "delete",
				description: "Delete an existing person",
				flags: []*Flag{
					{
						Name:  "help",
						Usage: "Show help message",
					},
					{
						Name:      "name",
						ShortName: "n",
						Usage:     "Name of the person",
					},
				},
			},
		},
	}

	// Capture output
	var buf bytes.Buffer
	ctx.PrintUsage(&buf)

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
	ctx := NewContext()
	var name string
	ctx.AddFlag(FlagString, "name", "n", &name, "Your name", true)

	if len(ctx.flags) != 2 {
		t.Errorf("Expected 2 flags, but got %v", len(ctx.flags))
	}

	if ctx.flags[1].Name != "name" {
		t.Errorf("Expected flag name to be 'name', but got '%v'", ctx.flags[1].Name)
	}

	if ctx.flags[1].ShortName != "n" {
		t.Errorf("Expected flag short name to be 'n', but got '%v'", ctx.flags[1].ShortName)
	}
}

func TestGlobalRequiredFlags(t *testing.T) {
	// test global required flag
	ctx := NewContext()

	var verbose bool
	var port int

	ctx.flags = []*Flag{
		{
			Name:      "verbose",
			FlagType:  FlagBool,
			Value:     &verbose,
			ShortName: "v",
			Usage:     "Enable verbose output",
			Required:  true,
		},
		{
			Name:      "port",
			FlagType:  FlagInt,
			Value:     &port,
			ShortName: "p",
			Usage:     "Port to listen on",
			Required:  true,
		},
	}

	// test cases
	argv1 := []string{
		"myapp",
		"--verbose",
		"--port", "8080",
	}

	if _, err := ctx.Parse(argv1); err != nil {
		t.Fatalf("Expected no error, but got '%v'", err)
	}

}
