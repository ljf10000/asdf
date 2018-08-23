package asdf

import (
	"fmt"
	"net"
	"strings"
	"unsafe"
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

func (me IpAddress) String() string {
	return me.ToString()
}

func (me IpAddress) ToString() string {
	address := uintptr(unsafe.Pointer(&me))
	ip := MakeSlice(address, 4, 4)

	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (me *IpAddress) FromString(s string) error {
	*me = IpAddressFromString(s)

	return nil
}

func IpAddressFromString(s string) IpAddress {
	ip := [4]byte{}

	fmt.Sscanf(s, "%d.%d.%d.%d", &ip[0], &ip[1], &ip[2], &ip[3])

	return *(*IpAddress)(SlicePointer(ip[:]))
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
			arr := strings.Split(addr.String(), "/")
			ipaddrs = append(ipaddrs, arr[0])
		}
	}

	return ipaddrs
}
