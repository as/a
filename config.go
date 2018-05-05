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
			0: Tag0,
		},
		Ctl: events,
	}

	ColConfig = &tag.Config{
		Margin:     image.Pt(15, 0),
		Filesystem: newfsclient(),
		Facer:      font.NewFace,
		FaceHeight: *ftsize,
		Color: [3]frame.Color{
			0: Tag1,
		},
		Ctl: events,
	}

	TagConfig = &tag.Config{
		Margin:     image.Pt(15, 0),
		Filesystem: newfsclient(),
		Facer:      font.NewFace,
		FaceHeight: *ftsize,
		Color: [3]frame.Color{
			0: Tag2,
			1: Body2,
		},
		Image: true,
		Ctl:   events,
	}
)
