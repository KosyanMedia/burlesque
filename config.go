package main

import (
	"flag"
)

const (
	DefaultProductionStorage = "burlesque.kch#opts=c#zcomp=gz#msiz=524288000"
)

var (
	Config struct {
		Storage string
		Env     string
		Port    int
		Rollbar string
	}
)

func SetupConfig() {
	Config.Storage = *flag.String("storage", "-", "Kyoto Cabinet storage path (e.g. "+DefaultProductionStorage+")")
	Config.Env = *flag.String("environment", "development", "Process environment: development or production")
	Config.Port = *flag.Int("port", 4401, "Server HTTP port")
	Config.Rollbar = *flag.String("rollbar", "", "Rollbar token")
	flag.Parse()

	if Config.Env == "production" && Config.Storage == "-" {
		Config.Storage = DefaultProductionStorage
	}
}
