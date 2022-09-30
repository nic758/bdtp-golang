package cli

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Author struct {
	Name  string // The Authors name
	Email string // The Authors email
}

type App struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string
	// Full name of command for help, defaults to Name
	HelpName string
	// Description of the program.
	Usage string
	// Text to override the USAGE section of help
	UsageText string
	// Description of the program argument format.
	ArgsUsage string
	// Version of the program
	Version string
	// Description of the program
	Description string
	// List of commands to execute
	Commands []Command
	// List of flags to parse
	Flags  []Flag
	Action interface{}

	// Execute this function if the proper command cannot be found
	//CommandNotFound CommandNotFoundFunc
	// Execute this function if an usage error occurs
	//OnUsageError OnUsageErrorFunc
	// Compilation date
	Compiled time.Time
	// List of all authors who contributed
	Authors []Author
	// Copyright of the binary if any
	Copyright string
	// Name of Author (Note: Use App.Authors, this is deprecated)
	Author string
	// Email of Author (Note: Use App.Authors, this is deprecated)
	Email string
	// Writer writer to write output to
	Writer io.Writer
	// ErrWriter writes error output
	ErrWriter io.Writer
}

func NewApp() *App {
	return &App{
		Name:      filepath.Base(os.Args[0]),
		HelpName:  filepath.Base(os.Args[0]),
		Usage:     "bdtp cli app",
		UsageText: "",
		Version:   "0.0.0",
		Compiled:  compileTime(),
		Action:    nil,
	}
}

func (a *App) Run(rawArgs []string) (err error) {
	ctx := NewContext(a)

	command := rawArgs[1:]
	if len(command) > 0 {
		n := command[0]
		c := a.Command(n)
		if c != nil {
			return c.Run(ctx, command[1:])
		}
	}

	if a.Action == nil {
		a.Action = helpCommand.Action
	}

	return HandleAction(a.Action, ctx)
}

func (a *App) Command(name string) *Command {
	for _, c := range a.Commands {
		if c.HasName(name) {
			return &c
		}
	}

	return nil
}

func HandleAction(action interface{}, ctx *Context) (err error) {
	if a, ok := action.(ActionFunc); ok {
		return a(ctx)
	} else {
		return errors.New("cannot handle action")
	}
}
func compileTime() time.Time {
	info, err := os.Stat(os.Args[0])
	if err != nil {
		return time.Now()
	}
	return info.ModTime()
}
