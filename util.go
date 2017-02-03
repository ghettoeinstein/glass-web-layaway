package main

import (
	"log"
	"os"
)

func setupLogging() {
	f, err := os.OpenFile("logs/glassLogs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error  opening up logfile %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
}
