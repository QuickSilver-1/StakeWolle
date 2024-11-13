package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()

	file, err := os.OpenFile("log/worker.log", os.O_APPEND, 0666)

	log.SetFlags(log.Ldate | log.Ltime)
	if err != nil {
		file, _ := os.OpenFile("log/file.log", os.O_APPEND, 0666)
		log.SetOutput(file)
		log.Fatal("Failed to open log file: ", err)
	}

	log.SetOutput(file)



}