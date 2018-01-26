package asdf

import (
	"crypto/sha256"
)

const DigestSize = 32

type IDigest interface {
	Digest(buf []byte) []byte
}

type DeftDigest struct{}

func (me *DeftDigest) Digest(buf []byte) []byte {
	hashDigest.Reset()
	hashDigest.Write(buf)

	return hashDigest.Sum(nil)
}

var hashDigest = sha256.New()
var DeftDigester = &DeftDigest{}
