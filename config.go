package main

import (
	"flag"
)

const (
	// With compression: burlesque.kch#opts=c#zcomp=gz#msiz=524288000
	DefaultProductionStorage = "burlesque.kch#msiz=524288000"
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
	flag.StringVar(&Config.Storage, "storage", "-", "Kyoto Cabinet storage path (e.g. "+DefaultProductionStorage+")")
	flag.StringVar(&Config.Env, "environment", "development", "Process environment: development or production")
	flag.IntVar(&Config.Port, "port", 4401, "Server HTTP port")
	flag.StringVar(&Config.Rollbar, "rollbar", "", "Rollbar token")
	flag.Parse()

	if Config.Env == "production" && Config.Storage == "-" {
		Config.Storage = DefaultProductionStorage
	}
}
