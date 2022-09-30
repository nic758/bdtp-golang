package cli

import (
	"log"
)

var helpCommand = Command{
	Name:    "help",
	Aliases: []string{"h"},
	Usage:   "Shows a list of commands of help for a command",
	Action: func(ctx *Context) error {
		//TODO
		log.Println("TODO: show available commands and flags")
		return nil
	},
}
