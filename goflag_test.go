package goflag

import (
	"bytes"
	"errors"
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
			Handler: func(userdata any) error {
				fmt.Printf("Hello %s, you are %d years old and %d cm tall\n", name, age, height)
				return nil
			},
		},
		{
			name:        "test2",
			description: "Test command 2",
			Handler: func(userdata any) error {
				fmt.Printf("Hello %s, you are %d years old and %d cm tall\n", name, age, height)
				return nil
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
	cli.SubCommand("test", "Test command", func(userdata any) error {
		fmt.Println("Test command")
		return nil
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

	deployCmd := cli.SubCommand("deploy", "Deploy something", func(userdata any) error { return nil })

	deployCmd.SubCommand("staging", "Deploy to staging", func(userdata any) error {
		output = fmt.Sprintf("staging count=%d dry=%v", count, dryRun)
		return nil
	}).
		Int("count", "c", &count, "Number of replicas").Required().
		Bool("dry-run", "d", &dryRun, "Dry run mode")

	deployCmd.SubCommand("prod", "Deploy to production", func(userdata any) error {
		output = fmt.Sprintf("prod count=%d dry=%v", count, dryRun)
		return nil
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

	if err := leaf.Handler(nil); err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	if output != "staging count=3 dry=true" {
		t.Errorf("unexpected handler output: %q", output)
	}
}

func TestPositionalArgs(t *testing.T) {
	cli := New("test", "Positional args test")

	var format string
	var captured []string

	cli.SubCommand("convert", "Convert files", func(userdata any) error {
		captured = userdata.([]string)
		return nil
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
	if err := leaf.Handler(positional); err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
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

// ---------------------------------------------------------------------------
// Feature 2 — Handler error propagation
// ---------------------------------------------------------------------------

// TestHandlerErrorPropagation verifies that ParseAndInvoke returns errors from
// the matched subcommand handler.
func TestHandlerErrorPropagation(t *testing.T) {
	cli := New("test", "Handler error test")
	sentinel := errors.New("handler failed")

	cli.SubCommand("fail", "Always fails", func(userdata any) error {
		return sentinel
	})

	err := cli.ParseAndInvoke([]string{"test", "fail"}, nil, nil)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

// TestHandlerNilErrorOnSuccess verifies that a successful handler returns nil.
func TestHandlerNilErrorOnSuccess(t *testing.T) {
	cli := New("test", "Handler success test")

	cli.SubCommand("ok", "Always succeeds", func(userdata any) error {
		return nil
	})

	if err := cli.ParseAndInvoke([]string{"test", "ok"}, nil, nil); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// Feature 3 — Mutually exclusive flags
// ---------------------------------------------------------------------------

// TestExclusiveFlagsOnSubCommand verifies that providing two mutually exclusive
// flags on a subcommand returns an error.
func TestExclusiveFlagsOnSubCommand(t *testing.T) {
	cli := New("test", "Exclusive flags test")
	var jsonOut, yamlOut bool

	cli.SubCommand("export", "Export data", func(userdata any) error { return nil }).
		Bool("json", "j", &jsonOut, "Output as JSON").
		Bool("yaml", "y", &yamlOut, "Output as YAML").
		ExclusiveFlags("json", "yaml")

	// Both flags set — should error.
	_, err := cli.Parse([]string{"test", "export", "--json", "--yaml"})
	if err == nil {
		t.Fatal("expected error when both exclusive flags are set, got nil")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("unexpected error text: %v", err)
	}

	// Only one flag set — should succeed.
	jsonOut, yamlOut = false, false
	_, err = cli.Parse([]string{"test", "export", "--json"})
	if err != nil {
		t.Errorf("expected no error with a single exclusive flag, got %v", err)
	}
}

// TestExclusiveFlagsOnCLI verifies the same constraint at the root level.
func TestExclusiveFlagsOnCLI(t *testing.T) {
	cli := New("test", "Root exclusive flags test")
	var aVal, bVal bool

	cli.Bool("alpha", "a", &aVal, "Alpha mode")
	cli.Bool("beta", "b", &bVal, "Beta mode")
	cli.ExclusiveFlags("alpha", "beta")

	_, err := cli.Parse([]string{"test", "--alpha", "--beta"})
	if err == nil {
		t.Fatal("expected error for conflicting root flags, got nil")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("unexpected error text: %v", err)
	}
}

// TestExclusiveFlagsNeitherSet verifies that providing neither flag in a group
// is accepted (exclusivity only prohibits simultaneous use, not absence).
func TestExclusiveFlagsNeitherSet(t *testing.T) {
	cli := New("test", "Neither exclusive flag set")
	var aVal, bVal bool

	cli.SubCommand("run", "Run something", func(userdata any) error { return nil }).
		Bool("fast", "f", &aVal, "Fast mode").
		Bool("safe", "s", &bVal, "Safe mode").
		ExclusiveFlags("fast", "safe")

	_, err := cli.Parse([]string{"test", "run"})
	if err != nil {
		t.Errorf("expected no error when neither exclusive flag is set, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// Feature 1 — Persistent flags
// ---------------------------------------------------------------------------

// TestPersistentFlagInheritedByChild verifies that a persistent flag declared
// on a parent subcommand is visible and parseable by a nested child.
func TestPersistentFlagInheritedByChild(t *testing.T) {
	cli := New("test", "Persistent flag test")
	var format string
	var count int

	parentCmd := cli.SubCommand("process", "Process data", func(userdata any) error { return nil })

	// Persistent flag on parent — all children should accept it.
	parentCmd.PersistentString("format", "f", &format, "Output format")

	parentCmd.SubCommand("items", "Process items", func(userdata any) error { return nil }).
		Int("count", "c", &count, "Number of items").Required()

	argv := []string{"test", "process", "items", "--count", "5", "--format", "csv"}
	leaf, err := cli.Parse(argv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if leaf == nil {
		t.Fatal("expected leaf subcommand, got nil")
	}
	if leaf.Name() != "items" {
		t.Errorf("expected leaf 'items', got %q", leaf.Name())
	}
	if format != "csv" {
		t.Errorf("expected format=csv (via persistent flag), got %q", format)
	}
	if count != 5 {
		t.Errorf("expected count=5, got %d", count)
	}
}

// TestPersistentFlagShadowedByChild verifies that a child's own flag with the
// same name as a parent's persistent flag takes precedence (no duplicate-flag
// error, child's pointer is populated).
func TestPersistentFlagShadowedByChild(t *testing.T) {
	cli := New("test", "Persistent flag shadowing test")
	var parentFormat string
	var childFormat string

	parentCmd := cli.SubCommand("do", "Do something", func(userdata any) error { return nil })
	parentCmd.PersistentString("format", "f", &parentFormat, "Output format (inherited)")

	// Child declares its own --format that points to childFormat.
	parentCmd.SubCommand("sub", "Subcommand", func(userdata any) error { return nil }).
		String("format", "F", &childFormat, "Output format (override)")

	argv := []string{"test", "do", "sub", "--format", "json"}
	_, err := cli.Parse(argv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Child's own flag takes precedence.
	if childFormat != "json" {
		t.Errorf("expected childFormat=json, got %q", childFormat)
	}
	// Parent's pointer is not touched.
	if parentFormat != "" {
		t.Errorf("expected parentFormat to remain empty, got %q", parentFormat)
	}
}

// TestPersistentFlagRequiredOnChild verifies that a required persistent flag
// that is not supplied yields a parse error.
func TestPersistentFlagRequiredOnChild(t *testing.T) {
	cli := New("test", "Required persistent flag test")
	var token string

	parentCmd := cli.SubCommand("api", "API commands", func(userdata any) error { return nil })
	parentCmd.PersistentString("token", "t", &token, "Auth token").Required()

	parentCmd.SubCommand("list", "List resources", func(userdata any) error { return nil })

	// --token omitted — should fail.
	_, err := cli.Parse([]string{"test", "api", "list"})
	if err == nil {
		t.Fatal("expected error for missing required persistent flag, got nil")
	}
	if !strings.Contains(err.Error(), "token") {
		t.Errorf("expected error to mention 'token', got: %v", err)
	}
}

// TestPrintUsageShowsInheritedFlags verifies that the help output for a child
// subcommand includes an "Inherited flags" section when a parent has persistent
// flags.
func TestPrintUsageShowsInheritedFlags(t *testing.T) {
	var token string
	parent := &SubCMD{
		name:        "api",
		description: "API commands",
		Handler:     func(userdata any) error { return nil },
		persistentFlags: []*Flag{
			{name: "token", shortName: "t", flagType: flagString, value: &token, usage: "Auth token"},
		},
	}
	child := &SubCMD{
		name:        "list",
		description: "List resources",
		Handler:     func(userdata any) error { return nil },
		parent:      parent,
		flags: []*Flag{
			{name: "help", shortName: "h", flagType: flagBool, usage: "Print help message and exit"},
		},
	}

	var buf bytes.Buffer
	child.PrintUsage(&buf)

	out := buf.String()
	if !strings.Contains(out, "Inherited flags") {
		t.Errorf("expected 'Inherited flags' section in child usage output, got:\n%s", out)
	}
	if !strings.Contains(out, "--token") {
		t.Errorf("expected '--token' in inherited flags section, got:\n%s", out)
	}
}
