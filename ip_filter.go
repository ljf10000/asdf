package asdf

import (
	"strconv"
)

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
	if me.HasAny() {
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
	if me.HasAny() {
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
	} else if Empty != me.Ip {
		err := obj.Ip.FromString(me.Ip)
		if nil != err {
			return nil, ErrSprintf("parse ip-filter's ip error: %s", err)
		}

		obj.setType(IpResAddress)
	} else if Empty != me.Subnet {
		err := obj.Subnet.FromString(me.Subnet)
		if nil != err {
			return nil, ErrSprintf("parse ip-filter's subnet error: %s", err)
		}

		obj.setType(IpResSubnet)
	} else if Empty != me.Zone {
		err := obj.Zone.FromString(me.Zone)
		if nil != err {
			return nil, ErrSprintf("parse ip-filter's zone error: %s", err)
		}

		obj.setType(IpResZone)
	} else if len(me.List) > 0 {
		var ip IpAddress

		for k, v := range me.List {
			err := ip.FromString(v)
			if nil != err {
				return nil, ErrSprintf("parse ip-filter's list[%d] error: %s", k, err)
			}

			obj.Map[ip] = true
		}

		obj.setType(IpResMap)
	} else {
		obj.SetAny()
	}

	return obj, nil
}

/******************************************************************************/

type IpPairFilter struct {
	Hit  [2]int
	Pair [2]IpFilter
}

func (me *IpPairFilter) String() string {
	return "[" +
		strconv.Itoa(me.Hit[0]) + ":" + me.Pair[0].String() + ", " +
		strconv.Itoa(me.Hit[1]) + ":" + me.Pair[1].String() + "]"
}

func (me *IpPairFilter) IsMatch(a, b IpAddress) bool {
	if me.Pair[0].Match(a) && me.Pair[1].Match(b) {
		me.Hit[0]++

		return true
	} else if me.Pair[0].Match(b) && me.Pair[1].Match(a) {
		me.Hit[1]++

		return true
	} else {
		return false
	}
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
		s += ", " + strconv.Itoa(k) + ":" + v.String()
	}

	if len(s) > 0 {
		s = s[2:]
	}

	return "[" + s + "]"
}

func (me IpPairFilters) IsMatch(a, b IpAddress) (*IpPairFilter, bool) {
	for _, v := range me {
		if v.IsMatch(a, b) {
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
		s += ", " + strconv.Itoa(k) + ":" + v.String()
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
func MakeIpPairCache() IpPairCache {
	return IpPairCache{
		Cache: map[Ip2Tuple]bool{},
	}
}

type IpPairCache struct {
	// 0: 正向命中, sip/dip 命中 ip2tuple
	// 1: 反向命中, dip/sip 命中 ip2tuple
	Hit   [2]int
	Cache map[Ip2Tuple]bool
}

func (me *IpPairCache) String() string {
	s := Empty

	for k, _ := range me.Cache {
		s += ", (" + k.String() + ")"
	}

	return "[" +
		strconv.Itoa(me.Hit[0]) + ", " +
		strconv.Itoa(me.Hit[1]) + ", " +
		SkipString(s, 2) + "]"
}

func (me *IpPairCache) isMatch(sip, dip IpAddress) bool {
	tuple := Ip2Tuple{
		Sip: sip,
		Dip: dip,
	}

	return me.Cache[tuple]
}

func (me *IpPairCache) Add(sip, dip IpAddress) {
	tuple := Ip2Tuple{
		Sip: sip,
		Dip: dip,
	}

	me.Cache[tuple] = true
}

func (me *IpPairCache) IsMatch(sip, dip IpAddress) bool {
	if me.isMatch(sip, dip) {
		me.Hit[0]++

		return true
	} else if me.isMatch(dip, sip) {
		me.Hit[1]++

		return true
	} else {
		return false
	}
}
