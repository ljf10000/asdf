package asdf

import (
	"encoding/binary"
)

type ENetworkDir int

const (
	ToServer ENetworkDir = 0
	ToClient ENetworkDir = 1
)

func Ntohs(bin []byte) uint16 {
	return binary.BigEndian.Uint16(bin)
}

func Ntohl(bin []byte) uint32 {
	return binary.BigEndian.Uint32(bin)
}

func Ntohll(bin []byte) uint64 {
	return binary.BigEndian.Uint64(bin)
}

func Htons(bin []byte, v uint16) {
	binary.BigEndian.PutUint16(bin, v)
}

func Htonl(bin []byte, v uint32) {
	binary.BigEndian.PutUint32(bin, v)
}

func Htonll(bin []byte, v uint64) {
	binary.BigEndian.PutUint64(bin, v)
}
