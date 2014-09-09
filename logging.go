package main

import (
	"log"
	"os"
	"runtime"
)

var (
	logger *log.Logger
)

func SetupLogging() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	Log("Burlesque v%s started", Version)
	Log("GOMAXPROCS is set to %d", runtime.GOMAXPROCS(-1))
	Log("Storage path: %s", Config.Storage)
	Log("Server is running at http://127.0.0.1:%d", Config.Port)
}

func Log(format string, args ...interface{}) {
	logger.Printf("[INFO]  "+format, args...)
}

func Error(err error, format string, args ...interface{}) {
	logger.Printf("[ERROR] "+format, args...)
	logger.Printf("        %s", err.Error())
}
