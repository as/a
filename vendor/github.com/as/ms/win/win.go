package win

import (
	"image"
	"os"
)

// Rect returns the windows bounds on screen
func Rect() (image.Rectangle, bool) {
	w, err := Open(os.Getpid())
	if err != nil {
		return image.ZR, false
	}
	r, err := w.Bounds()
	if err != nil {
		return image.ZR, false
	}
	return r, true
}

// ClientRect returns the sub-rectangle of the client area (the intersection of
// the window border and the window). The bounds are relative to the window bounds,
// so this rectangle remains constant if the window is moved without resize
func ClientRect() (image.Rectangle, bool) {
	w, err := Open(os.Getpid())
	println("pid is", os.Getpid(), "and w is", w)
	if err != nil {
		panic(err.Error())
		return image.ZR, false
	}
	r, err := w.Client()
	if err != nil {
		panic(err.Error())
		return image.ZR, false
	}
	return r, true
}

// ClientAbs returns the absolute client area.
func ClientAbs() image.Rectangle {
	wr, _ := Rect()
	cr, _ := ClientRect()
	return wr.Add(Border(wr, cr))
}

func Border(wr, cr image.Rectangle) image.Point {
	return image.Point{
		X: wr.Max.X - cr.Max.X - wr.Min.X,
		Y: wr.Max.Y - cr.Max.Y - wr.Min.Y,
	}
}
func FromPID(p int) (wids []uintptr, err error) {
	return fromPID(p)
}
func Move(wid int, to image.Rectangle, paint bool) (err error) {
	return move(wid, to, paint)
}
