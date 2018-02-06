package asdf

import (
	"os"
)

const (
	ENV_SERVER = "SERVER"
	ENV_PORT   = "PORT"
	ENV_DIR    = "DIR"
	ENV_DEBUG  = "DEBUG"
)

type Env struct {
	Server string
	Port   string
	Dir    string
	Debug  bool
}

func (me *Env) Load() {
	if v := os.Getenv(ENV_SERVER); Empty != v {
		me.Server = v

		Log.Info("env: SERVER=%s", v)
	}

	if v := os.Getenv(ENV_PORT); Empty != v {
		me.Port = v

		Log.Info("env: PORT=%s", v)
	}

	if v := os.Getenv(ENV_DIR); Empty != v {
		me.Dir = v

		Log.Info("env: DIR=%s", v)
	}

	if v := os.Getenv(ENV_DEBUG); Empty != v {
		me.Debug = true

		Log.Info("env: DEBUG=%s", v)
	}
}
