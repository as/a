package main

import (
	"image"

	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/shiny/screen"
	"github.com/as/ui"
)

var (
	winSize = image.Pt(1024, 768)
)

func frameinstall() (ui.Dev, screen.Window, *screen.Device, font.Face) {
	if *oled {
		usedarkcolors()
	}
	frame.ForceUTF8 = *utf8
	frame.ForceElastic = *elastic
	dev, err := ui.Init(&screen.NewWindowOptions{
		Width: winSize.X, Height: winSize.Y,
		Title: "A",
	})
	ckfault(err)
	if dev == nil {
		panic("no device")
	}
	return dev, dev.Window(), screen.Dev, font.NewFace(*ftsize)
}
