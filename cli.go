package asdf

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

const SizeofCliServerRecvBuffer = 64 * SizeofK
const CLI_NETWORK = "unix"

/******************************************************************************/

type CliCmd struct {
	Options CmdOptions
	Handler func(args []string) (string, error)
}

func (me *CliCmd) help() string {
	return me.Options.Help()
}

func (me *CliCmd) match(args []string) bool {
	return me.Options.match(args)
}

/******************************************************************************/

type CliCommand []*CliCmd

func (me CliCommand) help() string {
	s := Empty

	for _, cmd := range me {
		s += cmd.help() + Crlf
	}

	return s
}

func (me CliCommand) match(args []string) *CliCmd {
	for _, cmd := range me {
		if cmd.match(args) {
			return cmd
		}
	}

	return nil
}

func (me CliCommand) exec(args []string) (string, error) {
	cmd := me.match(args)
	if nil != cmd {
		return cmd.Handler(args)
	}

	return me.help(), ErrNoSupport
}

/******************************************************************************/

type CliServer struct {
	Path    string
	Command CliCommand
}

func (me *CliServer) Run() {
	os.Remove(me.Path)

	recvbuf := make([]byte, SizeofCliServerRecvBuffer)

	addr := &net.UnixAddr{
		Name: me.Path,
		Net:  CLI_NETWORK,
	}

	listener, err := net.ListenUnix(CLI_NETWORK, addr)
	if err != nil {
		Panic("cli server listen error:%s", err)
	}

	Log.Info("OK: cli server ...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			Log.Error("cli server accept error:", err)

			continue
		}

		n, err := conn.Read(recvbuf)
		if err != nil {
			Log.Error("cli server read error:", err)

			continue
		}

		go me.handle(conn, string(recvbuf[:n]))
	}
}

func (me *CliServer) handle(conn net.Conn, line string) {
	Log.Debug("cli server recv: %s", line)

	args := strings.Split(line, " ")

	var errs string
	s, err := me.Command.exec(args)
	if nil != err {
		errs = err.Error()
	}

	r := &CliResponse{
		Err:  errs,
		Info: s,
	}

	bin, _ := json.Marshal(r)

	conn.Write(bin)
	conn.Close()
}

/******************************************************************************/

type CliResponse struct {
	Err  string
	Info string
}

type CliClient struct {
	Path string
}

func (me *CliClient) Exec(args []string) (string, error) {
	addr := &net.UnixAddr{
		Name: me.Path,
		Net:  CLI_NETWORK,
	}

	conn, err := net.DialUnix(CLI_NETWORK, nil, addr)
	if err != nil {
		return Empty, ErrSprintf("dial cli server error:%s", err)
	}

	line := strings.Join(args, " ")
	_, err = conn.Write([]byte(line))
	if nil != err {
		return Empty, ErrSprintf("write cli server error:%s", err)
	}

	recvbuf := make([]byte, SizeofCliServerRecvBuffer)
	n, err := conn.Read(recvbuf)
	if nil != err {
		return Empty, ErrSprintf("read cli server error:%s", err)
	}

	r := &CliResponse{}
	if err := json.Unmarshal(recvbuf[:n], r); nil != err {
		return Empty, err
	}

	if Empty == r.Err {
		return r.Info, nil
	} else {
		return Empty, errors.New(r.Err)
	}
}

func (me *CliClient) Execf(args []string) error {
	s, err := me.Exec(args)
	if nil != err {
		return err
	}

	fmt.Printf("%s", s)

	return nil
}
