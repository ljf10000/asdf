package asdf

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestIp(t *testing.T) {
	ip := IpAddress(0)
	bin := [4]byte{}

	ipstring := "192.168.0.1"
	ip.FromString(ipstring)

	binary.BigEndian.PutUint32(bin[:], uint32(ip))

	fmt.Printf("ipstring=%s, ip=%x, bin=%v"+Crlf, ipstring, uint32(ip), bin[:])
}
