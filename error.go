package asdf

import (
	"errors"
)

var Error 			= errors.New(Empty)

var ErrNoSupport	= errors.New("no support")
var ErrNoFound 		= errors.New("no found")
var ErrNoMatch		= errors.New("no match")
var ErrNoSpace 		= errors.New("no space")

var ErrBadObj 		= errors.New("bad obj")
var ErrNilObj 		= errors.New("nil obj")

var ErrTooShortBuffer 	= errors.New("too short buffer")
var ErrBadPktLen 		= errors.New("invalid packet length")
var ErrPktLenNoMatchBufferLen 	= errors.New("packet length not match buffer length")