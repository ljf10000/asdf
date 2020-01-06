package asdf

import (
	"net"
	"time"
)

type IConn interface {
	IString

	net.Conn
}

func connIoWithTimeout(conn IConn,
	read bool,
	buf []byte,
	timeout time.Duration,
	errTimeoutRedirect error) (int, error) {
	//--------------------------------------------

	if timeout > 0 {
		t := time.Now().Add(timeout)

		if read {
			conn.SetReadDeadline(t)
		} else {
			conn.SetWriteDeadline(t)
		}
	}

	var err error
	var act string
	n := 0
	total := 0
	size := len(buf)

	for ; total < size; total += n {
		if read {
			n, err = conn.Read(buf[total:])
			act = "read"
		} else {
			n, err = conn.Write(buf[total:])
			act = "write"
		}

		if nil != err {
			if IsTimeout(err) {
				Log.Error("%s %s timeout[%d ms]", conn, act, timeout/time.Millisecond)

				if nil != errTimeoutRedirect {
					return total, errTimeoutRedirect
				}
			} else {
				Log.Error("%s %s error: %s", conn, act, err)
			}

			return total, err
		}

		want := size - total
		if n < want {
			Log.Debug("%s %s size[%d] < want[%d]", conn, act, n, want)
		}
	}

	return total, nil
}

func ConnReadWithTimeout(conn IConn, buf []byte, timeout time.Duration, errTimeoutRedirect error) (int, error) {
	return connIoWithTimeout(conn, true, buf, timeout, errTimeoutRedirect)
}

func ConnWriteWithTimeout(conn IConn, buf []byte, timeout time.Duration, errTimeoutRedirect error) (int, error) {
	return connIoWithTimeout(conn, false, buf, timeout, errTimeoutRedirect)
}
