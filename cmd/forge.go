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

var forgeCommand = cli.Command{
	Name:   "forge",
	Usage:  "forge a file in the blockchain",
	Action: forge,
	Flags:  CommonClientFlags,
}

func generatePolygonAddress() string {
	// Generate a new private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Derive the public key from the private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// Generate the address from the public key
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("Generated Polygon address: %s\n", address)

	return address
}

func generateWavesAddress() string {
	c := wavesplatform.NewWavesCrypto()
	seed := c.RandomSeed()
	pair := c.KeyPair(seed)

	fmt.Println("new seed: %s\n", seed)
	fmt.Println("new secret: %s\n", pair.PrivateKey)
	fmt.Println("new public: %s\n", pair.PublicKey)
	newAdd := c.Address(pair.PublicKey, 84)
	fmt.Println("new address: %s\n", newAdd)

	return string(newAdd)
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

	newAdd := generatePolygonAddress()
	log.Printf("file in bytes: %d", len(d))
	fmt.Println("add", newAdd)
	a := client.SavaDataToChain("POL", newAdd, d)

	log.Printf("Data stored at POL%s", a)
	return nil
}
