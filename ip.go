package asdf

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"unsafe"
)

func IpAddressFromString(s string) (IpAddress, error) {
	ip := [4]byte{}

	n, err := fmt.Sscanf(s, "%d.%d.%d.%d", &ip[0], &ip[1], &ip[2], &ip[3])
	if 4 != n {
		return 0, ErrSprintf("ip address(%s) parse error: n(%d) not 4", s, n)
	} else if nil != err {
		return 0, err
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
const IP_ANY IpAddress = 0

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

func (me IpAddress) MaskEx(length IpMaskLen) IpAddress {
	return IpAddress((1<<length - 1) << (32 - length))
}

func (me IpAddress) Mask(length IpMaskLen) IpAddress {
	switch length {
	case 8:
		return 0xff000000
	case 16:
		return 0xffff0000
	case 24:
		return 0xffffff00
	case 32:
		return 0xffffffff
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
const IpSubnetSplit = "/"

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

const IpZoneSplit = "-"

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

func (me IpZone) Match(v IpZone) bool {
	return v.Begin.InZone(me) || v.End.InZone(me)
}

func (me IpZone) Intersect(v IpZone) IpZone {
	if v.Begin.InZone(me) {
		if v.End.InZone(me) {
			// |--------- me ---------|
			//     |----- v -----|
			return v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return IpZone{
				Begin: v.Begin,
				End:   me.End,
			}
		}
	} else {
		if v.End.InZone(me) {
			//     |--------- me ---------|
			// |----- v -----|
			return IpZone{
				Begin: me.Begin,
				End:   v.End,
			}
		} else {
			//                  |--------- me ---------|
			// |----- v -----|              or              |----- v -----|
			return IpZone{}
		}
	}
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

/******************************************************************************/

const (
	IpResAddress IpResType = 0
	IpResSubnet  IpResType = 1
	IpResZone    IpResType = 2
	IpResMap     IpResType = 3
	IpResAny     IpResType = 4
	IpResEnd     IpResType = 5
)

type IpResType byte

func (me IpResType) TypeMask() IpResType {
	return 1 << me
}

func NewIpFilter() *IpFilter {
	return &IpFilter{
		Map: map[IpAddress]bool{},
	}
}

type IpFilter struct {
	TypeMask IpResType
	Ip       IpAddress
	Subnet   IpSubnet
	Zone     IpZone
	Map      map[IpAddress]bool
}

func (me *IpFilter) String() string {
	if me.hasType(IpResAny) {
		return "any: true"
	} else {
		s := Empty

		if me.hasType(IpResAddress) {
			s += ", ip: " + me.Ip.String()
		}

		if me.hasType(IpResSubnet) {
			s += ", subnet: " + me.Subnet.String()
		}

		if me.hasType(IpResZone) {
			s += ", zone: " + me.Zone.String()
		}

		if me.hasType(IpResMap) {
			ss := Empty

			for k, _ := range me.Map {
				ss += ", " + k.String()
			}

			if len(ss) > 0 {
				s += ", map:[" + ss[2:] + "]"
			}
		}

		if len(s) > 0 {
			s = s[2:]
		}

		return s
	}
}

func (me *IpFilter) hasType(Type IpResType) bool {
	mask := Type.TypeMask()

	return mask == (mask & me.TypeMask)
}

func (me *IpFilter) setType(Type IpResType) {
	me.TypeMask |= Type.TypeMask()
}

func (me *IpFilter) HasAny() bool {
	return me.hasType(IpResAny)
}

func (me *IpFilter) SetAny() {
	me.setType(IpResAny)
}

func (me *IpFilter) SetIp(ip IpAddress) {
	me.setType(IpResAddress)
	me.Ip = ip
}

func (me *IpFilter) SetSubnet(subnet IpSubnet) {
	me.setType(IpResSubnet)
	me.Subnet = subnet
}

func (me *IpFilter) SetZone(zone IpZone) {
	me.setType(IpResZone)
	me.Zone = zone
}

func (me *IpFilter) AddIp(ip IpAddress) {
	me.setType(IpResMap)
	me.Map[ip] = true
}

func (me *IpFilter) Match(ip IpAddress) bool {
	if me.hasType(IpResAny) {
		return true
	}

	if me.hasType(IpResAddress) && ip == me.Ip {
		return true
	}

	if me.hasType(IpResSubnet) && ip.InSubnet(me.Subnet) {
		return true
	}

	if me.hasType(IpResZone) && ip.InZone(me.Zone) {
		return true
	}

	if me.hasType(IpResMap) {
		if _, ok := me.Map[ip]; ok {
			return true
		}
	}

	return false
}

func (me *IpFilter) Itoa() *IpFilterStr {
	obj := &IpFilterStr{}

	if me.HasAny() {
		obj.Any = true
	} else {
		if me.hasType(IpResAddress) {
			obj.Ip = me.Ip.String()
		}

		if me.hasType(IpResSubnet) {
			obj.Subnet = me.Subnet.String()
		}

		if me.hasType(IpResZone) {
			obj.Zone = me.Zone.String()
		}

		if me.hasType(IpResMap) {
			for k, _ := range me.Map {
				obj.List = append(obj.List, k.String())
			}
		}
	}

	return obj
}

type IpFilterStr struct {
	Any    bool     `json:"any"`
	Ip     string   `json:"ip"`
	Subnet string   `json:"subnet"`
	Zone   string   `json:"zone"`
	List   []string `json:"list"`
}

func (me *IpFilterStr) String() string {
	if me.Any {
		return "any: true"
	} else {
		s := Empty

		if Empty != me.Ip {
			s += ", ip: " + me.Ip
		}

		if Empty != me.Subnet {
			s += ", subnet: " + me.Subnet
		}

		if Empty != me.Zone {
			s += ", zone: " + me.Zone
		}

		if len(me.List) > 0 {
			ss := Empty

			for _, v := range me.List {
				ss += ", " + v
			}

			if len(ss) > 0 {
				s += ", list:[" + ss[2:] + "]"
			}
		}

		if len(s) > 0 {
			s = s[2:]
		}

		return s
	}
}

func (me *IpFilterStr) Atoi() (*IpFilter, error) {
	obj := NewIpFilter()

	if me.Any {
		obj.SetAny()
	} else {
		if Empty != me.Ip {
			err := obj.Ip.FromString(me.Ip)
			if nil != err {
				return nil, ErrSprintf("parse ip-filter's ip error: %s", err)
			}

			obj.setType(IpResAddress)
		}

		if Empty != me.Subnet {
			err := obj.Subnet.FromString(me.Subnet)
			if nil != err {
				return nil, ErrSprintf("parse ip-filter's subnet error: %s", err)
			}

			obj.setType(IpResSubnet)
		}

		if Empty != me.Zone {
			err := obj.Zone.FromString(me.Zone)
			if nil != err {
				return nil, ErrSprintf("parse ip-filter's zone error: %s", err)
			}

			obj.setType(IpResZone)
		}

		if len(me.List) > 0 {
			var ip IpAddress

			for k, v := range me.List {
				err := ip.FromString(v)
				if nil != err {
					return nil, ErrSprintf("parse ip-filter's list[%d] error: %s", k, err)
				}

				obj.Map[ip] = true
			}

			obj.setType(IpResMap)
		}
	}

	return obj, nil
}

/******************************************************************************/

type IpPairFilter struct {
	Pair [2]IpFilter
}

func (me *IpPairFilter) String() string {
	return "[" + me.Pair[0].String() + ", " + me.Pair[1].String() + "]"
}

func (me *IpPairFilter) Match(a, b IpAddress) bool {
	return (me.Pair[0].Match(a) && me.Pair[1].Match(b)) ||
		(me.Pair[0].Match(b) && me.Pair[1].Match(a))
}

func (me *IpPairFilter) Itoa() *IpPairFilterStr {
	return &IpPairFilterStr{
		Pair: [2]IpFilterStr{*me.Pair[0].Itoa(), *me.Pair[1].Itoa()},
	}
}

type IpPairFilterStr struct {
	Pair [2]IpFilterStr `json:"pair"`
}

func (me *IpPairFilterStr) String() string {
	return "[" + me.Pair[0].String() + ", " + me.Pair[1].String() + "]"
}

func (me *IpPairFilterStr) Atoi() (*IpPairFilter, error) {
	var err error
	var Pair [2]*IpFilter

	for i := 0; i < 2; i++ {
		Pair[i], err = me.Pair[i].Atoi()
		if nil != err {
			return nil, err
		}
	}

	return &IpPairFilter{
		Pair: [2]IpFilter{*Pair[0], *Pair[1]},
	}, nil
}

/******************************************************************************/

type IpPairFilters []*IpPairFilter

func (me IpPairFilters) String() string {
	s := Empty

	for k, v := range me {
		s += fmt.Sprintf(", %d:%s", k, v)
	}

	if len(s) > 0 {
		s = s[2:]
	}

	return "[" + s + "]"
}

func (me IpPairFilters) Match(a, b IpAddress) (*IpPairFilter, bool) {
	for _, v := range me {
		if v.Match(a, b) {
			return v, true
		}
	}

	return nil, false
}

func (me IpPairFilters) Itoa() IpPairFilterStrs {
	var obj IpPairFilterStrs

	for _, v := range me {
		obj = append(obj, v.Itoa())
	}

	return obj
}

type IpPairFilterStrs []*IpPairFilterStr

func (me IpPairFilterStrs) String() string {
	s := Empty

	for k, v := range me {
		s += fmt.Sprintf(", %d:%s", k, v)
	}

	if len(s) > 0 {
		s = s[2:]
	}

	return "[" + s + "]"
}

func (me IpPairFilterStrs) Atoi() (IpPairFilters, error) {
	var obj IpPairFilters

	for _, v := range me {
		tmp, err := v.Atoi()
		if nil != err {
			return nil, err
		}

		obj = append(obj, tmp)
	}

	return obj, nil
}

/******************************************************************************/

type IpPairCache map[Ip2Tuple]bool

func (me IpPairCache) isMatch(sip, dip IpAddress) bool {
	tuple := Ip2Tuple{
		Sip: sip,
		Dip: dip,
	}

	_, ok := me[tuple]

	return ok
}

func (me IpPairCache) AddCache(sip, dip IpAddress) {
	tuple := Ip2Tuple{
		Sip: sip,
		Dip: dip,
	}

	me[tuple] = true
}

func (me IpPairCache) IsMatch(sip, dip IpAddress) bool {
	return me.isMatch(sip, dip) || me.isMatch(dip, sip)
}
