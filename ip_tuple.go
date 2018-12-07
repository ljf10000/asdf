package asdf

import (
	"unsafe"
)

const (
	SizeofPortTuple = 2 * SizeofInt16
	SizeofIp2Tuple  = 2 * SizeofInt32
	SizeofIp4Tuple  = 3 * SizeofInt32
	SizeofIp5Tuple  = SizeofIp2Tuple + SizeofPortTuple + 4
	SizeofIp6Tuple  = SizeofIp4Tuple + 4
)

var (
	zPortTuple = PortTuple{}
	zIp2Tuple  = Ip2Tuple{}
	zIp4Tuple  = Ip4Tuple{}
	zIp5Tuple  = Ip5Tuple{}
	zIp6Tuple  = Ip6Tuple{}

	scPortTuple = NewSizeChecker("PortTuple", unsafe.Sizeof(zPortTuple), SizeofPortTuple)
	scIp2Tuple  = NewSizeChecker("Ip2Tuple", unsafe.Sizeof(zIp2Tuple), SizeofIp2Tuple)
	scIp4Tuple  = NewSizeChecker("Ip4Tuple", unsafe.Sizeof(zIp4Tuple), SizeofIp4Tuple)
	scIp5Tuple  = NewSizeChecker("Ip5Tuple", unsafe.Sizeof(zIp5Tuple), SizeofIp5Tuple)
	scIp6Tuple  = NewSizeChecker("Ip6Tuple", unsafe.Sizeof(zIp6Tuple), SizeofIp6Tuple)
)

/******************************************************************************/

func MakePortTuple(sport, dport uint16) PortTuple {
	return PortTuple{
		Sport: sport,
		Dport: dport,
	}
}

type PortTuple struct {
	Sport uint16 `json:"sport"`
	Dport uint16 `json:"dport"`
}

func (me *PortTuple) String() string {
	return "sport:" + Utoa16(me.Sport) +
		", dport:" + Utoa16(me.Dport)
}

func (me *PortTuple) Reverse() PortTuple {
	return PortTuple{
		Sport: me.Dport,
		Dport: me.Sport,
	}
}

func (me *PortTuple) Zero() {
	*me = zPortTuple
}

func (me *PortTuple) Compare(obj *PortTuple) int {
	if cmp := CompareUint16(me.Sport, obj.Sport); 0 != cmp {
		return cmp
	}

	if cmp := CompareUint16(me.Dport, obj.Dport); 0 != cmp {
		return cmp
	}

	return 0
}

func (me *PortTuple) Eq(obj *PortTuple) bool {
	return me.Sport == obj.Sport && me.Dport == obj.Dport
}

func (me *PortTuple) Slice() []byte {
	return StructSlice(unsafe.Pointer(me), SizeofPortTuple)
}

func (me *PortTuple) Bkdr() Bkdr {
	return GenBkdr(me.Slice())
}

func (me *PortTuple) Index() int {
	return int(me.Bkdr())
}

/******************************************************************************/

func MakeIp2Tuple(sip, dip IpAddress) Ip2Tuple {
	return Ip2Tuple{
		Sip: sip,
		Dip: dip,
	}
}

type Ip2Tuple struct {
	Sip IpAddress // network sort
	Dip IpAddress // network sort
}

func (me *Ip2Tuple) String() string {
	return "sip:" + me.Sip.String() +
		", dip:" + me.Dip.String()
}

func (me *Ip2Tuple) Reverse() Ip2Tuple {
	return Ip2Tuple{
		Sip: me.Dip,
		Dip: me.Sip,
	}
}

func (me *Ip2Tuple) Zero() {
	*me = zIp2Tuple
}

func (me *Ip2Tuple) IsGood() bool {
	return me.Sip > 0 && me.Dip > 0
}

func (me *Ip2Tuple) Compare(obj *Ip2Tuple) int {
	if cmp := me.Sip.Compare(obj.Sip); 0 != cmp {
		return cmp
	}

	if cmp := me.Dip.Compare(obj.Dip); 0 != cmp {
		return cmp
	}

	return 0
}

func (me *Ip2Tuple) Eq(obj *Ip2Tuple) bool {
	return me.Sip == obj.Sip && me.Dip == obj.Dip
}

func (me *Ip2Tuple) Slice() []byte {
	return StructSlice(unsafe.Pointer(me), SizeofIp2Tuple)
}

func (me *Ip2Tuple) Bkdr() Bkdr {
	return GenBkdr(me.Slice())
}

func (me *Ip2Tuple) Index() int {
	return int(me.Bkdr())
}

type Ip2TupleStr struct {
	Sip string `json:"sip"`
	Dip string `json:"dip"`
}

func (me *Ip2TupleStr) String() string {
	return "sip:" + me.Sip +
		", dip:" + me.Dip
}

func (me *Ip2TupleStr) Atoi() (Ip2Tuple, error) {
	tuple := Ip2Tuple{}

	err := tuple.Sip.FromString(me.Sip)
	if nil != err {
		return zIp2Tuple, err
	}

	err = tuple.Dip.FromString(me.Dip)
	if nil != err {
		return zIp2Tuple, err
	}

	return tuple, nil
}

func (me *Ip2Tuple) Itoa() Ip2TupleStr {
	return Ip2TupleStr{
		Sip: me.Sip.String(),
		Dip: me.Dip.String(),
	}
}

/******************************************************************************/

func MakeIp4Tuple(sip, dip IpAddress, id uint16, proto IpProto) Ip4Tuple {
	return Ip4Tuple{
		Ip2Tuple: Ip2Tuple{
			Sip: sip,
			Dip: dip,
		},
		Id:    id,
		Proto: proto,
	}
}

// network sort
type Ip4Tuple struct {
	Ip2Tuple

	Id    uint16
	Proto IpProto
	_     byte
}

func (me *Ip4Tuple) String() string {
	return me.Ip2Tuple.String() +
		", id:" + Utoa16(me.Id) +
		", proto:" + Utoa8(byte(me.Proto))
}

func (me *Ip4Tuple) Reverse() Ip4Tuple {
	return Ip4Tuple{
		Ip2Tuple: me.Ip2Tuple.Reverse(),
		Id:       me.Id,
		Proto:    me.Proto,
	}
}

func (me *Ip4Tuple) Zero() {
	*me = zIp4Tuple
}

func (me *Ip4Tuple) IsGood() bool {
	return me.Proto > 0
}

func (me *Ip4Tuple) Compare(obj *Ip4Tuple) int {
	if cmp := me.Ip2Tuple.Compare(&obj.Ip2Tuple); 0 != cmp {
		return cmp
	}

	if cmp := CompareUint16(me.Id, obj.Id); 0 != cmp {
		return cmp
	}

	if cmp := me.Proto.Compare(obj.Proto); 0 != cmp {
		return cmp
	}

	return 0
}

func (me *Ip4Tuple) Eq(obj *Ip4Tuple) bool {
	return me.Ip2Tuple.Eq(&obj.Ip2Tuple) && me.Id == obj.Id && me.Proto == obj.Proto
}

func (me *Ip4Tuple) Slice() []byte {
	return StructSlice(unsafe.Pointer(me), SizeofIp4Tuple)
}

func (me *Ip4Tuple) Bkdr() Bkdr {
	return GenBkdr(me.Slice())
}

func (me *Ip4Tuple) Index() int {
	return int(me.Bkdr())
}

type Ip4TupleStr struct {
	Ip2TupleStr

	Id    uint16  `json:"id"`
	Proto IpProto `json:"proto"`
	_     byte
}

func (me *Ip4TupleStr) String() string {
	return me.Ip2TupleStr.String() +
		", id:" + Utoa16(me.Id) +
		", proto:" + Utoa8(byte(me.Proto))
}

func (me *Ip4TupleStr) Atoi() (Ip4Tuple, error) {
	tuple := Ip4Tuple{
		Id:    me.Id,
		Proto: me.Proto,
	}

	var err error

	tuple.Ip2Tuple, err = me.Ip2TupleStr.Atoi()
	if nil != err {
		return zIp4Tuple, err
	}

	return tuple, nil
}

func (me *Ip4Tuple) Itoa() Ip4TupleStr {
	return Ip4TupleStr{
		Ip2TupleStr: me.Ip2Tuple.Itoa(),
		Id:          me.Id,
		Proto:       me.Proto,
	}
}

/******************************************************************************/

func MakeIp5Tuple(sip, dip IpAddress, sport, dport uint16, proto IpProto) Ip5Tuple {
	return Ip5Tuple{
		Ip2Tuple: Ip2Tuple{
			Sip: sip,
			Dip: dip,
		},
		PortTuple: PortTuple{
			Sport: sport,
			Dport: dport,
		},
		Proto: proto,
	}
}

type IIp5Tuple interface {
	Ip5Tuple() Ip5Tuple
}

type Ip5Tuple struct {
	Ip2Tuple
	PortTuple

	Proto IpProto
	_     [3]byte
}

func (me *Ip5Tuple) String() string {
	return me.Ip2Tuple.String() +
		", " + me.PortTuple.String() +
		", proto:" + Utoa8(byte(me.Proto))
}

func (me *Ip5Tuple) Reverse() Ip5Tuple {
	return Ip5Tuple{
		Ip2Tuple:  me.Ip2Tuple.Reverse(),
		PortTuple: me.PortTuple.Reverse(),
		Proto:     me.Proto,
	}
}

func (me *Ip5Tuple) Zero() {
	*me = zIp5Tuple
}

func (me *Ip5Tuple) IsGood() bool {
	return me.Proto > 0
}

func (me *Ip5Tuple) Compare(obj *Ip5Tuple) int {
	if cmp := me.Ip2Tuple.Compare(&obj.Ip2Tuple); 0 != cmp {
		return cmp
	}

	if cmp := me.PortTuple.Compare(&obj.PortTuple); 0 != cmp {
		return cmp
	}

	return me.Proto.Compare(obj.Proto)
}

func (me *Ip5Tuple) Eq(obj *Ip5Tuple) bool {
	return me.Proto == obj.Proto &&
		me.Ip2Tuple.Eq(&obj.Ip2Tuple) &&
		me.PortTuple.Eq(&obj.PortTuple)
}

func (me *Ip5Tuple) Slice() []byte {
	return StructSlice(unsafe.Pointer(me), SizeofIp5Tuple)
}

func (me *Ip5Tuple) Bkdr() Bkdr {
	return GenBkdr(me.Slice())
}

func (me *Ip5Tuple) Index() int {
	return int(me.Bkdr())
}

type Ip5TupleStr struct {
	Ip2TupleStr
	PortTuple

	Proto IpProto `json:"proto"`
	_     [3]byte
}

func (me *Ip5TupleStr) String() string {
	return me.Ip2TupleStr.String() +
		", " + me.PortTuple.String() +
		", proto:" + Utoa8(byte(me.Proto))
}

func (me *Ip5TupleStr) Atoi() (Ip5Tuple, error) {
	tuple := Ip5Tuple{
		PortTuple: me.PortTuple,
		Proto:     me.Proto,
	}

	var err error

	tuple.Ip2Tuple, err = me.Ip2TupleStr.Atoi()
	if nil != err {
		return zIp5Tuple, err
	}

	return tuple, nil
}

func (me *Ip5Tuple) Itoa() Ip5TupleStr {
	return Ip5TupleStr{
		Ip2TupleStr: me.Ip2Tuple.Itoa(),
		PortTuple:   me.PortTuple,
		Proto:       me.Proto,
	}
}

type Ip5Tuples map[Ip5Tuple]uint64

func (me Ip5Tuples) AddEx(tuple Ip5Tuple) uint64 {
	count, ok := me[tuple]
	if ok {
		me[tuple] = count + 1

		return count + 1
	}

	r := tuple.Reverse()
	count, ok = me[r]
	if ok {
		me[r] = count + 1

		return count + 1
	}

	me[tuple] = 1

	return 1
}

func (me Ip5Tuples) Add(tuple Ip5Tuple) uint64 {
	count, ok := me[tuple]
	if ok {
		me[tuple] = count + 1

		return count + 1
	} else {
		me[tuple] = 1

		return 1
	}
}

func (me Ip5Tuples) ToList() Ip5TupleList {
	count := len(me)

	list := make(Ip5TupleList, 0, count)

	for k, v := range me {
		st := Ip5TupleStat{
			Ip5Tuple: k,
			Count:    v,
		}

		list = append(list, st.String())
	}

	return list
}

type Ip5TupleStat struct {
	Ip5Tuple

	Count uint64 `json:"count"`
}

func (me *Ip5TupleStat) String() string {
	return me.Ip5Tuple.String() +
		", count:" + Utoa64(me.Count)
}

type Ip5TupleList []string

/******************************************************************************/

func MakeIp6Tuple(sip, dip IpAddress, id uint16, proto IpProto, bodySize, offset uint16) Ip6Tuple {
	return Ip6Tuple{
		Ip4Tuple: Ip4Tuple{
			Ip2Tuple: Ip2Tuple{
				Sip: sip,
				Dip: dip,
			},
			Id:    id,
			Proto: proto,
		},
		IpBodySize: bodySize,
		Offset:     offset,
	}
}

type Ip6Tuple struct {
	Ip4Tuple

	IpBodySize uint16
	Offset     uint16
}

func (me *Ip6Tuple) String() string {
	return me.Ip4Tuple.String() +
		", offset" + Utoa16(me.Offset) +
		", size:" + Utoa16(me.IpBodySize)
}

func (me *Ip6Tuple) Reverse() Ip6Tuple {
	return Ip6Tuple{
		Ip4Tuple:   me.Ip4Tuple.Reverse(),
		IpBodySize: me.IpBodySize,
		Offset:     me.Offset,
	}
}

func (me *Ip6Tuple) Zero() {
	me.Ip4Tuple.Zero()

	me.IpBodySize = 0
	me.Offset = 0
}

func (me *Ip6Tuple) IsGood() bool {
	return me.Ip4Tuple.IsGood()
}

func (me *Ip6Tuple) Compare(obj *Ip6Tuple) int {
	if cmp := me.Ip4Tuple.Compare(&obj.Ip4Tuple); 0 != cmp {
		return cmp
	}

	if cmp := CompareUint16(me.IpBodySize, obj.IpBodySize); 0 != cmp {
		return cmp
	}

	if cmp := CompareUint16(me.Offset, obj.Offset); 0 != cmp {
		return cmp
	}

	return 0
}

func (me *Ip6Tuple) Eq(obj *Ip6Tuple) bool {
	return me.IpBodySize == obj.IpBodySize &&
		me.Offset == obj.Offset &&
		me.Ip4Tuple.Eq(&obj.Ip4Tuple)
}

func (me *Ip6Tuple) Slice() []byte {
	return StructSlice(unsafe.Pointer(me), SizeofIp6Tuple)
}

func (me *Ip6Tuple) Bkdr() Bkdr {
	return GenBkdr(me.Slice())
}

func (me *Ip6Tuple) Index() int {
	return int(me.Bkdr())
}

type Ip6TupleStr struct {
	Ip4TupleStr

	IpBodySize uint16 `json:"bodysize"`
	Offset     uint16 `json:"offset"`
}

func (me *Ip6TupleStr) String() string {
	return me.Ip4TupleStr.String() +
		", offset" + Utoa16(me.Offset) +
		", size:" + Utoa16(me.IpBodySize)
}

func (me *Ip6TupleStr) Atoi() (Ip6Tuple, error) {
	tuple := Ip6Tuple{
		IpBodySize: me.IpBodySize,
		Offset:     me.Offset,
	}

	var err error

	tuple.Ip4Tuple, err = me.Ip4TupleStr.Atoi()
	if nil != err {
		return zIp6Tuple, err
	}

	return tuple, nil
}

func (me *Ip6Tuple) Itoa() Ip6TupleStr {
	return Ip6TupleStr{
		Ip4TupleStr: me.Ip4Tuple.Itoa(),
		IpBodySize:  me.IpBodySize,
		Offset:      me.Offset,
	}
}
