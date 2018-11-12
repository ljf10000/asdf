package asdf

import (
	"fmt"
	"unsafe"
)

const SizeofIp2Tuple = 2 * SizeofInt32

var (
	zIp2Tuple = Ip2Tuple{}

	scIp2Tuple = NewSizeChecker("Ip2Tuple", unsafe.Sizeof(zIp2Tuple), SizeofIp2Tuple)
)

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
	return fmt.Sprintf("sip:%s, dip:%s",
		me.Sip,
		me.Dip)
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
	return fmt.Sprintf("sip:%s, dip:%s",
		me.Sip,
		me.Dip)
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
