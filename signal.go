package asdf

import (
	"os"
	"os/signal"
	"syscall"
)

var sigDeft = []os.Signal{
	syscall.SIGTERM,
	syscall.SIGINT,
}

type Signal = os.Signal

var chSignal = make(chan Signal, 8)

func SignalService(sigs []Signal, handle func(sig Signal)) {
	if 0 == len(sigs) {
		sigs = sigDeft
	}

	signal.Notify(chSignal, sigs...)
	Log.Crit("setup signal[%v]", sigs)

	for {
		select {
		case sig := <-chSignal:
			Log.Crit("recv signal: %s", sig)

			handle(sig)
		}
	}
}
