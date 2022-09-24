package cmd

import (
	"errors"
	"github.com/nic758/bdtp-golang/bdtp"
	"github.com/nic758/bdtp-golang/cli"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"io/ioutil"
	"log"
	"os"
)

var CommonClientFlags = []cli.Flag{
	&cli.StringFlag{
		Name:   "host",
		Value:  "localhost:4444",
		Usage:  "",
		EnvVar: "BDTP_HOST",
	},
}

var forgeCommand = cli.Command{
	Name:   "forge",
	Usage:  "forge a file in the blockchain",
	Action: forge,
	Flags:  CommonClientFlags,
}

func forge(ctx *cli.Context) error {
	c := wavesplatform.NewWavesCrypto()
	seed := c.RandomSeed()
	pair := c.KeyPair(seed)

	//log.Printf("new seed: %s\n", seed)
	//log.Printf("new secret: %s\n", pair.PrivateKey)
	//log.Printf("new public: %s\n", pair.PublicKey)
	newAdd := c.Address(pair.PublicKey, 84)
	//log.Printf("new address: %s\n", newAdd)

	client := bdtp.NewClient(os.Getenv("BDTP_HOST"))
	if ctx.Args.Last() == "" {
		return errors.New("A file path must be provided")
	}

	f, err := os.Open(ctx.Args.Last())
	if err != nil {
		log.Fatal(err)
		return err
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("file in bytes: %d", len(d))
	a := client.SavaDataToChain("WAV", string(newAdd), d)

	log.Printf("Data stored at WAV%s", a)
	return nil
}
