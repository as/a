// Package ui is a wrapper around a graphical driver like shiny
// Basically, I don't like having two extra arguments per function
// so this is just an initialization package as well as a struct
// to hold the screen and window pointers.
package ui

import (
	"image"

	"github.com/as/shiny/driver"
	"github.com/as/shiny/screen"
)

type Item interface {
	Buffer() screen.Buffer
	Send(e interface{})
	SendFirst(e interface{})
	NextEvent() (e interface{})
}

type Win interface {
	Item
	Blank()
	Bounds() image.Rectangle
	Bytes() []byte
	Dirty() bool
	Fill()
	Len() int64
	Refresh()
	Size() image.Point
	Upload()
	Resize(size image.Point)
	Move(sp image.Point)
	Write(p []byte) (n int, err error)
}

type Dev interface {
	NewBuffer(size image.Point) (screen.Buffer, error)
	Screen() screen.Screen
	Window() screen.Window
}

type dev struct {
	scr    screen.Screen
	events screen.Window
	killc  chan bool
}

// Config configures the ui device
type Config = screen.NewWindowOptions

// TODO(as): write a gofix for this
// for now just type alias to avoid
// breaking compatibility with old
// programs

//type Config struct{
//	Width, Height int
//	Title         string
//}

func Init(conf *Config) (device Dev, err error) {
	errc := make(chan error)
	go func(errc chan error) {
		driver.Main(func(scr screen.Screen) {
			wind, err := scr.NewWindow((*screen.NewWindowOptions)(conf))
			if err != nil {
				errc <- err
			}
			d := &dev{scr, wind, make(chan bool)}
			device = d
			errc <- err
			<-d.killc
		})
	}(errc)
	return device, <-errc
}
func (d *dev) Screen() screen.Screen { return d.scr }
func (d *dev) Window() screen.Window { return d.events }
func (d *dev) NewBuffer(size image.Point) (screen.Buffer, error) {
	return d.scr.NewBuffer(size)
}

type Node struct {
	sp, size image.Point
}

func (n *Node) Move(pt image.Point) {
	n.sp = pt
}

func (n *Node) Resize(pt image.Point) {
	n.size = pt
}

func (n *Node) Size() image.Point {
	return n.size
}

func (n *Node) Pad() image.Point {
	return n.sp.Add(n.Size())
}
