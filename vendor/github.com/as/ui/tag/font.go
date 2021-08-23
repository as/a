package tag

import (
	"github.com/as/font"
)

type facer interface {
	Face() font.Face
	SetFont(font.Face)
}

// Height returns the recommended minimum pixel height for a tag label
// given the face height in pixels.
func Height(facePix int) int {
	if facePix == 0 {
		facePix = DefaultConfig.FaceHeight
	}
	return facePix + facePix/2 + facePix/3
}

// Face sets the font face
func (w *Tag) Face() font.Face{
	if f, ok := w.Window.(facer); ok{
		return f.Face()
	}
	return nil
}
// SetFont sets the font face
func (w *Tag) SetFont(ft font.Face) {
	if f, ok := w.Window.(facer); ok{
		f.SetFont(ft)
		w.dirty = true
		w.Mark()
	}
}
