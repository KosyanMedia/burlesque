package main

import (
	"github.com/stvp/rollbar"
	logpkg "log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type (
	Message []byte
	Key     []byte
)

func NewKey(queue string, index uint) Key {
	istr := strconv.FormatUint(uint64(index), 10)
	key := strings.Join([]string{queue, istr}, "_")
	return Key(key)
}

var (
	log *logpkg.Logger
)

func HandleShutdown() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch
		SaveState()
		log.Printf("State successfully persisted")
		storage.Close()
		rollbar.Wait()
		log.Println("Storage closed")
		log.Printf("Server stopped")
		os.Exit(1)
	}()
}

func main() {
	log = logpkg.New(os.Stdout, "", logpkg.Ldate|logpkg.Lmicroseconds)

	rollbar.Token = "c91028beb8434b66882f59f55f22659d" // klit access token
	rollbar.Environment = cfg.Env

	SetupConfig()
	SetupStorage()
	SetupServer()
	HandleShutdown()
	LoadState()
	go KeepStatePersisted()
	go PersistMessages()

	log.Printf("GOMAXPROCS = %d", runtime.GOMAXPROCS(-1))
	log.Printf("Starting HTTP server on port %d", cfg.Port)

	http.ListenAndServe(cfg.PortString(), nil)
}
