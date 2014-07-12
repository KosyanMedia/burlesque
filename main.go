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
		rollbar.Wait()
		Log("Storage closed")
		Log("Server stopped")
		os.Exit(1)
	}()
}

func main() {
	SetupLogging()
	SetupConfig()
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
