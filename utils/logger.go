package utils

import (
	"log"
	"os"
)

var (
	Log *log.Logger
)

func Initlog(config *Config) func() {
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	if config.Debug {
		return func() {}
	}

	f, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Log.Printf("Error opening log file - %v", err)
		return func() {}
	}
	println("LogFile: " + config.LogFile)
	Log.SetOutput(f)
	Log.Println("-------------")
	Log.Println("Start logging")
	return func() {
		e := f.Close()
		if e != nil {
			Log.Printf("Problem closing the log file: %s\n", e)
		}
	}
}
