package asdf

import (
	"crypto/sha256"
	"io"
)

const DigestSize = 32

type IDigest interface {
	io.Writer

	Digest(buf []byte) []byte
}

type DeftDigest struct {
	digest []byte
}

func (me *DeftDigest) NewDigest() []byte {
	var buf [DigestSize]byte

	return me.Digest(buf[:])
}

func (me *DeftDigest) Digest(buf []byte) []byte {
	me.Write(buf)

	return me.digest
}

func (me *DeftDigest) Write(buf []byte) (int, error) {
	hashDigest.Reset()
	hashDigest.Write(buf)

	me.digest = hashDigest.Sum(nil)

	return len(buf), nil
}

var hashDigest = sha256.New()
var DeftDigester = &DeftDigest{}
