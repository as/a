package main

import (
	"os"
	"runtime"
	"strconv"

	"github.com/as/font"
	"github.com/as/ui/tag"
)

type facer interface {
	Face() font.Face
	SetFont(font.Face)
}

var fontmap = make(map[facer]int)

var fontfuncs = [...]func(int) font.Face{
	font.NewFace,
	font.NewGoMedium,
	font.NewGoRegular,
	font.NewGoMono,
}

func nextFace(f facer) {
	fontmap[f]++
	fn := fontfuncs[fontmap[f]%len(fontfuncs)]
	if f, ok := f.(*tag.Tag); ok {
		f.Config.Facer = fn
	}
	f.SetFont(fn(f.Face().Height()))
}

func defaultFaceSize() int {
	if s := os.Getenv("fontsize"); s != "" {
		// user specified a font size, so let's
		// just go with it
		v, err := strconv.Atoi(s)
		if err == nil {
			return v
		}
	}
	switch runtime.GOOS {
	case "darwin":
		// darwin begets a larger font; might not be enough;
		// TODO(as): proper DPI-aware scaling
		return 13
	default:
		// de-facto standard for best looking font size
		// based on references I don't have anymore
		return 11
	}
}
