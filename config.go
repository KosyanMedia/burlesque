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
	}
)

var (
	cfg = Config{}
)

func SetupConfig() {
	cfg.Storage = *flag.String("storage", "-", "Kyoto Cabinet storage path (e.g. storage.kch#zcomp=gz#capsiz=524288000)")
	cfg.Env = *flag.String("environment", "development", "Process environment: development or production")
	cfg.Port = *flag.Int("port", 4401, "HTTP port to listen")
	flag.Parse()
}

func (c Config) PortString() string {
	return fmt.Sprintf(":%d", cfg.Port)
}
