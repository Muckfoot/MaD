package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

func initLogs() (logf *os.File) {
	logf, err := os.OpenFile("errors.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logf.Close()
	log.SetOutput(logf)
	return logf
}

func checkErr(err error) {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		log.Printf("file: %s, line: %d, error: %s", file, line, err)
	}
}

func getCfg() Configuration {
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Configuration
	json.Unmarshal(raw, &c)
	return c
}
