package tag

import (
	"image"

	"github.com/as/shiny/screen"
)

type Window interface {
	Bytes() []byte
	Write(p []byte) (n int, err error)

	Select(q0, q1 int64)
	SetOrigin(org int64, exact bool)
	Dot() (int64, int64)
	Len() int64
	Origin() int64
	Graphical() bool
	//	IndexOf(image.Point) int64
	//	PointOf(int64) (image.Point)

	Scroll(dl int)
	Insert(p []byte, q0 int64) (n int)
	//	WriteAt(p []byte, at int64) (n int, err error)
	Delete(q0, q1 int64) (n int)

	Fill()
	Blank()
	Refresh()
	Dirty() bool
	Bounds() image.Rectangle
	Move(image.Point)
	Resize(image.Point)
	Close() error
	Upload()
	Window() screen.Window

	//	SetFont(font.Face)
}
