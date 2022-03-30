package utils

import "encoding/binary"

func ConvertInt64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))

	return b
}

func ConvertInt16ToBytes(i int16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(i))

	return b
}

func ConvertInt32ToBytes(i int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))

	return b
}
