package main

import (
	"github.com/as/font"
	"github.com/as/ui/tag"
)

var fontmap = make(map[font.Facer]int)

var fontfuncs = [...]func(int) font.Face{
	font.NewFace,
	font.NewGoMedium,
	font.NewGoRegular,
	font.NewGoMono,
}

func nextFace(f font.Facer) {
	fontmap[f]++
	fn := fontfuncs[fontmap[f]%len(fontfuncs)]
	if f, ok := f.(*tag.Tag); ok {
		f.Config.Facer = fn
	}
	f.SetFont(fn(f.Face().Height()))
}
