package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/KosyanMedia/burlesque/hub"
	"github.com/KosyanMedia/burlesque/storage"
)

const (
	version = "0.2.0"
)

var (
	theStorage *storage.Storage
	theHub     *hub.Hub
	config     struct {
		storage string
		port    int
	}
)

func main() {
	flag.StringVar(&config.storage, "storage", "-", "Kyoto Cabinet storage path (e.g. burlesque.kch#dfunit=8#msiz=512M)")
	flag.IntVar(&config.port, "port", 4401, "Server HTTP port")
	flag.Parse()

	theStorage, err := storage.New(config.storage)
	if err != nil {
		panic(err)
	}

	theHub = hub.New(theStorage)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch
		theStorage.Close()
		os.Exit(0)
	}()

	setupLogging()
	startServer()
}
