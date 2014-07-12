package main

import (
	"flag"
	"fmt"
)

type (
	Config struct {
		Storage string
		Env     string
		Port    int
		Rollbar string
	}
)

var (
	cfg = Config{}
)

func SetupConfig() {
	cfg.Storage = *flag.String("storage", "-", "Kyoto Cabinet storage path (e.g. storage.kch#zcomp=gz#capsiz=524288000)")
	cfg.Env = *flag.String("environment", "development", "Process environment: development or production")
	cfg.Port = *flag.Int("port", 4401, "Server HTTP port")
	cfg.Rollbar = *flag.String("rollbar", "", "Rollbar token")
	flag.Parse()
}

func (c Config) PortString() string {
	return fmt.Sprintf(":%d", cfg.Port)
}
