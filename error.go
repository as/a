package main

import (
	"log"
	"os"
)

var logFunc = log.Printf

func init() {
	log.SetFlags(log.Llongfile)
	log.SetPrefix("a: ")
}

func logf(fm string, v ...interface{}) {
	logFunc(fm, v...)
}

func ckfault(err error) {
	if err == nil {
		return
	}
	logf("fault: %s", err)
	os.Exit(1) // TODO(as): if we're in graphical mode, or have files open, we cant do this
}

func setLogFunc(f func(string, ...interface{})) {
	logFunc = f
}
