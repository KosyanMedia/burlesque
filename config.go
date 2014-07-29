package main

import (
	"flag"
)

const (
	// With compression: burlesque.kch#opts=c#zcomp=gz#msiz=524288000
	DefaultProductionStorage = "burlesque.kch#dfunit=8#msiz=512M"
)

var (
	Config struct {
		Storage string
		Env     string
		Port    int
	}
)

func SetupConfig() {
	flag.StringVar(&Config.Storage, "storage", "-", "Kyoto Cabinet storage path (e.g. "+DefaultProductionStorage+")")
	flag.StringVar(&Config.Env, "environment", "production", "Process environment: production or development")
	flag.IntVar(&Config.Port, "port", 4401, "Server HTTP port")
	flag.Parse()

	if Config.Env == "production" && Config.Storage == "-" {
		Config.Storage = DefaultProductionStorage
	}
}
