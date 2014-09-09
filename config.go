package main

import (
	"flag"
)

var (
	config struct {
		storage string
		port    int
	}
)

func setupConfig() {
	flag.StringVar(&config.storage, "storage", "-", "Kyoto Cabinet storage path (e.g. burlesque.kch#dfunit=8#msiz=512M)")
	flag.IntVar(&config.port, "port", 4401, "Server HTTP port")
	flag.Parse()
}
