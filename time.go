package asdf

import (
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type Time32 uint32

func (me Time32) Unix() time.Time {
	return time.Unix(int64(me), 0)
}

func (me *Time32) Read(s string) error {
	tm, err := time.Parse(TimeFormat, s)
	if nil != err {
		return err
	}

	*me = Time32(tm.Unix())
	return nil
}

func NowTime32() Time32 {
	return Time32(time.Now().Unix())
}
