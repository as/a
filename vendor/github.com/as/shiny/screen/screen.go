// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package screen // import "github.com/as/shiny/screen"

import (
	"image"
	"image/color"
	"image/draw"
	"unicode/utf8"

	"github.com/as/shiny/event/key"
	"github.com/as/shiny/event/lifecycle"
	"github.com/as/shiny/event/mouse"
	"github.com/as/shiny/event/paint"
	"github.com/as/shiny/event/size"
	"github.com/as/shiny/math/f64"
)

// Screen creates Buffers, Textures and Windows.
type Screen interface {
	NewBuffer(size image.Point) (Buffer, error)
	NewTexture(size image.Point) (Texture, error)
	NewWindow(opts *NewWindowOptions) (Window, error)
}

type Buffer interface {
	Release()
	Size() image.Point
	Bounds() image.Rectangle
	RGBA() *image.RGBA
}

type Texture interface {
	Release()
	Size() image.Point
	Bounds() image.Rectangle
	Uploader
}

// Window is a top-level, double-buffered GUI window.
type Window interface {
	Device() *Device
	Release()
	Uploader
	Drawer
	Publish() PublishResult
}

type (
	Lifecycle = lifecycle.Event
	Scroll    = mouse.Event
	Mouse     = mouse.Event
	Key       = key.Event
	Size      = size.Event
	Paint     = paint.Event
)

var Dev = &Device{
	Scroll:    make(chan Scroll, 1),
	Mouse:     make(chan Mouse, 1),
	Key:       make(chan Key, 1),
	Size:      make(chan Size, 1),
	Paint:     make(chan Paint, 1),
	Lifecycle: make(chan Lifecycle, 1),
}

type Device struct {
	Lifecycle chan Lifecycle
	Scroll    chan Scroll
	Mouse     chan Mouse
	Key       chan Key
	Size      chan Size
	Paint     chan Paint
}

func SendMouse(e Mouse) {
	select {
	case Dev.Mouse <- e:
	default:
		// TODO: Retry on failure, but only if it's a press or release
		// note that this may hang the user, so a better fix should be
		// in order
		if e.Button != mouse.ButtonNone && e.Direction != mouse.DirNone {
			Dev.Mouse <- e
		}
	}
}

func SendKey(e Key) {
	select {
	case Dev.Key <- e:
	}
}

func SendSize(e Size) {
	select {
	case Dev.Size <- e:
	default:
	}
}

func SendPaint(e Paint) {
	select {
	case Dev.Paint <- e:
	default:
	}
}

func SendScroll(e Scroll) {
	for {
		select {
		case Dev.Scroll <- e:
			return
		default:
			select {
			case <-Dev.Scroll:
			default:
			}
		}
	}
}

func SendLifecycle(e Lifecycle) {
	select {
	case Dev.Lifecycle <- e:
	}

}

// PublishResult is the result of an Window.Publish call.
type PublishResult struct {
	BackBufferPreserved bool
}

// NewWindowOptions are optional arguments to NewWindow.
// TODO(as): NewWindowOptions could be named better
type NewWindowOptions struct {
	Width, Height int
	Title         string

	// Overlay, if true, attempts to create a new window over top
	// of the parent process's existing window (similar to running
	// a graphical application in Plan9 over top an existing Rio
	// window).
	Overlay bool
}

func (o *NewWindowOptions) GetTitle() string {
	if o == nil {
		return ""
	}
	return sanitizeUTF8(o.Title, 4096)
}

func sanitizeUTF8(s string, n int) string {
	if n < len(s) {
		s = s[:n]
	}
	i := 0
	for i < len(s) {
		r, n := utf8.DecodeRuneInString(s[i:])
		if r == 0 || (r == utf8.RuneError && n == 1) {
			break
		}
		i += n
	}
	return s[:i]
}

// Uploader is something you can upload a Buffer to.
type Uploader interface {
	Upload(dp image.Point, src Buffer, sr image.Rectangle)
	Fill(dr image.Rectangle, src color.Color, op draw.Op)
}

type Drawer interface {
	Draw(src2dst f64.Aff3, src Texture, sr image.Rectangle, op draw.Op, opts *DrawOptions)
	DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op, opts *DrawOptions)
	Copy(dp image.Point, src Texture, sr image.Rectangle, op draw.Op, opts *DrawOptions)
	Scale(dr image.Rectangle, src Texture, sr image.Rectangle, op draw.Op, opts *DrawOptions)
}

const (
	Over = draw.Over
	Src  = draw.Src
)

type DrawOptions struct {
}
