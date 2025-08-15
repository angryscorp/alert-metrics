package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

func NewGracefulShutdownNotifier() chan struct{} {
	ch := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigint
		ch <- struct{}{}
	}()

	return ch
}
