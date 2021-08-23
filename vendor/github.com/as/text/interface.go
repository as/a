package text

import (
	"image"

	"github.com/as/event"
)

type Selector interface {
	Select(q0, q1 int64)
	Dot() (q0, q1 int64)
}

type File interface {
	Named
	Editor
	Put() error
	Get() error
}

type Named interface {
	Name() string
}

type Editor interface {
	Buffer
	Selector
}

type Dirt interface {
	Mark()
	Dirty() bool
}

type Plane interface {
	Bounds() image.Rectangle
	Size() image.Point
}

type Projector interface {
	PointOf(q int64) (pt image.Point)
	IndexOf(pt image.Point) (q int64)
}

type Ruler interface {
	Measure(s string) int
	Height() int
}

type Shower interface {
	Show(min int) bool
}
type Jumper interface {
	Jump(mouseFunc func(image.Point))
}

type Win interface {
	Editor
	Scroller
	Plane
	Mark()
}

type Sweeper interface {
	Plane
	Projector
	Scroller
	Selector
}

type Scroller interface {
	Origin() int64
	SetOrigin(int64, bool)
	Fill()
	Scroll(int)
}

type Inverse struct {
	e interface{}
}

type History interface {
	Next() interface{}
	Prev() interface{}
	Event() interface{}
	Add(e interface{})
	Commit()
	Apply()
}

type Logger interface {
	Write(event.Record) (err error)
	ReadAt(at int64) (event.Record, error)
	Len() int64
}

type Sender interface {
	Send(interface{})
	SendFirst(interface{})
}
