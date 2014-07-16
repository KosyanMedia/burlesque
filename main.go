package main

import (
	"fmt"
	"github.com/stvp/rollbar"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func HandleShutdown() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch

		SaveState()
		Log("State successfully persisted")

		CloseStorage()

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

	port := fmt.Sprintf(":%d", Config.Port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		Error(err, "Error starting server on port %d", Config.Port)
	}
}
