package globals

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var GlobalShutdown bool = false

func HandleGlobalShutdown(wg *sync.WaitGroup) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	// Wait for Interrupt
	<-osSignal
	GlobalShutdown = true
	wg.Done()
	return
}
