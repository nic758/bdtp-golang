package blockchain

import (
	"bytes"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/stretchr/testify/assert"
	"testing"
)

//all values are from waves official doc at https://docs.waves.tech/en/blockchain/waves-protocol/cryptographic-practical-details#description
const (
	privK   = "7VLYNhmuvAo5Us4mNGxWpzhMSdSSdEbEPFUDKSnA6eBv"
	pubK    = "EENPV1mRhUD9gSKbcWt84cqnfSGQP5LkCu5gMBfAanYH"
	address = "3N9Q2sdkkhAnbR4XCveuRaSMLiVtvebZ3wp" //mainnet
	txBytes = "Ht7FtLJBrnukwWtywum4o1PbQSNyDWMgb4nXR5ZkV78krj9qVt17jz74XYSrKSTQe6wXuPdt3aCvmnF5hfjhnd1gyij36hN1zSDaiDg3TFi7c7RbXTHDDUbRgGajXci8PJB3iJM1tZvh8AL5wD4o4DCo1VJoKk2PUWX3cUydB7brxWGUxC6mPxKMdXefXwHeB4khwugbvcsPgk8F6YB"
)

func Test_generate_address(t *testing.T) {
	add := GetWavesAddress(base58.Decode(pubK), int8('T'))
	assert.Equal(t, address, base58.Encode(add))
}

func Test_bodyBytes(t *testing.T) {
	s := base58.Decode(txBytes)
	ty := s[0]
	senderPubK := s[1:33]
	amountAssetFlag := byte(0)
	//amountAssetId := s[34:66]
	feeAssetFlag := byte(0)
	//feeAssetId := s[67:99]
	ts := s[99:107]
	amount := s[107:115]
	fee := s[115:123]
	recipient := s[123:149]
	attLen := s[149:151]
	attach := s[151:]

	buf := &bytes.Buffer{}
	buf.WriteByte(ty)
	buf.Write(senderPubK)
	buf.WriteByte(amountAssetFlag)
	//buf.Write(amountAssetId)
	buf.WriteByte(feeAssetFlag)
	//buf.Write(feeAssetId)
	buf.Write(ts)
	buf.Write(amount)
	buf.Write(fee)
	buf.Write(recipient)
	buf.Write(attLen)
	buf.Write(attach)

	expected := buf.Bytes()

	actualTx := &TransactionTransfer{
		Version:         2,
		Type:            4,
		SenderPk:        "EENPV1mRhUD9gSKbcWt84cqnfSGQP5LkCu5gMBfAanYH",
		AmountAssetFlag: 0,
		FeeAssetFlag:    0,
		Timestamp:       1479287120875,
		Amount:          1,
		Fee:             1,
		Recipient:       "3NBVqYXrapgJP9atQccdBPAgJPwHDKkh6A8",
		AttachmentLen:   int16(4),
		Attachment:      "2VfUX",
	}

	actual := actualTx.toBytes()
	assert.Equal(t, base58.Encode(expected), base58.Encode(actual))
}
