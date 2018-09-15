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
	TagConfig.FaceHeight = *ftsize
	if *utf8 {
		GridConfig.Tag.Frame.Flag |= frame.FrUTF8
		ColConfig.Tag.Frame.Flag |= frame.FrUTF8
		TagConfig.Tag.Frame.Flag |= frame.FrUTF8
		TagConfig.Body.Frame.Flag |= frame.FrUTF8
		TagConfig.Facer = func(n int) font.Face {
			return font.NewRune(font.NewGoMedium(n))
		}
	}
	if *elastic {
		TagConfig.Body.Frame.Flag |= frame.FrElastic
	}
	dev, err := ui.Init(&ui.Config{
		Width: winSize.X, Height: winSize.Y,
		Title:   "A",
		Overlay: true,
	})
	ckfault(err)
	if dev == nil {
		panic("no device")
	}
	return dev, dev.Window(), screen.Dev, font.NewFace(*ftsize)
}
