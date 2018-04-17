package asdf

import (
	"encoding/json"
	"errors"
	"fmt"
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
	return FileName(CurrentDirFile(string(me)))
}

func FileShortName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if os.IsPathSeparator(path[i]) {
			return path[i:]
		}
	}

	return Empty
}

func (me FileName) ShortName() string {
	return FileShortName(me.String())
}

func (me FileName) MkDir() error {
	dir := filepath.Dir(string(me))

	return os.MkdirAll(dir, 0755)
}

func (me FileName) Append(buf []byte) error {
	f, err := os.OpenFile(string(me), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		Log.Error("open %s error: %s", me, err)

		return err
	}

	_, err = f.Write(buf)
	if nil != err {
		Log.Error("write %s error: %s", me, err)

		return err
	}

	f.Close()

	return nil
}

func (me FileName) AppendLine(line string) error {
	return me.Append([]byte(line + Crlf))
}

func (me FileName) Save(buf []byte) error {
	err := ioutil.WriteFile(string(me), buf, FilePermNormal)
	if nil != err {
		Log.Error("save %s error: %s", me, err)
	}

	return err
}

func (me FileName) Delete() error {
	err := os.Remove(string(me))
	if nil != err {
		Log.Error("delete %s error: %s", me, err)
	}

	return err
}

func (me FileName) Touch(Time Time32) error {
	if 0 == Time {
		Time = NowTime32()
	}

	tm := time.Unix(int64(Time), 0)

	err := os.Chtimes(string(me), tm, tm)
	if nil != err {
		Log.Error("change %s time error: %s", me, err)
	}

	return err
}

func (me FileName) Saves(texts []string, crlf bool) error {
	f, err := os.Create(string(me))
	if nil != err {
		Log.Error("create %s error: %s", me, err)

		return err
	}
	defer f.Close()

	for _, text := range texts {
		_, err := f.WriteString(text)
		if nil != err {
			Log.Error("writes %s error: %s", me, err)

			return err
		}

		if crlf {
			_, err := f.WriteString(Crlf)
			if nil != err {
				Log.Error("writes %s error: %s", me, err)

				return err
			}
		}
	}

	return nil
}

func (me FileName) Load() ([]byte, error) {
	buf, err := ioutil.ReadFile(string(me))
	if nil != err {
		Log.Error("load %s error: %s", me, err)
	}

	return buf, err
}

func (me FileName) LoadByLine(lineHandle func(line string) error) error {
	line := Empty

	if err := lineHandle(line); nil != err {
		return err
	}

	return nil
}

func (me FileName) SaveJson(obj interface{}) error {
	buf, err := json.MarshalIndent(obj, Empty, "\t")
	if nil != err {
		Log.Error("save %s json error: %s", me, err)

		return err
	}

	return me.Save(buf)
}

func (me FileName) LoadJson(obj interface{}) error {
	buf, err := me.Load()
	if nil != err {
		Log.Error("load %s json error: %s", me, err)

		return err
	}

	err = json.Unmarshal(buf, obj)
	if nil != err {
		Log.Error("unmarshal %s json error: %s", me, err)

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
		Log.Error("load pidfile %s error:%v", me, err)

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
	} else if _, err := os.Stat(string(me)); nil != err {
		return false
	} else {
		return true
	}
}

func (me FileName) FileSize() (int64, error) {
	if Empty == me {
		return 0, errors.New("empty filename")
	} else if info, err := os.Stat(string(me)); nil != err {
		return 0, err
	} else {
		return info.Size(), nil
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

func (me FileName) DirScanSimple(scan func(path string, info os.FileInfo) error) error {
	if !me.DirExist() {
		return errors.New(fmt.Sprintf("dir:%s not exist", me))
	}

	return filepath.Walk(me.String(), func(path string, info os.FileInfo, err error) error {
		if nil != err {
			return err
		} else if info.IsDir() {
			return filepath.SkipDir
		}

		return scan(path, info)
	})
}

func (me FileName) DirUnserialize(creator func() IUnserialize) error {
	return filepath.Walk(string(me), func(path string, f os.FileInfo, err error) error {
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
