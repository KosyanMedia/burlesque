package main

import (
	"github.com/stvp/rollbar"
	"log"
	"os"
)

var (
	logger *log.Logger
)

func SetupLogging() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	rollbar.Token = cfg.Rollbar
	rollbar.Environment = cfg.Env
}

func Log(format string, args ...interface{}) {
	logger.Printf("[INFO] "+format, args...)
}

func Debug(format string, args ...interface{}) {
	if cfg.Env == "development" {
		logger.Printf("[DEBUG]"+format, args...)
	}
}

func Error(err error, format string, args ...interface{}) {
	logger.Printf("[ERROR]"+format, args...)

	if cfg.Env == "development" {
		panic(err)
	}
	if cfg.Rollbar != "" {
		rollbar.Error("error", err)
	}
}
