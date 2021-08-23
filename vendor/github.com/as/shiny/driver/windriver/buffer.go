// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package windriver

import (
	"image"
	"image/draw"
	"syscall"

	"github.com/as/shiny/driver/win32"
)

type bufferImpl struct {
	hbitmap   syscall.Handle
	buf, buf2 []byte
	rgba      image.RGBA
	size      image.Point
}

//func (b *bufferImpl) Resize(size image.Point) *bufferImpl       {
//	real := b.rgba.Bounds().Size()
//	curr := b.size
//	if size.X > real.X || size.Y > real.Y{
//		// Client wants a larger rectangle than we can provide
//		return nil
//	}
//	if curr.X*curr.Y / 3 > size.X*size.Y{
//		// Very small rectangle
//		returnnil
//	}
//}
func (b *bufferImpl) Size() image.Point       { return b.size }
func (b *bufferImpl) Bounds() image.Rectangle { return image.Rectangle{Max: b.size} }
func (b *bufferImpl) RGBA() *image.RGBA       { return &b.rgba }
func (b *bufferImpl) Release() {
	go b.cleanUp()
}

func (b *bufferImpl) cleanUp() {
	if b.rgba.Pix != nil {
		b.rgba.Pix = nil
		win32.DeleteObject(b.hbitmap)
	}
}

func (b *bufferImpl) blitToDC(dc syscall.Handle, dp image.Point, sr image.Rectangle) error {
	return copyBitmapToDC(dc, sr.Add(dp.Sub(sr.Min)), b.hbitmap, sr, draw.Src)
}
