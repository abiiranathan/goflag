package goflag

import (
	"fmt"
	"io"
	"time"
)

type Getter interface {
	Get(name string) any
	GetString(name string) string
	GetInt(name string) int
	GetInt64(name string) int64
	GetFloat32(name string) float32
	GetFloat64(name string) float64
	GetBool(name string) bool
	GetStringSlice(name string) []string
	GetRune(name string) rune
	GetTime(name string) time.Time
	GetDuration(name string) time.Duration
	GetIntSlice(name string) []int
}

// A subcommand. It can have its own flags.
type subcommand struct {
	name        string   // Subcommand name. used as a key to find the subcommand.
	description string   // Description of what this subcommand does.
	Handler     Handler  // Subcommand callback handler. Will be invoked by user if it matches.
	flags       []*gflag // subcommand flags.
}

// Create a new subcommand.
func SubCommand(name, description string, handler Handler) *subcommand {
	return &subcommand{
		name:        name,
		description: description,
		Handler:     handler,
	}
}

// Add a flag to a subcommand.
func (cmd *subcommand) AddFlag(f *gflag) *subcommand {
	if f.name == "" {
		panic("flag name can't be empty")
	}
	cmd.flags = append(cmd.flags, f)
	return cmd
}

func (cmd *subcommand) Get(name string) any {
	for _, flag := range cmd.flags {
		if flag.name == name || flag.shortName == name {
			return flag.value
		}
	}
	panic(fmt.Sprintf("flag %q not found", name))
}

// Add validator to last flag in the subcommand chain. If no flag exists, it panics.
func (cmd *subcommand) Validate(validator func(any) (bool, string)) *subcommand {
	if len(cmd.flags) == 0 {
		panic("no flags to validate")
	}
	cmd.flags[len(cmd.flags)-1].validator = validator
	return cmd
}

func (cmd *subcommand) PrintUsage(w io.Writer) {
	printSubCommand(cmd, w, 0, "")
}
