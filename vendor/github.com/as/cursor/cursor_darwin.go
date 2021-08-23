// +build darwin
// +build 386 amd64
// +build !ios

package cursor

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#include <pthread.h>
#include <stdint.h>
#include <stdlib.h>
void setmouse(float x, float y);
void acmesetmouse(float x, float y);
NSRect bounds();
NSRect winbounds();
CGFloat pixelScale();
CGFloat pointScale();
*/
import "C"

import (
	"image"
)

func ScalePix(n int) int {
	return int(float64(n) * float64(C.pixelScale()))
}

func Bounds() image.Rectangle {
	nr := C.bounds()
	r := image.Rect(
		int(nr.origin.x),
		int(nr.origin.y),
		int(nr.origin.x+nr.size.width),
		int(nr.origin.y+nr.size.height),
	)
	return r
}

func WinBounds() image.Rectangle {
	nr := C.winbounds()
	r := image.Rect(
		int(nr.origin.x),
		int(nr.origin.y),
		int(nr.origin.x+nr.size.width),
		int(nr.origin.y+nr.size.height),
	)
	fr := Bounds()
	//s := int(C.pixelScale())
	r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	r.Min.Y = fr.Max.Y + (-r.Min.Y)*2
	r.Max.Y = fr.Max.Y + (-r.Max.Y)*2
	r.Min.X *= 2
	r.Max.X *= 2
	return r
}

func MoveTo(p image.Point) bool {
	p = p.Div(2)
	p = p.Add(WinBounds().Min.Div(2))
	x := float32(p.X)
	y := float32(p.Y)
	C.setmouse(C.float(x+4), C.float(y+25+4))
	return true
}

func moveTo(p image.Point) bool {
	return false
}
