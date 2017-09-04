package asdf

import (
	"time"
)

type Time32 uint32

func (me Time32) Unix() time.Time {
	return time.Unix(int64(me), 0)
}

func NowTime32() Time32 {
	return Time32(time.Now().Unix())
}
