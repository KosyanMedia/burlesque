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
	rollbar.Token = Config.Rollbar
	rollbar.Environment = Config.Env
}

func Log(format string, args ...interface{}) {
	logger.Printf("[INFO] "+format, args...)
}

func Debug(format string, args ...interface{}) {
	if Config.Env == "development" {
		logger.Printf("[DEBUG] "+format, args...)
	}
}

func Error(err error, format string, args ...interface{}) {
	logger.Printf("[ERROR] "+format, args...)

	if Config.Env == "development" {
		panic(err)
	}
	if Config.Rollbar != "" {
		rollbar.Error("error", err)
	}
}
