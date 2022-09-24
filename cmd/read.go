package cmd

import (
	"errors"
	"github.com/nic758/bdtp-golang/bdtp"
	"github.com/nic758/bdtp-golang/cli"
	"log"
	"os"
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
		return errors.New("A pointers must be provided")
	}

	b := client.FetchDataFromChain(bdtp.Pointer(ctx.Args.Last()))
	log.Printf("Data: \n" + string(b))

	return nil
}
