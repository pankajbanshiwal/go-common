package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

func Handle(handler func()) <-chan struct{} {
	ch := make(chan struct{}, 1)

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
		<-shutdown
		if handler != nil {
			handler()
		}
		ch <- struct{}{}
	}()

	return ch
}
