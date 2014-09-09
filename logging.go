package main

import (
	loglib "log"
	"os"
	"runtime"
)

var (
	logger *loglib.Logger
)

func setupLogging() {
	logger = loglib.New(os.Stdout, "", loglib.Ldate|loglib.Lmicroseconds)

	log("Burlesque v%s started", version)
	log("GOMAXPROCS is set to %d", runtime.GOMAXPROCS(-1))
	log("Storage path: %s", config.storage)
	log("Server is running at http://127.0.0.1:%d", config.port)
}

func log(format string, args ...interface{}) {
	logger.Printf("[INFO]  "+format, args...)
}

func alert(err error, format string, args ...interface{}) {
	logger.Printf("[ERROR] "+format, args...)
	logger.Printf("        %s", err.Error())
}
