package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

func Await() {
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)
	for sig := range sigCh {
		switch sig {
		case syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT:
			return
		}
	}
}
