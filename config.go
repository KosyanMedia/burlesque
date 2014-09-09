package main

import (
	"flag"
)

var (
	Config struct {
		Storage string
		Port    int
	}
)

func SetupConfig() {
	flag.StringVar(&Config.Storage, "storage", "-", "Kyoto Cabinet storage path (e.g. burlesque.kch#dfunit=8#msiz=512M)")
	flag.IntVar(&Config.Port, "port", 4401, "Server HTTP port")
	flag.Parse()
}
