package cmd

import (
	"errors"
	"os"

	"github.com/nic758/bdtp-golang/bdtp"
	"github.com/nic758/bdtp-golang/cli"
)

var readCommand = cli.Command{
	Name:   "read",
	Usage:  "reads data from the provided pointer",
	Action: read,
	Flags:  CommonClientFlags,
}

func read(ctx *cli.Context) error {
	client := bdtp.NewClient(os.Getenv("BDTP_HOST"))
	if ctx.Args.Last() == "" {
		return errors.New("a pointers must be provided")
	}

	b := client.FetchDataFromChain(bdtp.Pointer(ctx.Args.Last()))
	os.WriteFile("out", b, 0644)

	return nil
}
