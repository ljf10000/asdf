package asdf

import (
	"fmt"
)

const (
	cmdPrefixMaybe = '['
	cmdPrefixMust  = '{'
)

type CmdOptions []string

func (me CmdOptions) Help() string {
	s := Tab

	for _, opt := range me {
		s += opt + Space
	}

	return s
}

func (me CmdOptions) check() (int, error) {
	count := len(me)

	// only last option support maybe
	if cmdPrefixMaybe == me[count-1][0] {
		count--
	}

	for i := 0; i < count; i++ {
		if cmdPrefixMaybe == me[i][0] {
			return 0, ErrBadFormat
		}
	}

	return count, nil
}

func (me CmdOptions) match(args []string) bool {
	for idx, opt := range me {
		switch opt[0] {
		case cmdPrefixMust, cmdPrefixMaybe: // ignore
		default:
			if opt != args[idx] {
				return false
			}
		}
	}

	return true
}

type Cmd struct {
	argc    int // ignore [xxx]
	Options CmdOptions
	Handler func(args []string) error
}

func (me *Cmd) Help() string {
	return me.Options.Help()
}

func (me *Cmd) init() error {
	argc, err := me.Options.check()
	if nil != err {
		return err
	}

	me.argc = argc

	return nil
}

func (me *Cmd) match(args []string) bool {
	argc := len(args)
	count := len(me.Options)

	if argc == count || argc == me.argc {
		return me.Options.match(args)
	} else {
		return false
	}
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

func (me *Command) init() error {
	for _, cmd := range me.Cmds {
		if err := cmd.init(); nil != err {
			return err
		}
	}

	return nil
}

func (me *Command) Run(args []string) error {
	err := me.init()
	if nil != err {
		return err
	}

	cmd := me.match(args)
	if nil != cmd {
		return cmd.Handler(args)
	}

	fmt.Println(me.Help())

	return nil
}
