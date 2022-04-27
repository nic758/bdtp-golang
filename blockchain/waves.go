package blockchain

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/nic758/bdtp-golang/utils"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

//TODO
const (
	Waves_prefix = "WAV"

	waves_env_seed             = "WAVES_SEED"
	waves_testnet_default_node = "https://nodes-testnet.wavesnodes.com"
	waves_mainnet_default_node = "https://nodes.wavesnodes.com"
	waves_mainnet              = 87
	waves_testnet              = 84
	waves_data_maxLength       = 140
	waves_min_amount           = 1

	Error_Env                = "Env environment variable cannot be found. This node cannot broadcast data to the selected Blockchain. please provide a valid value.\n"
	Error_blockchain_network = "The provided network is not supported\n"
)

type Waves struct{}
type TransactionTransfer struct {
	Version         int    `json:"version"`
	Type            int    `json:"type"`
	SenderPk        string `json:"senderPublicKey"`
	AmountAssetFlag int    `json:"amountAssetFlag"`
	FeeAssetFlag    int    `json:"feeAssetFlag"`
	Timestamp       int64  `json:"timestamp"`
	Amount          int64  `json:"amount"`
	Fee             int64  `json:"fee"`
	Recipient       string `json:"recipient"`
	AttachmentLen   int16
	Attachment      string `json:"attachment"`
	Signature       string `json:"signature"`

	//not used when posting a tx
	Height int    `json:"height"`
	Id     string `json:"id"`
	Sender string `json:"sender"`
}

func (Waves) FetchData(address []byte) ([]byte, error) {
	log.Printf("fetching data from waves blockchain at %s\n", base58.Encode(address))
	net := address[1]
	netUrl := ""

	if net == waves_mainnet {
		netUrl = waves_mainnet_default_node
	}

	if net == waves_testnet {
		netUrl = waves_testnet_default_node
	}

	if netUrl == "" {
		return nil, errors.New(fmt.Sprintf("%s", Error_blockchain_network))
	}

	add := base58.Encode(address)
	//TODO: limit is 1000 transactions so max data is 140bytes*1000
	//we can process in batch if file is too long.
	r, err := http.Get(fmt.Sprintf("%s/transactions/address/%s/limit/%d", netUrl, add, 1000))
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode != http.StatusOK {
		log.Println(string(body))
		return make([]byte, 0), nil
	}

	var txs [][]TransactionTransfer
	json.Unmarshal(body, &txs)

	if len(txs) == 0 || len(txs[0]) == 0 {
		log.Printf("No data found at %s", add)
		return make([]byte, 0), nil
	}
	sort.Slice(txs[0][:], func(i, j int) bool {
		return txs[0][i].Timestamp < txs[0][j].Timestamp
	})

	dataLen := txs[0][0].Amount
	c := int64(0)
	buf := bytes.Buffer{}
	buf.Write(utils.ConvertInt32ToBytes(int32(dataLen)))
	for i := 0; c < dataLen; i++ {
		if i >= len(txs[0]) {
			log.Printf("Corrupted data.. only %d found on %d asked.", c, dataLen)
			return nil, errors.New("maybe one of the transactions are not confirm in the blockchain")
		}

		tx := txs[0]
		encodedData := tx[i].Attachment
		data := base58.Decode(encodedData)
		c += int64(len(data))
		buf.Write(data)
	}

	return buf.Bytes(), nil
}
func (Waves) ForgeData(address []byte, data []byte) error {
	log.Printf("posting data to waves blockchain at %s...!\n", base58.Encode(address))
	net := address[1]
	netUrl := ""

	if net == waves_mainnet {
		netUrl = waves_mainnet_default_node
	}

	if net == waves_testnet {
		netUrl = waves_testnet_default_node
	}

	if netUrl == "" {
		return errors.New(fmt.Sprintf("%s", Error_blockchain_network))
	}

	seed, ok := os.LookupEnv(waves_env_seed)
	if !ok {
		//TODO create specific error type.
		return errors.New(fmt.Sprintf("%s%s need to be set.", Error_Env, waves_env_seed))
	}

	c := wavesplatform.NewWavesCrypto()
	keyPair := c.KeyPair(wavesplatform.Seed(seed))

	pubK := base58.Decode(string(keyPair.PublicKey))
	privK := base58.Decode(string(keyPair.PrivateKey))

	return sendTransactions(netUrl, pubK, privK, address, data)
}

func sendTransactions(net string, pubK []byte, k ed25519.PrivateKey, recipientAddress, data []byte) error {
	offset := 0
	amount := int64(len(data))
	for offset <= len(data) {
		end := offset + waves_data_maxLength
		if end > len(data) {
			end = len(data)
		}

		tx := CreateAndSignTransaction(pubK, k, recipientAddress, data[offset:end], amount)

		jsonTx, err := json.Marshal(tx)
		if err != nil {
			//TODO: application should not crash
			log.Fatal(err)
		}

		req := bytes.NewReader(jsonTx)
		route := fmt.Sprintf("%s/%s", net, "transactions/broadcast")
		r, err := http.Post(route, "application/json", req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close()

		if r.StatusCode != http.StatusOK {
			log.Printf(string(body))
			return errors.New("failed to post data to waves blockchain")
		}
		if err != nil {
			log.Fatal(err)
		}

		//TODO: we might want to return something to the user.
		respTx := TransactionTransfer{}
		err = json.Unmarshal(body, &respTx)
		if err != nil {
			log.Fatal(err)
		}

		offset += waves_data_maxLength
		amount = waves_min_amount
	}
	return nil
}

func CreateAndSignTransaction(pubK []byte, privK []byte, recipient, data []byte, amount int64) TransactionTransfer {
	tx := &TransactionTransfer{
		Version:         1,
		Type:            4,
		SenderPk:        base58.Encode(pubK),
		AmountAssetFlag: 0,
		FeeAssetFlag:    0,
		Timestamp:       time.Now().UnixMilli(),
		Amount:          amount,
		Fee:             100000,
		Recipient:       base58.Encode(recipient),
		AttachmentLen:   int16(len(data)),
		Attachment:      base58.Encode(data),
	}

	tx.sign(privK)

	return *tx
}

func (t *TransactionTransfer) sign(k ed25519.PrivateKey) {
	txData := t.toBytes()
	c := wavesplatform.NewWavesCrypto()
	sig := c.SignBytes(txData, wavesplatform.PrivateKey(base58.Encode(k)))

	//sig := Sign(k, txData)
	t.Signature = base58.Encode(sig)
}

func (t TransactionTransfer) toBytes() []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(byte(t.Type))
	//buf.WriteByte(byte(2))
	buf.Write(base58.Decode(t.SenderPk))
	buf.WriteByte(byte(0))
	buf.WriteByte(byte(0))
	buf.Write(utils.ConvertInt64ToBytes(t.Timestamp))
	buf.Write(utils.ConvertInt64ToBytes(t.Amount))
	buf.Write(utils.ConvertInt64ToBytes(t.Fee))
	buf.Write(base58.Decode(t.Recipient))
	buf.Write(utils.ConvertInt16ToBytes(t.AttachmentLen))
	buf.Write(base58.Decode(t.Attachment))

	return buf.Bytes()
}

func GetWavesAddress(publicKey []byte, net int8) []byte {
	b := make([]byte, 0)
	t := int8(1)
	h := secureHash(publicKey)
	b = append(b, byte(t), byte(net))
	b = append(b, h[:20]...)

	hh := secureHash(b)

	return append(b, hh[:4]...)
}
func secureHash(data []byte) [32]byte {
	b := blake2b.Sum256(data)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(b[:])
	copy(b[:32], hash.Sum(nil))
	return b
}
