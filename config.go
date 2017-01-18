package asdf

import (
	"net/http"
	"strconv"
	"time"
)

const (
	HttpReadTimeout  = 10 * time.Second
	HttpWriteTimeout = 10 * time.Second
	HttpMaxHeader    = 1 << 20
)

type ProtoJsonCfg struct {
	Host string
	Port int
}

type HttpJsonCfg struct {
	ProtoJsonCfg

	Name         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxHeader    int
}

func NewHttpJsonCfg(name, host string, port int) *HttpJsonCfg {
	return &HttpJsonCfg{
		ProtoJsonCfg: ProtoJsonCfg{
			Host: host,
			Port: port,
		},

		Name:         name,
		ReadTimeout:  HttpReadTimeout,
		WriteTimeout: HttpWriteTimeout,
		MaxHeader:    HttpMaxHeader,
	}
}

func (me *HttpJsonCfg) Server(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:           me.Host + ":" + strconv.Itoa(me.Port),
		Handler:        handler,
		ReadTimeout:    me.ReadTimeout,
		WriteTimeout:   me.WriteTimeout,
		MaxHeaderBytes: me.MaxHeader,
	}
}

func (me *HttpJsonCfg) Run(handler http.Handler) {
	server := me.Server(handler)

	Log.Debug("%s listen @%s", me.Name, server.Addr)

	go server.ListenAndServe()
}
