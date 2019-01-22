package asdf

import (
	"net"
	"strconv"
	"strings"
	"unsafe"
)

/******************************************************************************/
const (
	IP_ANY IpAddress = 0

	IpAddressSplit = "."
	IpSubnetSplit  = "/"
	IpZoneSplit    = "-"
)

func IpAddressFromString(s string) (IpAddress, error) {
	ip := [4]byte{}

	split := strings.Split(s, IpAddressSplit)
	n := len(split)
	if 4 != n {
		return 0, ErrSprintf("ip address(%s) parse error: n(%d) not 4", s, n)
	}

	for i := 0; i < 4; i++ {
		v, err := strconv.Atoi(split[i])
		if nil != err {
			return 0, ErrSprintf("ip address(%s) parse error: %s", s, err)
		}

		ip[i] = byte(v)
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

/******************************************************************************/

type IpAddress uint32

func (me IpAddress) IsGood() bool {
	return IP_ANY != me
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

	return strconv.Itoa(int(ip[0])) + IpAddressSplit +
		strconv.Itoa(int(ip[1])) + IpAddressSplit +
		strconv.Itoa(int(ip[2])) + IpAddressSplit +
		strconv.Itoa(int(ip[3]))
}

func (me *IpAddress) FromString(s string) error {
	ip, err := IpAddressFromString(s)
	if nil != err {
		return err
	}

	*me = ip

	return nil
}

func (me IpAddress) MaskEx(length IpMaskLen) IpAddress {
	return IpAddress((1<<length - 1) << (32 - length))
}

const (
	IpMaskA   IpAddress = 0xff000000
	IpMaskB   IpAddress = 0xffff0000
	IpMaskC   IpAddress = 0xffffff00
	IpMaskAll IpAddress = 0xffffffff
)

func (me IpAddress) Mask(length IpMaskLen) IpAddress {
	switch length {
	case 8:
		return IpMaskA
	case 16:
		return IpMaskB
	case 24:
		return IpMaskC
	case 32:
		return IpMaskAll
	default:
		return me.MaskEx(length)
	}
}

func (me IpAddress) Network(mask IpAddress) IpAddress {
	return me & mask
}

func (me IpAddress) Host(mask IpAddress) IpAddress {
	return me & ^mask
}

func (me IpAddress) Match(ip, mask IpAddress) bool {
	return me.Network(mask) == ip.Network(mask)
}

func (me IpAddress) inRange(a, b IpAddress) bool {
	return a <= me && me <= b
}

func (me IpAddress) InRange(a, b IpAddress) bool {
	if a < b {
		return me.inRange(a, b)
	} else {
		return me.inRange(b, a)
	}
}

func (me IpAddress) InZone(z IpZone) bool {
	return me.inRange(z.Begin, z.End)
}

func (me IpAddress) InSubnet(v IpSubnet) bool {
	mask := v.Mask()

	return me.Network(mask) == v.Network()
}

/******************************************************************************/

type IpMaskLen byte

func (me IpMaskLen) IsGood() bool {
	return me > 0 && me <= 32
}

func (me IpMaskLen) String() string {
	return strconv.Itoa(int(me))
}

func (me *IpMaskLen) FromString(s string) error {
	v, err := strconv.Atoi(s)
	if nil != err {
		return err
	}

	*me = IpMaskLen(v)

	return nil
}

/******************************************************************************/

type IpSubnet struct {
	Ip  IpAddress `json:"ip"`
	Len IpMaskLen `json:"len"`
}

func (me IpSubnet) String() string {
	return me.Ip.String() + IpSubnetSplit + me.Len.String()
}

func (me *IpSubnet) FromString(s string) error {
	// n, err := fmt.Sscanf(s, "%s/%d", &sIp, &Len)'
	// fuck, n always 1, why ?

	split := strings.Split(s, IpSubnetSplit)
	n := len(split)
	if 2 != n {
		return ErrSprintf("ip subnet(%s) parse error: n(%d) not 2", s, n)
	}

	var Len IpMaskLen
	var Ip IpAddress

	err := Ip.FromString(split[0])
	if nil != err {
		return err
	}

	err = Len.FromString(split[1])
	if nil != err {
		return err
	}

	me.Ip = Ip
	me.Len = Len

	return nil
}

func (me IpSubnet) IsGood() bool {
	return me.Ip.IsGood() || me.Len.IsGood()
}

func (me *IpSubnet) Zero() {
	me.Ip = 0
	me.Len = 0
}

func (me IpSubnet) Mask() IpAddress {
	return me.Ip.Mask(me.Len)
}

func (me IpSubnet) Network() IpAddress {
	return me.Ip.Network(me.Mask())
}

func (me IpSubnet) Host() IpAddress {
	return me.Ip.Host(me.Mask())
}

func (me IpSubnet) Include(v IpSubnet) bool {
	// |--------- me ---------|
	//      |----- v -----|
	mask := me.Mask()

	return me.Ip.Network(mask) == v.Ip.Network(mask)
}

/******************************************************************************/

type IpZone struct {
	Begin IpAddress `json:"begin"`
	End   IpAddress `json:"end"`
}

func (me IpZone) String() string {
	return me.Begin.String() + IpZoneSplit + me.End.String()
}

func (me *IpZone) FromString(s string) error {
	split := strings.Split(s, IpZoneSplit)
	n := len(split)
	if 2 != n {
		return ErrSprintf("ip zone(%s) parse error: n(%d) not 2", s, n)
	}

	var Begin, End IpAddress

	err := Begin.FromString(split[0])
	if nil != err {
		return err
	}

	err = End.FromString(split[1])
	if nil != err {
		return err
	}

	me.Begin = Begin
	me.End = End

	return nil
}

func (me IpZone) IsGood() bool {
	return me.Begin.IsGood() || me.End.IsGood()
}

func (me *IpZone) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me IpZone) Include(z IpZone) bool {
	// |--------- me ---------|
	//      |----- z -----|
	return me.Begin <= z.Begin && me.End >= z.End
}

func (me IpZone) Compare(v IpZone) int {
	if me.End < v.Begin {
		// |--------- me ---------|
		//                            |----- v -----|
		return -1
	} else if me.Begin > v.End {
		//                  |--------- me ---------|
		// |----- v -----|
		return 1
	} else {
		//            |--------- me ---------|
		// |----- v -----|
		//                 |----- v -----|
		//                                 |----- v -----|
		return 0
	}
}

func (me IpZone) Match(v IpZone) bool {
	return 0 == me.Compare(v)
}

func (me IpZone) Intersect(v IpZone) IpZone {
	if 0 != me.Compare(v) {
		return IpZone{}
	}

	// get max begin
	begin := me.Begin
	if me.Begin < v.Begin {
		begin = v.Begin
	}

	// get min end
	end := me.End
	if me.End > v.End {
		end = v.End
	}

	return IpZone{
		Begin: begin,
		End:   end,
	}
}
