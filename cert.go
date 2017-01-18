package asdf

import (
	"encoding/base64"
)

//
// base64 encoded
//
type BuildinCert struct {
	Cert string
	Key  string
	Ca   string
}

//
// after base64 decoded
//
type RawCert struct {
	Cert []byte
	Key  []byte
	Ca   []byte
}

func (me *RawCert) FromBuildin(buildin *BuildinCert) {
	me.Cert, _ = base64.StdEncoding.DecodeString(buildin.Cert)
	me.Key, _ = base64.StdEncoding.DecodeString(buildin.Key)
	me.Ca, _ = base64.StdEncoding.DecodeString(buildin.Ca)
}
