package main

import (
	//wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/nic758/bdtp-golang/bdtp"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	c := wavesplatform.NewWavesCrypto()
	seed := c.RandomSeed()
	pair := c.KeyPair(seed)

	log.Printf("new seed: %s\n", seed)
	log.Printf("new secret: %s\n", pair.PrivateKey)
	log.Printf("new public: %s\n", pair.PublicKey)
	newAdd := c.Address(pair.PublicKey, 84)
	log.Printf("new address: %s\n", newAdd)

	//fetchAddress("3N5wR9pAjJpzrdZ44dg5Qnpps4zf5SBujiT")
	//fetchAddress("3N5zVw8UA8xjFWsvtWUB4zAwFRRCDFeHg8y")
	//custom(string(newAdd))
	//saveString(string(newAdd), "premiere demo")
	//fetchAddress("3MzmtsXMowUwjm27e2am8qHCoCDtHfr54sJ")
	//saveFile(string(newAdd))
	fetchAddress("3MzPNnn4Zq1WNzqzd1QHKuSEgB7Ncc2Ri7h")

}

func custom(address string) {
	saveFile(address)
	fetchAddress(address)
}

func fetchAddress(address string) {
	client := bdtp.NewClient("localhost:4444")
	b := client.FetchDataFromChain("WAV", address)
	log.Printf(string(b))
}

func saveString(address, s string) {
	client := bdtp.NewClient("localhost:4444")
	client.SavaDataToChain("WAV", address, []byte(s))
}

func saveFile(address string) {
	client := bdtp.NewClient("localhost:4444")

	f, err := os.Open("bdtp-file.txt")
	if err != nil {
		log.Fatal(err)
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("file in bytes: %d", len(d))
	a := client.SavaDataToChain("WAV", string(address), d)

	log.Printf("Data stored at %s", a)
}
