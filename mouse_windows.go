package main

import (
	"image"
	"os"

	"github.com/as/cursor"
	"github.com/as/ms/win"
)

var (
	winfd    win.Window
	winfderr error
)

func tryWindow() {
	if winfd != 0 && winfderr == nil {
		return
	}
	winfd, winfderr = win.Open(os.Getpid())
}

func moveMouse(pt image.Point) {
	tryWindow()
	abs, err := winfd.Client()
	if err != nil {
		winfderr = err
		return
	}
	pt = pt.Add(abs.Min)
	cursor.MoveTo(pt)
}
