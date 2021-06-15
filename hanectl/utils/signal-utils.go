package utils

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

type CallbackFun func()

func CatchOsSignals(cleanupCallback CallbackFun, reloadCallback CallbackFun) {
	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	signal.Notify(sigs)
	// method invoked upon seeing signal
	go func() {
		for {
			s := <-sigs
			switch s {
			case syscall.SIGUSR1:
				// kill -SIGUSR1 <pid>
				// ignore me !!
			case syscall.SIGUSR2:
				// kill -SIGUSR2 <pid>
				reloadCallback()
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGKILL:
				fallthrough
			case syscall.SIGTERM:
				fallthrough
			case syscall.SIGQUIT:
				fallthrough
			case syscall.SIGHUP:
				log.Debug().Msgf("received signal: %s, terminating", s)
				cleanupCallback()
				os.Exit(0)
			default:
				//log.Debug().Msgf("received signal: %s, continuing", s)
			}
		}
	}()
}
