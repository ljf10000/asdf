package asdf

import (
	"time"
)

const (
	DATE_SPLIT = "-"

	INVALID_YEAR  = YEAR_MAX + 1
	INVALID_MONTH = 0
	INVALID_DAY   = 0

	MONTH_PER_YEAR = 12
	DAY_PER_MONTH  = 31

	YEAR_MIN = 0
	YEAR_MAX = 65534

	MONTH_MIN = 1
	MONTH_MAX = MONTH_PER_YEAR

	DAY_MIN = 1
	DAY_MAX = 31
)

/******************************************************************************/

type Year uint16

func (me Year) String() string {
	return Utoa16(uint16(me))
}

func (me Year) IsGood() bool {
	return YEAR_MIN <= me && me <= YEAR_MAX
}

/******************************************************************************/

type Month byte

func (me Month) String() string {
	return Utoa8(byte(me))
}

func (me Month) IsGood() bool {
	return MONTH_MIN <= me && me <= MONTH_MAX
}

/******************************************************************************/

type Day byte

func (me Day) String() string {
	return Utoa8(byte(me))
}

func (me Day) IsGood() bool {
	return DAY_MIN <= me && me <= DAY_MAX
}

/******************************************************************************/

type Date struct {
	Year  `json:"year"`
	Month `json:"month"`
	Day   `json:"day"`
}

func MakeDate(t time.Time) Date {
	year, month, day := t.Date()

	return Date{
		Year:  Year(year),
		Month: Month(month),
		Day:   Day(day),
	}
}

func (me *Date) String() string {
	return me.Year.String() +
		DATE_SPLIT + me.Month.String() +
		DATE_SPLIT + me.Day.String()
}

func (me *Date) IsGood() bool {
	return me.Year.IsGood() && me.Month.IsGood() && me.Day.IsGood()
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
