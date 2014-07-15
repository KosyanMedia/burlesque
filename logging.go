package main

import (
	"github.com/stvp/rollbar"
	"log"
	"os"
	"runtime"
)

var (
	logger *log.Logger
)

func SetupLogging() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	rollbar.Token = Config.Rollbar
	rollbar.Environment = Config.Env

	Log("Burlesque started in %s environment", Config.Env)
	Log("GOMAXPROCS is set to %d", runtime.GOMAXPROCS(-1))
	Log("Storage path: %s", Config.Storage)
	Log("Server is running at http://127.0.0.1:%d", Config.Port)
	if Config.Rollbar == "" {
		Log("Rollbar is not configured and disabled")
	} else {
		Log("Rollbar token: %s", Config.Rollbar)
	}
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
