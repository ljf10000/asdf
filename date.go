package asdf

import (
	"time"
)

const DATE_SPLIT = "-"

type Date struct {
	Year  uint16
	Month byte
	Day   byte
}

func MakeDate(t time.Time) Date {
	year, month, day := t.Date()

	return Date{
		Year:  uint16(year),
		Month: byte(month),
		Day:   byte(day),
	}
}

func (me *Date) String() string {
	return Utoa16(me.Year) +
		DATE_SPLIT + Utoa8(me.Month) +
		DATE_SPLIT + Utoa8(me.Day)
}

func (me *Date) Load(t time.Time) {
	*me = MakeDate(t)
}

func (me *Date) Unix() Time32 {
	t := time.Date(int(me.Year), time.Month(me.Month), int(me.Day),
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
