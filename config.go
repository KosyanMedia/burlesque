package main

import (
	"flag"
)

var (
	Config struct {
		Storage string
		Env     string
		Port    int
	}
)

func SetupConfig() {
	flag.StringVar(&Config.Storage, "storage", "-", "Kyoto Cabinet storage path (e.g. burlesque.kch#dfunit=8#msiz=512M)")
	flag.StringVar(&Config.Env, "environment", "production", "Process environment: production or development")
	flag.IntVar(&Config.Port, "port", 4401, "Server HTTP port")
	flag.Parse()
}
