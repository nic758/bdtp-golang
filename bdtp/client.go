package bdtp

import (
	"bufio"
	"encoding/binary"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/nic758/bdtp-golang/utils"
	"log"
	"net"
)

type Pointer string

func (p Pointer) GetChain() Pointer {
	return p[:3]
}

func (p Pointer) GetAddress() Pointer {
	return p[3:]
}

type bdtpClient struct {
	ad string
}

func NewClient(address string) *bdtpClient {
	return &bdtpClient{ad: address}
}

//address should be freshly generated.
func (c *bdtpClient) SavaDataToChain(chain, address string, data []byte) string {
	conn, err := net.Dial("tcp", c.ad)
	if err != nil {
		log.Fatal("Cannot connect to the server")
	}

	buf := bufio.NewWriter(conn)
	buf.Write([]byte(chain))
	buf.Write(base58.Decode(address))

	size := utils.ConvertInt32ToBytes(int32(len(data)))
	buf.Write(size)
	buf.Write(data)
	buf.Flush()

	dataSize := make([]byte, 4)
	if _, err = conn.Read(dataSize); err != nil {
		log.Fatal(err)
	}

	l := binary.BigEndian.Uint32(dataSize)
	d := make([]byte, l)
	if _, err = conn.Read(d); err != nil {
		log.Fatal(err)
	}

	//log.Printf(string(d))

	if err = conn.Close(); err != nil {
		log.Fatal(err)
	}
	return address
}

func (c *bdtpClient) FetchDataFromChain(pointer Pointer) []byte {
	chain := pointer.GetChain()
	address := pointer.GetAddress()

	conn, err := net.Dial("tcp", c.ad)
	if err != nil {
		return nil
	}
	buf := bufio.NewWriter(conn)
	buf.Write([]byte(chain))
	buf.Write(base58.Decode(string(address)))

	size := utils.ConvertInt32ToBytes(int32(0))
	buf.Write(size)
	buf.Flush()
	if err != nil {
		//TODO
		log.Printf(err.Error())
		return nil
	}

	dataSize := make([]byte, 4)
	if _, err = conn.Read(dataSize); err != nil {
		log.Fatal(err)
	}

	l := binary.BigEndian.Uint32(dataSize)
	d := make([]byte, l)
	if _, err = conn.Read(d); err != nil {
		log.Println(err)
		log.Println("Data may no be confirmed on the blockchain.")
		return nil
	}

	conn.Close()
	return d
}
