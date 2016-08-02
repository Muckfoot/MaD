package main

import (
	"log"
	"runtime"
)

func checkErr(err error) {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		log.Printf("file: %s, line: %d, error: %s", file, line, err)
	}
}
