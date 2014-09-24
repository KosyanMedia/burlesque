package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/KosyanMedia/burlesque/hub"
	"github.com/KosyanMedia/burlesque/server"
	"github.com/KosyanMedia/burlesque/storage"
)

const (
	version = "0.2.0"
)

func main() {
	var (
		storagePath string
		port        int
	)

	flag.StringVar(&storagePath, "storage", "-", "Kyoto Cabinet storage path (e.g. burlesque.kch#dfunit=8#msiz=512M)")
	flag.IntVar(&port, "port", 4401, "Server HTTP port")
	flag.Parse()

	store, err := storage.New(storagePath)
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-shutdown
		store.Close()
		os.Exit(0)
	}()

	fmt.Printf("Burlesque v%s started\n", version)
	fmt.Printf("GOMAXPROCS is set to %d\n", runtime.GOMAXPROCS(-1))
	fmt.Printf("Storage path: %s\n", storagePath)
	fmt.Printf("Server is running at http://127.0.0.1:%d\n", port)

	h := hub.New(store)
	s := server.New(port, h)

	s.Start()
}
