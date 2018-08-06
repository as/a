package main

import (
	"io"
	"sync"
)

type Funnel struct {
	sync.Mutex
	io.Writer
}

func (f *Funnel) Read(p []byte) (n int, err error) {
	return
}

func (f *Funnel) Write(p []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	return f.Writer.Write(p)
}

func (f *Funnel) Unlock() {
	f.Mutex.Unlock()
	repaint()
}
