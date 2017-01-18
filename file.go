package asdf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type FileName string

func (me FileName) Load() ([]byte, error) {
	buf, err := ioutil.ReadFile(string(me))
	if nil == err {
		Log.Debug("load %s:%s\n", me, buf)
	}

	return buf, err
}

func (me FileName) LoadJson(obj interface{}) error {
	if buf, err := me.Load(); nil != err {
		return err
	} else {
		return json.Unmarshal(buf, obj)
	}
}

func (me FileName) DirExist() bool {
	if Empty == me {
		return false
	} else if f, err := os.Stat(string(me)); nil != err {
		return false
	} else {
		return f.IsDir()
	}
}
