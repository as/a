// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package windriver

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/as/shiny/driver/win32"
	"github.com/as/shiny/screen"
)

var theScreen = &screenImpl{}

type screenImpl struct {
	windows *windowImpl
}

func (*screenImpl) NewBuffer(size image.Point) (screen.Buffer, error) {
	const (
		maxInt32  = 0x7fffffff
		maxBufLen = maxInt32
	)
	if size.X < 0 || size.X > maxInt32 || size.Y < 0 || size.Y > maxInt32 || int64(size.X)*int64(size.Y)*4 > maxBufLen {
		return nil, fmt.Errorf("windriver: invalid buffer size %v", size)
	}

	hbitmap, bitvalues, err := mkbitmap(size)
	if err != nil {
		return nil, err
	}
	bufLen := 4 * size.X * size.Y
	array := (*[maxBufLen]byte)(unsafe.Pointer(bitvalues))
	buf := (*array)[:bufLen:bufLen]
	buf2 := make([]byte, bufLen, bufLen)
	return &bufferImpl{
		hbitmap: hbitmap,
		buf:     buf,
		buf2:    buf2,
		rgba: image.RGBA{
			Pix:    buf2,
			Stride: 4 * size.X,
			Rect:   image.Rectangle{Max: size},
		},
		size: size,
	}, nil
}

func (*screenImpl) NewTexture(size image.Point) (screen.Texture, error) {
	return newTexture(size)
}

func (s *screenImpl) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	h, err := win32.NewWindow(opts)
	if err != nil {
		return nil, err
	}
	dc, err := win32.GetDC(h)
	if err != nil {
		return nil, err
	}
	_, err = win32.SetGraphicsMode(dc, win32.GmAdvanced)
	if err != nil {
		return nil, err
	}

	s.windows = &windowImpl{
		dc:   dc,
		hwnd: h,
	}
	if err = win32.Resize(h, win32.Pt(opts.Width, opts.Height)); err != nil {
		return nil, err
	}
	win32.Show(h)
	return s.windows, nil
}

/*
// experimental: attempt to overlay on top of existing window
func (s *screenImpl) newWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	var (
		err error
	)
	h := win32.GetConsoleWindow()
	if !opts.Overlay || h == 0 {
		// Won't bind to the parent window; create a new one
		if h, err = win32.NewWindow(opts); err != nil {
			return nil, err
		}
	}
	dc, err := win32.GetDC(h)
	if err != nil {
		return nil, err
	}
	_, err = _SetGraphicsMode(dc, _GM_ADVANCED)
	if err != nil {
		return nil, err
	}

	s.windows = &windowImpl{
		dc:   dc,
		hwnd: h,
	}
	if err = win32.ResizeClientRect(h, opts); err != nil {
		return nil, err
	}
	win32.Show(h)
	return s.windows, nil
}
*/
