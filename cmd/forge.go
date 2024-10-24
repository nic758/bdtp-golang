package cmd

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nic758/bdtp-golang/bdtp"
	"github.com/nic758/bdtp-golang/blockchain"
	"github.com/nic758/bdtp-golang/cli"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
)

var CommonClientFlags = []cli.Flag{
	&cli.StringFlag{
		Name:   "host",
		Value:  "localhost:4444",
		Usage:  "",
		EnvVar: "BDTP_HOST",
	},
}

var forgeFlags = []cli.Flag{
	&cli.StringFlag{
		Name:   "blockchain",
		Value:  "POL",
		Usage:  "specify the blockchain to forge data",
		EnvVar: "BDTP_CHAIN",
	},
}

var forgeCommand = cli.Command{
	Name:   "forge",
	Usage:  "forge a file in the blockchain",
	Action: forge,
	Flags:  append(CommonClientFlags, forgeFlags...),
}

func generatePolygonAddress() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("Generated Polygon address: %s\n", address)

	return address
}

func generateWavesAddress() string {
	c := wavesplatform.NewWavesCrypto()
	seed := c.RandomSeed()
	pair := c.KeyPair(seed)

	newAdd := c.Address(pair.PublicKey, 84)

	return string(newAdd)
}

func generatePointer() string {
	chain := os.Getenv("BDTP_CHAIN")

	switch chain {
	case blockchain.Polygon_prefix:
		return generatePolygonAddress()
	case blockchain.Waves_prefix:
		return generateWavesAddress()
	}

	log.Fatal("BDTP_CHAIN doesn't exist")
	return ""
}

func forge(ctx *cli.Context) error {
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

	newAdd := generatePointer()
	log.Printf("file in bytes: %d", len(d))
	fmt.Println("add", newAdd)
	_ = client.SavaDataToChain(os.Getenv("BDTP_CHAIN"), newAdd, d)

	log.Printf("Data stored at %s%s", os.Getenv("BDTP_CHAIN"), newAdd)

	return nil
}
