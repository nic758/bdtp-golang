package cmd

import (
	"github.com/nic758/bdtp-golang/cli"
	"log"
	"os"
	"path/filepath"
)

func newApp(name string) *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Commands = []cli.Command{srvCommand, forgeCommand, readCommand}

	return app
}

func Main(args []string) {
	if err := newApp(filepath.Base(args[0])).Run(args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
