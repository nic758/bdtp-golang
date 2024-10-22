package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	Polygon_prefix                 = "POL"
	polygon_testnet_default_node   = "https://polygon-amoy-bor-rpc.publicnode.com"
	polygon_testnet_block_explorer = "https://api-amoy.polygonscan.com/api"
	polygon_testnet                = 80002
	polygon_data_maxLength         = 100 * 1024 //100kByte
	polygon_min_amount             = 0
	polyscan_api_key_env           = "POLYSCAN_API_KEY"
)

type Polygon struct {
	pubKey  *ecdsa.PublicKey
	privKey *ecdsa.PrivateKey
	address common.Address
	apiKey  string
}

type PolygonTransaction struct {
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	Gas             string `json:"gas"`
	GasPrice        string `json:"gasPrice"`
	IsError         string `json:"isError"`
	TxReceiptStatus string `json:"txreceipt_status"`
	Input           string `json:"input"`
}

type ApiResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Result  []PolygonTransaction `json:"result"`
}

func NewPolygon() *Polygon {
	seed, ok := os.LookupEnv(bdtp_env_seed)
	if !ok || seed == "" {
		fmt.Println("Seed missing: unable to initialize polygon network")
		return nil
	}

	h := sha256.Sum256([]byte(seed))
	privKey, _ := crypto.ToECDSA(h[:])
	pubKey := privKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*pubKey)
	fmt.Println("Polygon initialized! \nAddress: ", fromAddress.Hex())

	apiKey, ok := os.LookupEnv(polyscan_api_key_env)
	if !ok {
		fmt.Printf("%s not set", polyscan_api_key_env)
		return nil
	}

	return &Polygon{
		privKey: privKey,
		pubKey:  pubKey,
		address: fromAddress,
		apiKey:  apiKey,
	}
}

func (p *Polygon) fetchBlockExplorer(address []byte) ([]PolygonTransaction, error) {
	query := url.Values{}
	query.Set("module", "account")
	query.Set("action", "txlist")
	query.Set("address", string(address))
	query.Set("startblock", "0")
	query.Set("endblock", "99999999")
	query.Set("page", "1")
	query.Set("offset", "100")
	query.Set("sort", "asc")
	query.Set("apikey", p.apiKey)
	requestURL := fmt.Sprintf("%s?%s", polygon_testnet_block_explorer, query.Encode())

	response, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
		return nil, err
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
		return nil, err
	}

	if apiResponse.Status != "1" {
		log.Printf("Error from Polygonscan: %s\n", apiResponse.Message)
		return nil, err
	}

	return apiResponse.Result, nil
}

func (p *Polygon) FetchData(address []byte) ([]byte, error) {
	log.Printf("fetching data from polygon blockchain at %s\n", address)
	txs, err := p.fetchBlockExplorer(address)
	if err != nil {
		return nil, err
	}
	if len(txs) <= 0 {
		return nil, errors.New("no tx found for this pointer")
	}

	totalBytes, err := strconv.Atoi(txs[0].Value)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	bytesW := 0
	var buf bytes.Buffer

	for _, tx := range txs {
		db := tx.Input[2:]
		dataBytes, err := hex.DecodeString(db)
		if err != nil {
			fmt.Println("error: ", err)
			return nil, err
		}
		n, err := buf.Write(dataBytes)
		if err != nil {
			return nil, err
		}

		bytesW += n
	}

	if totalBytes < bytesW {
		err := errors.New("bytes written mismatch total bytes")
		fmt.Println("error: ", err)

		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Polygon) sendTransaction(client *ethclient.Client, address common.Address, data []byte, txAmount int) (int, error) {
	nonce, err := client.PendingNonceAt(context.Background(), p.address)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
		return 0, err
	}

	value := big.NewInt(int64(txAmount))
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
		return 0, err
	}

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:    &address,
		From:  p.address,
		Value: value,
		Data:  data,
	})
	if err != nil {
		return 0, err
	}

	tx := types.NewTransaction(nonce, address, value, gasLimit, gasPrice, data)
	chainID := big.NewInt(polygon_testnet)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), p.privKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
		return 0, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("Failed to send transaction: %v", err)
		return 0, err
	}

	return len(data), nil
}

func (p *Polygon) ForgeData(address []byte, data []byte) error {
	/*	if len(data) > polygon_data_maxLength {
		err := errors.New("the data size exceeds the allowed limit")
		fmt.Println(err)
		return err
	}*/

	log.Printf("posting data to polygon blockchain at %s...!\n", address)

	//mainnet and testnet is not in recipientAddress like waves
	client, err := ethclient.Dial(polygon_testnet_default_node)
	if err != nil {
		return err
	}

	toAddress := common.HexToAddress(string(address))

	w := 0
	txAmount := len(data)
	for i := 0; i < len(data); i += polygon_data_maxLength {
		end := i + polygon_data_maxLength
		if end > len(data) {
			end = len(data)
		}

		subArr := data[w:end]
		c, err := p.sendTransaction(client, toAddress, subArr, txAmount)

		if err != nil {
			return err
		}

		w += c
		txAmount = polygon_min_amount
	}

	return nil
}
