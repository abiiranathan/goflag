package goflag

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cli := New("test", "My test app")
	var (
		name   = "World"
		age    = 20
		height = 0
	)

	cli.subcommands = []*SubCMD{
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
							return value != "", "name cannot be empty"
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
							return a.(int) > 0, "height must be greater than 0"
						},
					},
				},
			},
			Handler: func(userdata any) {
				fmt.Printf("Hello %s, you are %d years old and %d cm tall\n", name, age, height)
			},
		},
		{
			name:        "test2",
			description: "Test command 2",
			Handler: func(userdata any) {
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
							return a.(int) > 0, "height must be greater than 0"
						},
					},
				},
			},
		},
	}

	argv1 := []string{
		"myapp",
		"test",
		"--name", "John",
		"--age", "30",
		"--height", "100",
	}

	s1, err := cli.Parse(argv1)
	if err != nil {
		t.Fatal("expected no error:", err)
	}
	if s1 == nil {
		t.Fatal("expected subcommand, got nil")
	}
	if s1.name != "test" {
		t.Errorf("expected subcommand name 'test', got %q", s1.name)
	}

	// height=0 triggers the validator; expect an error.
	argv2 := []string{
		"myapp",
		"test2",
		"--age", "30",
		"--height", "0",
	}

	s2, err := cli.Parse(argv2)
	if err == nil {
		t.Fatal("expected error for height=0, got nil")
	}
	if s2 != nil {
		t.Fatal("expected nil subcommand on error, got non-nil")
	}
}

func TestAddCommand(t *testing.T) {
	cli := New("test", "My test app")
	var name string
	cli.SubCommand("test", "Test command", func(userdata any) {
		fmt.Println("Test command")
	}).Flag(flagString, "name", "n", &name, "Your name").Required()

	// cli.subcommands[0] is the auto-registered completion subcommand.
	if len(cli.subcommands) != 2 {
		t.Errorf("expected 2 subcommands, got %d", len(cli.subcommands))
	}
	if cli.subcommands[1].name != "test" {
		t.Errorf("expected subcommand name 'test', got %q", cli.subcommands[1].name)
	}
	if cli.subcommands[1].description != "Test command" {
		t.Errorf("expected description 'Test command', got %q", cli.subcommands[1].description)
	}
	// 2 flags: auto-inserted --help + "name".
	if len(cli.subcommands[1].flags) != 2 {
		t.Errorf("expected 2 flags, got %d", len(cli.subcommands[1].flags))
	}
	if cli.subcommands[1].flags[1].name != "name" {
		t.Errorf("expected flag name 'name', got %q", cli.subcommands[1].flags[1].name)
	}
}

func TestPrintUsage(t *testing.T) {
	cli := &CLI{
		flags: []*Flag{
			{name: "help", usage: "Show help message"},
			{name: "name", shortName: "n", usage: "Name of the person"},
			{name: "age", shortName: "a", usage: "Age of the person"},
		},
		subcommands: []*SubCMD{
			{
				name:        "add",
				description: "Add a new person",
				flags: []*Flag{
					{name: "help", usage: "Show help message"},
					{name: "name", shortName: "n", usage: "Name of the person"},
					{name: "age", shortName: "a", usage: "Age of the person"},
				},
			},
			{
				name:        "delete",
				description: "Delete an existing person",
				flags: []*Flag{
					{name: "help", usage: "Show help message"},
					{name: "name", shortName: "n", usage: "Name of the person"},
				},
			},
		},
	}

	var buf bytes.Buffer
	cli.PrintUsage(&buf)

	expected := []string{
		"--help", "-n", "--name", "-a", "--age",
		"delete", "add", "Global Flags:", "Subcommands:", "Usage:",
		"Name of the person", "Age of the person", "Show help message",
		"Delete an existing person", "Add a new person",
	}
	for _, e := range expected {
		if !strings.Contains(buf.String(), e) {
			t.Errorf("expected output to contain %q", e)
		}
	}
}

func TestAddFlag(t *testing.T) {
	cli := New("test", "My test app")
	var name string
	cli.addFlag(flagString, "name", "n", &name, "Your name").Required()

	if len(cli.flags) != 2 {
		t.Errorf("expected 2 flags (help + name), got %d", len(cli.flags))
	}
	if cli.flags[1].name != "name" {
		t.Errorf("expected flag name 'name', got %q", cli.flags[1].name)
	}
	if cli.flags[1].shortName != "n" {
		t.Errorf("expected short name 'n', got %q", cli.flags[1].shortName)
	}
}

func TestGlobalRequiredFlags(t *testing.T) {
	cli := New("test", "My test app")
	var verbose bool
	var port int

	cli.flags = []*Flag{
		{name: "verbose", flagType: flagBool, value: &verbose, shortName: "v", usage: "Enable verbose output", required: true},
		{name: "port", flagType: flagInt, value: &port, shortName: "p", usage: "Port to listen on", required: true},
	}

	argv1 := []string{"myapp", "--verbose", "--port", "8080"}
	if _, err := cli.Parse(argv1); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNestedSubCommands(t *testing.T) {
	cli := New("test", "Nested subcommand test")

	var output string
	var count int
	var dryRun bool

	deployCmd := cli.SubCommand("deploy", "Deploy something", func(userdata any) {})

	deployCmd.SubCommand("staging", "Deploy to staging", func(userdata any) {
		output = fmt.Sprintf("staging count=%d dry=%v", count, dryRun)
	}).
		Int("count", "c", &count, "Number of replicas").Required().
		Bool("dry-run", "d", &dryRun, "Dry run mode")

	deployCmd.SubCommand("prod", "Deploy to production", func(userdata any) {
		output = fmt.Sprintf("prod count=%d dry=%v", count, dryRun)
	}).
		Int("count", "c", &count, "Number of replicas").Required()

	argv := []string{"test", "deploy", "staging", "--count", "3", "--dry-run"}
	leaf, err := cli.Parse(argv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if leaf == nil {
		t.Fatal("expected a leaf subcommand, got nil")
	}
	if leaf.Name() != "staging" {
		t.Errorf("expected leaf name 'staging', got %q", leaf.Name())
	}

	leaf.Handler(nil)
	if output != "staging count=3 dry=true" {
		t.Errorf("unexpected handler output: %q", output)
	}
}

func TestPositionalArgs(t *testing.T) {
	cli := New("test", "Positional args test")

	var format string
	var captured []string

	cli.SubCommand("convert", "Convert files", func(userdata any) {
		captured = userdata.([]string)
	}).String("format", "f", &format, "Output format")

	// Everything after --format value should land in Args().
	argv := []string{"test", "convert", "--format", "png", "file1.jpg", "file2.jpg", "file3.jpg"}
	leaf, err := cli.Parse(argv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if leaf == nil {
		t.Fatal("expected subcommand, got nil")
	}

	positional := leaf.Args()
	if len(positional) != 3 {
		t.Fatalf("expected 3 positional args, got %d: %v", len(positional), positional)
	}

	// Invoke handler with the positional args as userdata.
	leaf.Handler(positional)
	if len(captured) != 3 || captured[0] != "file1.jpg" {
		t.Errorf("unexpected captured args: %v", captured)
	}
}

func TestRootPositionalArgs(t *testing.T) {
	cli := New("test", "Root positional args test")

	// No subcommand invoked — positional args land on cli.Args().
	argv := []string{"test", "alpha", "beta", "gamma"}
	leaf, err := cli.Parse(argv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if leaf != nil {
		t.Fatalf("expected nil subcommand (no registered subcommands), got %q", leaf.Name())
	}

	args := cli.Args()
	if len(args) != 3 {
		t.Fatalf("expected 3 root positional args, got %d: %v", len(args), args)
	}
	if args[0] != "alpha" || args[1] != "beta" || args[2] != "gamma" {
		t.Errorf("unexpected root positional args: %v", args)
	}
}
