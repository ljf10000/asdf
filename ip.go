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
	obj, ok := it.(IpAddress)
	if ok {
		return me == obj
	} else {
		return false
	}
}

func (me IpAddress) Compare(obj IpAddress) int {
	return CompareUint32(uint32(me), uint32(obj))
}

func (me IpAddress) String() string {
	return me.ToString()
}

func (me IpAddress) ToString() string {
	ip := MakeSlice(unsafe.Pointer(&me), 4, 4)

	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (me *IpAddress) FromString(s string) error {
	ip, err := IpAddressFromString(s)
	if nil != err {
		return err
	}

	*me = ip

	return nil
}

func IpAddressFromString(s string) (IpAddress, error) {
	ip := [4]byte{}

	n, err := fmt.Sscanf(s, "%d.%d.%d.%d", &ip[0], &ip[1], &ip[2], &ip[3])
	if nil != err {
		return 0, err
	} else if 4 != n {
		return 0, ErrBadIpAddress
	}

	return *(*IpAddress)(SlicePointer(ip[:])), nil
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
