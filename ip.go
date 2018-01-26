package asdf

import (
	"encoding/binary"
	"fmt"
	"net"
	. "strconv"
)

type IpAddress uint32

func (me IpAddress) IsGood() bool {
	return true
}

func (me IpAddress) Int() int {
	return int(me)
}

func (me IpAddress) Eq(it interface{}) bool {
	return true
}

func (me IpAddress) ToString() string {
	bin := [4]byte{}

	binary.BigEndian.PutUint32(bin[:], uint32(me))

	return Itoa(int(bin[0])) + "." +
		Itoa(int(bin[1])) + "." +
		Itoa(int(bin[2])) + "." +
		Itoa(int(bin[3]))
}

func (me *IpAddress) FromString(s string) error {
	*me = IpAddressFromString(s)

	return nil
}

func IpAddressFromString(s string) IpAddress {
	ip := [4]byte{}
	fmt.Sscanf(s, "%d.%d.%d.%d", &ip[0], &ip[1], &ip[2], &ip[3])

	return IpAddress(binary.BigEndian.Uint32(ip[:]))
}

func GetLocalAddress() []string {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		return nil
	}

	var ipaddrs []string

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && nil != ipnet.IP.To4() {
			ipaddrs = append(ipaddrs, addr.String())
		}
	}

	return ipaddrs
}
