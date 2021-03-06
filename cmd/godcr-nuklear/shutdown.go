package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// triggered after program execution is complete or if interrupt signal is received
var beginShutdown = make(chan bool)

// shutdownOps holds cleanup/shutdown functions that should be executed when shutdown signal is triggered
var shutdownOps []func()

var exitCode = 0

func listenForShutdownRequests() {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// listen for the initial interrupt request and trigger shutdown signal
	sig := <-interruptChannel
	log.Infof("Received %s signal. Shutting down...", sig)
	beginShutdown <- true

	// continue to listen for interrupt requests and log that shutdown has already been signaled
	for {
		<-interruptChannel
		log.Warnf(" Already shutting down... Please wait")
	}
}

func handleShutdownRequests(wg *sync.WaitGroup) {
	// make wait group wait till shutdownSignal is received and shutdownOps performed
	wg.Add(1)

	<-beginShutdown
	for _, shutdownOp := range shutdownOps {
		shutdownOp()
	}

	// shutdown complete
	wg.Done()

	os.Exit(exitCode)
}
