package bdtp

import (
	"encoding/binary"
	"fmt"
	"github.com/nic758/bdtp-golang/blockchain"
	"github.com/nic758/bdtp-golang/utils"
	"log"
	"net"
	"os"
)

type chainSet struct {
	addressLen int
	txDataLen  int
}

var Chains = map[string]chainSet{
	blockchain.Waves_prefix: {addressLen: 26, txDataLen: 140},
}

func start(p string) (net.Listener, error) {
	log.Printf("Starting blockchain data transfer protocol on PORT :%s", p)
	s := fmt.Sprintf(":%s", p)
	return net.Listen("tcp", s)
}

func listen(serv net.Listener) {
	log.Println("Waiting for connections...")
	conn, err := serv.Accept()
	log.Printf("Client:%s connected!\n", conn.RemoteAddr())
	if err != nil {
		log.Printf("Could not connect with %s REASON:\n%s", conn.RemoteAddr(), err)
		return
	}

	chainPrefix := make([]byte, 3)
	_, err = conn.Read(chainPrefix)
	if err != nil {

		log.Printf("Failed to read chain  prefix, Client: %s", conn.RemoteAddr())
		conn.Close()
		return
	}

	chain, ok := Chains[string(chainPrefix)]
	if !ok {
		log.Printf("Closing connection for client:%s\n", conn.RemoteAddr())
		log.Printf("chain:%s is not supported\n", string(chainPrefix))
		if err = conn.Close(); err != nil {
			log.Fatal(err)
		}

		return
	}

	//TODO: maybe check if address is available on that chain.
	address := make([]byte, chain.addressLen)
	_, err = conn.Read(address)
	if err != nil {
		log.Printf("Failed to read chain  address\rClient: %s", conn.RemoteAddr())
		conn.Close()
		return
	}

	dataSize := make([]byte, 4)
	if _, err = conn.Read(dataSize); err != nil {
		log.Printf("Failed to read data size\rClient: %s", conn.RemoteAddr())
		conn.Close()
		return
	}

	l := binary.BigEndian.Uint32(dataSize)
	bc := blockchain.Factory(string(chainPrefix))

	if l < 0 {
		if err = conn.Close(); err != nil {
			log.Fatal(err)
		}

		return
	}

	if l == 0 {
		r, err := bc.FetchData(address)
		if err != nil {

			//may return an error
			conn.Close()
			return
		}

		_, err = conn.Write(r)
		if err = conn.Close(); err != nil {
			//TODO: should not crash program
			log.Fatal(err)
		}

		return
	}

	//save data
	data := make([]byte, l)
	_, err = conn.Read(data)
	if err != nil {
		log.Printf("Cant read data from %s", conn.RemoteAddr())
		conn.Close()
		return
	}

	if err = bc.ForgeData(address, data); err != nil {
		//TODO: should not crash program
		log.Printf("error in forging proccess for %s", conn.RemoteAddr())
		log.Printf(err.Error())
		conn.Close()
		return
	}

	r := []byte("WE NEED TO SEND A CONFIRMATION, it migth be all tx ids.\n")
	_, err = conn.Write(utils.ConvertInt32ToBytes(int32(len(r))))
	_, err = conn.Write(r)
	if err = conn.Close(); err != nil {
		//TODO: should not crash program
		log.Fatal(err)
	}
}

//ALL DATA SHOULD NOT BE ENCODED!!!
func NewServer() error {
	p := os.Getenv("PORT")

	serv, err := start(p)
	if err != nil {
		log.Fatal(err)
	}

	for {
		listen(serv)
	}
}
