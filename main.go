package main

import (
	"os"
	"os/signal"
	"syscall"
)

const (
	Version = "0.1.3"
)

func HandleShutdown() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch

		SaveState()
		Log("State successfully persisted")

		CloseStorage()

		Log("Stopped")
		os.Exit(0)
	}()
}

func main() {
	SetupConfig()
	SetupLogging()
	SetupStorage()
	SetupServer()
	HandleShutdown()
	LoadState()
	go KeepStatePersisted()
	StartServer()
}
