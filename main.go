package main

import (
	"github.com/stvp/rollbar"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func HandleShutdown() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch

		SaveState()
		Log("State successfully persisted")

		storage.Close()
		Log("Storage closed")

		Log("Waiting for rollbar...")
		rollbar.Wait()

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
	go PersistMessages()

	Log("GOMAXPROCS = %d", runtime.GOMAXPROCS(-1))
	Log("Starting HTTP server on port %d", cfg.Port)

	http.ListenAndServe(cfg.PortString(), nil)
}
