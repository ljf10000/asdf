package asdf

import (
	"fmt"
)

type CmdOptions []string

func (me CmdOptions) Help() string {
	s := Tab

	for _, opt := range me {
		s += opt + Space
	}

	return s[:len(s)-1]
}

func (me CmdOptions) match(args []string) bool {
	count := len(args)
	if len(me) != count {
		return false
	}

	for idx, opt := range me {
		switch opt[0] {
		case '{', '[':
		default:
			if opt != args[idx] {
				return false
			}
		}
	}

	return true
}

type Cmd struct {
	Options CmdOptions
	Handler func(args []string) error
}

func (me *Cmd) Help() string {
	return me.Options.Help()
}

func (me *Cmd) match(args []string) bool {
	return me.Options.match(args)
}

type Command struct {
	Self string
	Cmds []*Cmd
}

func (me *Command) Help() string {
	s := me.Self + ":" + Crlf

	for _, cmd := range me.Cmds {
		s += cmd.Help() + Crlf
	}

	return s
}

func (me *Command) match(args []string) *Cmd {
	for _, cmd := range me.Cmds {
		if cmd.match(args) {
			return cmd
		}
	}

	return nil
}

func (me *Command) Run(args []string) error {
	cmd := me.match(args)
	if nil != cmd {
		return cmd.Handler(args)
	}

	fmt.Println(me.Help())

	return nil
}
