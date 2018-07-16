package main

import (
	"image"

	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/ui/tag"
)

var (
	GridConfig = &tag.Config{
		Margin:     image.Pt(15, 0),
		Filesystem: newfsclient(),
		Facer:      font.NewFace,
		FaceHeight: *ftsize,
		Color: [3]frame.Color{
			Palette["grid"],
		},
		Ctl: events,
	}

	ColConfig = &tag.Config{
		Margin:     image.Pt(15, 0),
		Filesystem: newfsclient(),
		Facer:      font.NewFace,
		FaceHeight: *ftsize,
		Color: [3]frame.Color{
			Palette["col"],
		},
		Ctl: events,
	}

	TagConfig = &tag.Config{
		Margin:     image.Pt(15, 0),
		Filesystem: newfsclient(),
		Facer:      font.NewFace,
		FaceHeight: *ftsize,
		Color: [3]frame.Color{
			Palette["tag"],
			Palette["win"],
		},
		Ctl: events,
	}
)
