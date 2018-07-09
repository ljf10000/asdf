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

func Ntoh16(bin []byte) uint16 {
	return binary.BigEndian.Uint16(bin)
}

func Ntohl(bin []byte) uint32 {
	return binary.BigEndian.Uint32(bin)
}

func Ntoh32(bin []byte) uint32 {
	return binary.BigEndian.Uint32(bin)
}

func Ntohll(bin []byte) uint64 {
	return binary.BigEndian.Uint64(bin)
}

func Ntoh64(bin []byte) uint64 {
	return binary.BigEndian.Uint64(bin)
}

func Htons(bin []byte, v uint16) {
	binary.BigEndian.PutUint16(bin, v)
}

func Hton16(bin []byte, v uint16) {
	binary.BigEndian.PutUint16(bin, v)
}

func Htonl(bin []byte, v uint32) {
	binary.BigEndian.PutUint32(bin, v)
}

func Hton32(bin []byte, v uint32) {
	binary.BigEndian.PutUint32(bin, v)
}

func Htonll(bin []byte, v uint64) {
	binary.BigEndian.PutUint64(bin, v)
}

func Hton64(bin []byte, v uint64) {
	binary.BigEndian.PutUint64(bin, v)
}

func HtonsE(v uint16) []byte {
	bin := [2]byte{}

	Htons(bin[:], v)

	return bin[:]
}

func Hton16E(v uint16) []byte {
	bin := [2]byte{}

	Hton16(bin[:], v)

	return bin[:]
}

func HtonlE(v uint32) []byte {
	bin := [4]byte{}

	Htonl(bin[:], v)

	return bin[:]
}

func Hton32E(v uint32) []byte {
	bin := [4]byte{}

	Hton32(bin[:], v)

	return bin[:]
}

func HtonllE(v uint64) []byte {
	bin := [8]byte{}

	Htonll(bin[:], v)

	return bin[:]
}

func Hton64E(v uint64) []byte {
	bin := [8]byte{}

	Hton64(bin[:], v)

	return bin[:]
}

func NtoH16(v uint16) uint16 {
	return (v&0x00ff)<<8 | (v&0xff00)>>8
}

func NtoH32(v uint32) uint32 {
	return (v&0x000000ff)<<24 |
		(v&0x0000ff00)<<8 |
		(v&0x00ff0000)>>8 |
		(v&0xff000000)>>24
}

func NtoH64(v uint64) uint64 {
	return (v&0x00000000000000ff)<<56 |
		(v&0x000000000000ff00)<<40 |
		(v&0x0000000000ff0000)<<24 |
		(v&0x00000000ff000000)<<8 |
		(v&0x000000ff00000000)>>8 |
		(v&0x0000ff0000000000)>>24 |
		(v&0x00ff000000000000)>>40 |
		(v&0xff00000000000000)>>56
}

func HtoN16(v uint16) uint16 {
	return NtoH16(v)
}

func HtoN32(v uint32) uint32 {
	return NtoH32(v)
}

func HtoN64(v uint64) uint64 {
	return NtoH64(v)
}
