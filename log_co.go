package asdf

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CoLogChan chan LogMsg

func OpenCoLogger(dir, prefix string, size int, show bool) error {
	ch := make(CoLogChan, size)

	abs, err := filepath.Abs(dir)
	if nil != err {
		return err
	} else if err := os.MkdirAll(abs, FilePermDir); nil != err {
		return err
	}

	logger := &CoLogger{
		level:  LogLevelDeft,
		dir:    abs,
		prefix: prefix,
		ch:     ch,
	}

	Log = logger

	go logger.run(ch, show)

	return nil
}

type CoLogger struct {
	level  LogLevel
	dir    string
	prefix string
	fd     *os.File
	ch     CoLogChan
	year   int
	month  time.Month
	day    int
}

func (me *CoLogger) open() error {
	var err error

	file := fmt.Sprintf("%s.%04d-%02d-%02d.log",
		me.prefix,
		me.year,
		me.month,
		me.day)

	me.fd, err = os.OpenFile(
		filepath.Join(me.dir, file),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if nil != err {
		return err
	}

	return nil
}

func (me *CoLogger) tryCut() error {
	year, month, day := time.Now().Date()

	if year != me.year || month != me.month || day != me.day {
		me.Close()

		me.year, me.month, me.day = year, month, day

		return me.open()
	}

	return nil
}

func (me *CoLogger) run(ch CoLogChan, show bool) {
	fmt.Printf("cologger running...\n")

	me.tryCut()

	for {
		msg, ok := <-ch
		if !ok {
			me.Close()

			panic("co-logger chan error")
		}

		s := msg.String()
		me.fd.WriteString(s)

		if show {
			fmt.Print(s)
		}

		err := me.tryCut()
		if nil != err {
			me.Close()

			panic("co-logger reopen error")
		}

		if msg.panic {
			me.fd.WriteString("Panic: " + s)
			me.Close()

			panic("co-logger recv panic")
		}
	}
}

func (me *CoLogger) GetLevel() LogLevel {
	return me.level
}

func (me *CoLogger) SetLevel(level LogLevel) {
	me.level = level
}

var CoLoggerPanicWait time.Duration = 15

func (me *CoLogger) Panic(format string, v ...interface{}) {
	me.ch <- LogMsg{
		s:     fmt.Sprintf(format+Crlf, v...),
		panic: true,
	}

	go func() {
		time.Sleep(CoLoggerPanicWait * time.Second)

		panic(fmt.Sprintf("co-logger[last] Panic: "+format+Crlf, v...))
	}()
}

func (me *CoLogger) Log(level LogLevel, format string, v ...interface{}) {
	me.ch <- LogMsg{
		s:     fmt.Sprintf(format+Crlf, v...),
		level: level,
	}
}

func (me *CoLogger) Emerg(format string, v ...interface{}) {
	if LogLevelEmerg <= me.level {
		me.Log(LogLevelEmerg, format, v...)
	}
}

func (me *CoLogger) Alert(format string, v ...interface{}) {
	if LogLevelAlert <= me.level {
		me.Log(LogLevelAlert, format, v...)
	}
}

func (me *CoLogger) Crit(format string, v ...interface{}) {
	if LogLevelCrit <= me.level {
		me.Log(LogLevelCrit, format, v...)
	}
}

func (me *CoLogger) Error(format string, v ...interface{}) {
	if LogLevelError <= me.level {
		me.Log(LogLevelError, format, v...)
	}
}

func (me *CoLogger) Warning(format string, v ...interface{}) {
	if LogLevelWarning <= me.level {
		me.Log(LogLevelWarning, format, v...)
	}
}

func (me *CoLogger) Warn(format string, v ...interface{}) {
	if LogLevelWarn <= me.level {
		me.Log(LogLevelWarn, format, v...)
	}
}

func (me *CoLogger) Notice(format string, v ...interface{}) {
	if LogLevelNotice <= me.level {
		me.Log(LogLevelNotice, format, v...)
	}
}

func (me *CoLogger) Info(format string, v ...interface{}) {
	if LogLevelInfo <= me.level {
		me.Log(LogLevelInfo, format, v...)
	}
}

func (me *CoLogger) Debug(format string, v ...interface{}) {
	if LogLevelDebug <= me.level {
		me.Log(LogLevelDebug, format, v...)
	}
}

func (me *CoLogger) Close() error {
	// me.ch = nil

	me.fd.Sync()

	err := me.fd.Close()

	me.fd = nil

	return err
}
