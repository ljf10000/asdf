package asdf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileName string

func (me FileName) String() string {
	return string(me)
}

func (me FileName) Append(buf []byte) error {
	return nil //todo
}

func (me FileName) Save(buf []byte) error {
	return ioutil.WriteFile(me.String(), buf, os.ModePerm)
}

func (me FileName) Load() ([]byte, error) {
	buf, err := ioutil.ReadFile(me.String())
	if nil == err {
		Log.Debug("load %s:%s\n", me, buf)
	}

	return buf, err
}

func (me FileName) SaveJson(obj interface{}) error {
	buf, err := json.Marshal(obj)
	if nil != err {
		return err
	}

	return me.Save(buf)
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
	} else if f, err := os.Stat(me.String()); nil != err {
		return false
	} else {
		return f.IsDir()
	}
}

func (me FileName) DirUnserialize(creator func() IUnserialize) error {
	return filepath.Walk(me.String(), func(path string, f os.FileInfo, err error) error {
		if nil == f {
			return err
		} else if f.IsDir() {
			return nil
		}

		obj := creator()

		if err := FileName(path).LoadJson(obj); nil != err {
			// skip error
		} else if err := obj.Unserialize(); nil != err {
			// skip error
		}

		return nil
	})
}
