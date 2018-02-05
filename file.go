package asdf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	FilePermExec   = 0755
	FilePermNormal = 0644
)

type FileName string

func (me FileName) String() string {
	return string(me)
}

func (me FileName) Abs() FileName {
	return FileName(CurrentDirFile(me.String()))
}

func (me FileName) Append(buf []byte) error {
	return nil //todo
}

func (me FileName) Save(buf []byte) error {
	err := ioutil.WriteFile(me.String(), buf, FilePermNormal)
	if nil != err {
		Log.Info("save %s error:%v\n", me, err)
	}

	return err
}

func (me FileName) Delete() error {
	err := os.Remove(me.String())
	if nil != err {
		Log.Info("delete %s error:%v\n", me, err)
	}

	return err
}

func (me FileName) Touch(Time Time32) error {
	if 0 == Time {
		Time = NowTime32()
	}

	tm := time.Unix(int64(Time), 0)

	err := os.Chtimes(me.String(), tm, tm)
	if nil != err {
		Log.Info("change %s time error:%v\n", me, err)
	}

	return err
}

func (me FileName) Saves(texts []string, crlf bool) error {
	f, err := os.Create(me.String())
	if nil != err {
		Log.Info("create %s error:%v\n", me, err)

		return err
	}
	defer f.Close()

	for _, text := range texts {
		_, err := f.WriteString(text)
		if nil != err {
			Log.Info("writes %s error:%v\n", me, err)

			return err
		}

		if crlf {
			_, err := f.WriteString(Crlf)
			if nil != err {
				Log.Info("writes %s error:%v\n", me, err)

				return err
			}
		}
	}

	return nil
}

func (me FileName) Load() ([]byte, error) {
	buf, err := ioutil.ReadFile(me.String())
	if nil != err {
		Log.Info("load %s error:%v\n", me, err)
	}

	return buf, err
}

func (me FileName) SaveJson(obj interface{}) error {
	buf, err := json.Marshal(obj)
	if nil != err {
		Log.Info("save %s json error:%v\n", me, err)

		return err
	}

	return me.Save(buf)
}

func (me FileName) LoadJson(obj interface{}) error {
	buf, err := me.Load()
	if nil != err {
		return err
	}

	err = json.Unmarshal(buf, obj)
	if nil != err {
		Log.Info("load %s json error:%v\n", me, err)

		return err
	}

	return nil
}

func (me FileName) ReadPid() int {
	b, err := me.Load()
	if nil != err {
		return 0
	}

	pidstr := string(b)
	pid, err := strconv.Atoi(pidstr)
	if nil != err {
		Log.Info("load pidfile %s error:%v\n", me, err)

		return 0
	}

	return pid
}

func (me FileName) WritePid() {
	pid := os.Getpid()
	pidstr := strconv.Itoa(pid)

	me.Saves([]string{pidstr}, false)
}

func (me FileName) Exist() bool {
	if Empty == me {
		return false
	} else if _, err := os.Stat(me.String()); nil != err {
		return false
	} else {
		return true
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
