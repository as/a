package main

import (
	"image"

	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/ui/tag"
)

var GridConfig = &tag.Config{
	Margin:     image.Pt(15, 0),
	Filesystem: newfsclient(),
	Facer:      font.NewFace,
	FaceHeight: *ftsize,
	Color: [3]frame.Color{
		0: frame.ATag0,
	},
	Image: true,
	Ctl:   events,
}

var TagConfig = &tag.Config{
	Margin:     image.Pt(15, 0),
	Filesystem: newfsclient(),
	Facer:      font.NewFace,
	FaceHeight: *ftsize,
	Color: [3]frame.Color{
		0: frame.ATag1,
	},
	Image: true,
	Ctl:   events,
}
