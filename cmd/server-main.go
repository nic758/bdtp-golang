package cmd

import (
	"github.com/nic758/bdtp-golang/bdtp"
	"github.com/nic758/bdtp-golang/cli"
)

var DefaultPort = "4444"

var srvFlags = []cli.Flag{
	&cli.StringFlag{
		Name:   "port",
		Value:  DefaultPort,
		Usage:  "bind a specific :PORT to the server",
		EnvVar: "PORT",
	},
	&cli.StringFlag{
		Name:   "bdtp_seed",
		Value:  "",
		Usage:  "the default seed for the server, if no seed is provided the server will be in read-only mode",
		EnvVar: "BDTP_SEED",
	},
}

var srvCommand = cli.Command{
	Name:    "server",
	Aliases: []string{"srv"},
	Usage:   "starts a bdtp server",
	Action:  startServer,
	Flags:   srvFlags,
}

func startServer(ctx *cli.Context) error {
	return bdtp.NewServer()
}
