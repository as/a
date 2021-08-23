package col

import (
	"image"
)

type Plane interface {
	Bounds() image.Rectangle
	Move(image.Point)
	Resize(image.Point)
	Dirty() bool
	Refresh()
}
