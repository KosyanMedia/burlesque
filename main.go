package main

import (
	"os"
	"os/signal"
	"syscall"
)

const (
	version = "0.1.3"
)

func handleShutdown() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch

		saveState()
		log("State successfully persisted")

		closeStorage()

		log("Stopped")
		os.Exit(0)
	}()
}

func main() {
	setupConfig()
	setupLogging()
	setupStorage()
	setupServer()
	handleShutdown()
	loadState()
	go keepStatePersisted()
	startServer()
}
