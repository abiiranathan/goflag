package goflag

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	ctx := NewContext()
	ctx.subcommands = []*subcommand{
		{
			name: "test",
			flags: []*gflag{
				{
					name:      "name",
					flagType:  flagString,
					value:     "",
					shortName: "n",
					required:  true,
					usage:     "Your name",
					validator: func(value any) (bool, string) {
						return value != "", "Name cannot be empty"
					},
				},
			},
		},
		{
			name: "test2",
			flags: []*gflag{
				{
					name:      "age",
					flagType:  flagInt,
					value:     20,
					shortName: "a",
					required:  true,
					usage:     "Your age",
					validator: func(value any) (bool, string) {
						return value.(int) > 0, "Age must be greater than 0"
					},
				},

				{
					name:      "height",
					flagType:  flagFloat64,
					value:     0.0,
					shortName: "h",
					required:  true,
					usage:     "Your height",
					validator: func(value any) (bool, string) {
						return value.(float64) > 0.0, "Height must be greater than 0"
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
	}

	argv2 := []string{
		"myapp",
		"test2",
		"--age", "30",
		"--height", "0",
	}

	s1, err := ctx.Parse(argv1)
	if err != nil {
		t.Fatal("Expected validation error for height, got nil\n", err)
	}

	if s1 == nil {
		t.Fatalf("Expected subcommand to be not nil, but got nil")
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
	if ctx.subcommands[0].flags[0].value != "John" {
		t.Errorf("Expected name flag value to be 'John', but got '%v'", ctx.subcommands[0].flags[0].value)
	}

	if ctx.subcommands[1].flags[0].value != 30 {
		t.Errorf("Expected age flag value to be 30, but got '%v'", ctx.subcommands[1].flags[0].value)
	}

	// Test the error cases and global flags
	ctx.flags = []*gflag{
		{
			name:      "verbose",
			flagType:  flagBool,
			value:     false,
			shortName: "v",
			required:  false,
			usage:     "Verbose mode",
			validator: nil,
		},
		{
			name:      "port",
			flagType:  flagInt,
			value:     8080,
			shortName: "p",
			required:  false,
			usage:     "Port number",
		},
		{
			name:      "host",
			flagType:  flagString,
			value:     "localhost",
			shortName: "h",
			required:  false,
			usage:     "Host name",
		},
		{
			name:      "timeout",
			flagType:  flagFloat64,
			value:     10.0,
			shortName: "t",
			required:  false,
			usage:     "Timeout in seconds",
		},
	}

	ctx.subcommands = []*subcommand{
		{
			name: "run",
			flags: []*gflag{
				{
					name:      "name",
					flagType:  flagString,
					value:     "",
					shortName: "n",
					required:  true,
					usage:     "Your name",
					validator: func(value any) (bool, string) {
						return value != "", "Name cannot be empty"
					},
				},
			},
		},
	}

	// test global flags
	argv3 := []string{
		"myapp",
		"--verbose",
		"--port", "8081",
		"--host", "localhost",
		"--timeout", "10.5",
		"run",
		"--name", "John",
	}

	s3, err := ctx.Parse(argv3)
	if err != nil {
		t.Fatalf("Expected no error, but got '%v'", err)
	}
	if s3 == nil {
		t.Fatalf("Expected subcommand to be not nil, but got nil")
	}

	if s3.name != "run" {
		t.Errorf("Expected subcommand name to be 'run', but got '%v'", s3.name)
	}

	if ctx.flags[0].value != true {
		t.Errorf("Expected verbose flag value to be true, but got '%v'", ctx.flags[0].value)
	}

	if ctx.flags[1].value != 8081 {
		t.Errorf("Expected port flag value to be 8081, but got '%v'", ctx.flags[1].value)
	}

	if ctx.flags[2].value != "localhost" {
		t.Errorf("Expected host flag value to be 'localhost', but got '%v'", ctx.flags[2].value)
	}

	if ctx.flags[3].value != 10.5 {
		t.Errorf("Expected timeout flag value to be 10.5, but got '%v'", ctx.flags[3].value)
	}

	if ctx.subcommands[0].flags[0].value != "John" {
		t.Errorf("Expected name flag value to be 'John', but got '%v'", ctx.subcommands[0].flags[0].value)
	}

	// test error cases
	argv4 := []string{
		"myapp",
		"--verbose",
		"--port", "8081",
	}

	s4, err := ctx.Parse(argv4)
	if err != nil {
		t.Fatalf("Expected no error, but got '%v'", err)
	}

	if s4 != nil {
		t.Fatalf("Expected subcommand to be nil, but got '%v'", s4)
	}

	// Test bool with no value
	argv5 := []string{
		"myapp",
		"--verbose",
		"--port", "8081",
		"run",
		"--name", "John",
	}

	s5, err := ctx.Parse(argv5)
	if err != nil {
		t.Fatalf("Expected no error, but got '%v'", err)
	}

	if s5 == nil {
		t.Fatalf("Expected subcommand to be not nil, but got nil")
	}

	if s5.name != "run" {
		t.Errorf("Expected subcommand name to be 'run', but got '%v'", s5.name)
	}

	if ctx.flags[0].value != true {
		t.Errorf("Expected verbose flag value to be true, but got '%v'", ctx.flags[0].value)
	}

	// Test that a regular flag with no value panics
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected name flag to panic, but it didn't")
		}
	}()

	ctx = NewContext()
	ctx.flags = []*gflag{
		{
			name:      "query",
			flagType:  flagString,
			value:     "",
			shortName: "q",
			required:  true,
			usage:     "Query string",
		},
	}

	argv6 := []string{
		"myapp",
		"-query",
	}

	// Should panic
	ctx.Parse(argv6)

}

func TestAddCommand(t *testing.T) {
	ctx := NewContext()
	ctx.AddSubCommand(&subcommand{
		name:        "test",
		description: "Test command",
		Handler: func(ctx Getter, subCommand Getter) {
			fmt.Println("Test command")
		},
		flags: []*gflag{
			{
				name:      "name",
				flagType:  flagString,
				value:     "",
				shortName: "n",
				required:  true,
				usage:     "Your name",
				validator: func(value any) (bool, string) {
					return value != "", "Name cannot be empty"
				},
			},
		},
	})

	if len(ctx.subcommands) != 1 {
		t.Errorf("Expected 1 subcommand, but got %v", len(ctx.subcommands))
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

	if ctx.subcommands[0].flags[0].name != "name" {
		t.Errorf("Expected flag name to be 'name', but got '%v'", ctx.subcommands[0].flags[0].name)
	}
}

func TestPrintUsage(t *testing.T) {
	// Set up test data
	ctx := &Context{
		flags: []*gflag{
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
				flags: []*gflag{
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
				flags: []*gflag{
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
	ctx.AddFlag(&gflag{
		name:      "name",
		flagType:  flagString,
		value:     "",
		shortName: "n",
		required:  true,
		usage:     "Your name",
		validator: func(value any) (bool, string) {
			return value != "", "Name cannot be empty"
		},
	})

	if len(ctx.flags) != 2 {
		t.Errorf("Expected 2 flag, but got %v", len(ctx.flags))
	}

	if ctx.flags[1].name != "name" {
		t.Errorf("Expected flag name to be 'name', but got '%v'", ctx.flags[0].name)
	}

	if ctx.flags[1].usage != "Your name" {
		t.Errorf("Expected flag usage to be 'Your name', but got '%v'", ctx.flags[0].usage)
	}
}

func TestGlobalRequiredFlags(t *testing.T) {
	// test global required flag
	ctx := NewContext()

	ctx.flags = []*gflag{
		{
			name:      "name",
			flagType:  flagString,
			value:     "",
			shortName: "n",
			required:  true,
			usage:     "Your name",
		},
	}

	// we need to pass some argv, otherwise ctx.Parse will return nil
	// and not panic as we expect
	argv := []string{"myapp", "--random", "value"}
	_, err := ctx.Parse(argv)
	fmt.Println(err)
	if err == nil {
		t.Errorf("Expected error not to be nil due to failed validation but got: %v", err)
	}
}

func TestSubcommandRequiredFlags(t *testing.T) {
	// test global required flag
	ctx := NewContext()

	ctx.subcommands = []*subcommand{
		{
			name: "test",
			flags: []*gflag{
				{
					name:      "name",
					flagType:  flagString,
					value:     "",
					shortName: "n",
					required:  true,
					usage:     "Your name",
				},
			},
		},
	}

}
