package asdf

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

type TcpStream net.TCPConn

func (me *TcpStream) TcpConn() *net.TCPConn {
	return (*net.TCPConn)(me)
}

func TcpStreamDial(addr *net.TCPAddr) (*TcpStream, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if nil != err {
		return nil, err
	} else {
		return NewTcpStream(conn), nil
	}
}

func NewTcpStream(conn *net.TCPConn) *TcpStream {
	return (*TcpStream)(conn)
}

func (me *TcpStream) ReadN(size int) ([]byte, error) {
	var err error

	if size <= 0 {
		return nil, ErrBadLen
	}

	buf := make([]byte, size)
	Len := 0
	all := 0

	conn := me.TcpConn()

	for {
		Len, err = conn.Read(buf[all:])
		if nil != err {
			Log.Info("conn read error:%v", err)

			return nil, err
		}

		all += Len
		if all == size {
			return buf, nil
		} else if all > size {
			Log.Info("conn read too more")

			return nil, ErrTooMore
		}

		// continue
	}
}

func (me *TcpStream) WriteN(buf []byte) error {
	if nil == buf {
		return ErrBadBuf
	}

	size := len(buf)
	Len := 0
	all := 0

	conn := me.TcpConn()

	var err error

	for {
		Len, err = conn.Write(buf[all:])
		if nil != err {
			Log.Info("conn write error:%v", err)

			return err
		}

		all += Len
		if all == size {
			return nil
		} else if all > size {
			Log.Info("conn write too more")

			return ErrTooMore
		}

		// continue
	}
}

// todo: check io.EOF
func (me *TcpStream) Read() ([]byte, error) {
	var buf []byte
	var err error

	buf, err = me.ReadN(SizeofInt32)
	if nil != err {
		Log.Info("conn read head error:%v", err)

		return nil, err
	}

	size := int(binary.BigEndian.Uint32(buf))
	if size <= 0 {
		Log.Info("conn read invalid head size")

		return nil, ErrBadLen
	}

	buf, err = me.ReadN(size)
	if nil != err {
		Log.Info("conn read body error:%v", err)

		return nil, err
	}

	return buf, nil
}

func (me *TcpStream) Write(buf []byte) error {
	if nil == buf {
		return ErrBadBuf
	}

	var err error
	var head [SizeofInt32]byte
	binary.BigEndian.PutUint32(head[:], uint32(len(buf)))

	err = me.WriteN(head[:])
	if nil != err {
		Log.Info("conn write head error:%v", err)

		return err
	}

	err = me.WriteN(buf)
	if nil != err {
		Log.Info("conn write body error:%v", err)

		return err
	}

	return nil
}

func (me *TcpStream) WriteMulti(bufs [][]byte) error {
	if nil == bufs {
		return ErrBadBuf
	}

	var err error
	var head [SizeofInt32]byte

	Len := 0
	for _, buf := range bufs {
		Len += len(buf)
	}

	binary.BigEndian.PutUint32(head[:], uint32(Len))

	err = me.WriteN(head[:])
	if nil != err {
		Log.Info("conn write head error:%v", err)
		return err
	}

	for _, buf := range bufs {
		err = me.WriteN(buf)
		if nil != err {
			Log.Info("conn write body error:%v", err)

			return err
		}
	}

	return nil
}

func (me *TcpStream) ReadJsonObj(obj interface{}) ([]byte, error) {
	buf, err := me.Read()
	if nil != err {
		return nil, err
	}

	err = json.Unmarshal(buf, obj)
	if nil != err {
		return nil, err
	}

	return buf, nil
}

func (me *TcpStream) WriteJsonObj(obj interface{}) error {
	buf, err := json.Marshal(obj)
	if nil != err {
		return err
	}

	return me.Write(buf)
}

type TcpStreamError struct {
	Err  int    `json:"err"`
	Info string `json:"info"`
}

func NewTcpStreamError(err int, info string) *TcpStreamError {
	return &TcpStreamError{
		Err:  err,
		Info: info,
	}
}

func (me *TcpStream) Error(errno int, err error) error {
	msg := NewTcpStreamError(errno, err.Error())

	me.WriteJsonObj(msg)

	return err
}

func (me *TcpStream) OK() error {
	msg := NewTcpStreamError(0, Empty)

	me.WriteJsonObj(msg)

	return nil
}

type TcpStreamJsonCmd struct {
	Cmd int `json:"cmd"`
}

type TcpStreamJsonMsg struct {
	Chan chan error
	Cmd  int
	Buf  []byte
}

func tcpStreamJsonCall(jch chan *TcpStreamJsonMsg, cmd int, buf []byte) error {
	ch := make(chan error)

	jch <- &TcpStreamJsonMsg{
		Chan: ch,
		Cmd:  cmd,
		Buf:  buf,
	}

	err := <-ch

	close(ch)

	return err
}

func TcpStreamJsonHandle(jch chan *TcpStreamJsonMsg, stream *TcpStream) error {
	req := &TcpStreamJsonCmd{}

	buf, err := stream.ReadJsonObj(req)
	if nil != err {
		return stream.Error(-1, err)
	}

	err = tcpStreamJsonCall(jch, req.Cmd, buf)
	if nil != err {
		return stream.Error(-2, err)
	}

	stream.OK()

	return stream.Close()
}

type TcpStreamJsonObjHandle func(obj interface{}, handle func() error) error
type TcpStreamJsonCmdHandle func(cmd int, handle TcpStreamJsonObjHandle) error

func TcpStreamJsonMain(jch chan *TcpStreamJsonMsg, cmdHandle TcpStreamJsonCmdHandle) {
	for {
		msg := <-jch

		objHandle := func(obj interface{}, handle func() error) error {
			err := json.Unmarshal(msg.Buf, obj)
			if nil != err {
				return err
			}

			return handle()
		}

		msg.Chan <- cmdHandle(msg.Cmd, objHandle)
	}
}
