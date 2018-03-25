package main

import "log"

var logFunc = log.Printf

func init() {
	log.SetFlags(log.Llongfile)
	log.SetPrefix("a: ")
}

func logf(fm string, v ...interface{}) {
	logFunc(fm, v...)
}

func setLogFunc(f func(string, ...interface{})) {
	logFunc = f
}
