package asdf

import (
	"time"
)

const DATE_SPLIT = "-"

type Date struct {
	year  uint16
	month byte
	day   byte
}

func MakeDate(t time.Time) Date {
	year, month, day := t.Date()

	return Date{
		year:  uint16(year),
		month: byte(month),
		day:   byte(day),
	}
}

func (me *Date) String() string {
	return Utoa16(me.year) +
		DATE_SPLIT + Utoa8(me.month) +
		DATE_SPLIT + Utoa8(me.day)
}

func (me *Date) Load(t time.Time) {
	*me = MakeDate(t)
}

func (me *Date) Unix() Time32 {
	t := time.Date(int(me.year), time.Month(me.month), int(me.day),
		0, 0, 0, 0, time.Local)

	return Time32(t.Unix())
}

func (me Date) AddDay(day int) Date {
	sec := me.Unix() + Time32(day)*24*3600
	t := time.Unix(int64(sec), 0)

	return MakeDate(t)
}

/*******************************************************************************/

type Datezone struct {
	Begin Date
	End   Date
}
