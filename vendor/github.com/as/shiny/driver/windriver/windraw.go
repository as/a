// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package windriver

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"syscall"
	"unsafe"

	"github.com/as/shiny/driver/win32"
)

func mkbitmap(size image.Point) (syscall.Handle, *byte, error) {
	bi := win32.BitmapInfo{
		Header: win32.BitmapInfoV4{
			Size:        uint32(unsafe.Sizeof(win32.BitmapInfoV4{})),
			Width:       int32(size.X),
			Height:      -int32(size.Y), // negative height to force top-down drawing
			Planes:      1,
			BitCount:    32,
			Compression: win32.BiRGB,
			SizeImage:   uint32(size.X * size.Y * 4),
		},
	}

	var ppvBits *byte
	bitmap, err := win32.CreateDIBSection(0, &bi, win32.DibRGBColors, &ppvBits, 0, 0)
	if err != nil {
		return 0, nil, err
	}
	return bitmap, ppvBits, nil
}

var blendOverFunc = win32.BlendFunc{
	Op:            win32.AcSrcOver,
	Flags:         0,
	SrcConstAlpha: 255,              // only use per-pixel alphas
	AlphaFormat:   win32.AcSrcAlpha, // premultiplied
}

func copyBitmapToDC(dc syscall.Handle, dr image.Rectangle, src syscall.Handle, sr image.Rectangle, op draw.Op) (retErr error) {
	memdc, err := win32.CreateCompatibleDC(dc)
	if err != nil {
		return err
	}
	defer win32.DeleteDC(memdc)

	_, err = win32.SelectObject(memdc, src)
	if err != nil {
		return err
	}

	switch op {
	case draw.Src:
		return win32.StretchBlt(dc, int32(dr.Min.X), int32(dr.Min.Y), int32(dr.Dx()), int32(dr.Dy()),
			memdc, int32(sr.Min.X), int32(sr.Min.Y), int32(sr.Dx()), int32(sr.Dy()), win32.SrcCopy)
	case draw.Over:
		return win32.AlphaBlend(dc, int32(dr.Min.X), int32(dr.Min.Y), int32(dr.Dx()), int32(dr.Dy()),
			memdc, int32(sr.Min.X), int32(sr.Min.Y), int32(sr.Dx()), int32(sr.Dy()), blendOverFunc.Uintptr())
	default:
		return fmt.Errorf("windriver: invalid draw operation %v", op)
	}
}

func fill(dc syscall.Handle, dr image.Rectangle, c color.Color, op draw.Op) error {
	const bgrmask = ((1 << 24) - 1)
	cr := win32.NewColorRef(c)

	if op == draw.Src {
		brush, err := win32.CreateSolidBrush(cr & bgrmask)
		if err != nil {
			return err
		}
		defer win32.DeleteObject(brush)

		return win32.FillRect(dc, &win32.Rectangle{
			win32.Point{int32(dr.Min.X), int32(dr.Min.Y)},
			win32.Point{int32(dr.Max.X), int32(dr.Max.Y)},
		}, brush)
	}

	// AlphaBlend will stretch the input image (using StretchBlt's
	// COLORONCOLOR mode) to fill the output rectangle. Testing
	// this shows that the result appears to be the same as if we had
	// used a MxN bitmap instead.
	sr := image.Rectangle{image.ZP, image.Point{1, 1}}
	bitmap, bitvalues, err := mkbitmap(sr.Max)
	if err != nil {
		return err
	}
	defer win32.DeleteObject(bitmap)

	*(*win32.ColorRef)(unsafe.Pointer(bitvalues)) = cr

	return copyBitmapToDC(dc, dr, bitmap, sr, draw.Over)
}
