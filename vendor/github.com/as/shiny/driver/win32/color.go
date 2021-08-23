package win32

import "unsafe"

type BlendFunc struct {
	Op            byte
	Flags         byte
	SrcConstAlpha byte
	AlphaFormat   byte
}

// Uintptr helps to pass bf to syscall.Syscall.
func (b BlendFunc) Uintptr() uintptr {
	return *((*uintptr)(unsafe.Pointer(&b)))
}

type RGBQuad struct {
	Blue     byte
	Green    byte
	Red      byte
	Reserved byte
}

type ColorRef uint32

func NewColorRef(c interface {
	RGBA() (r, g, b, a uint32)
}) ColorRef {
	type cr = ColorRef
	r, g, b, a := c.RGBA()
	return cr(r>>8) | cr(g>>8)<<8 | cr(b>>8)<<16 | cr(a>>8)<<24
}

func RGB(r, g, b byte) ColorRef {
	return ColorRef(r) | ColorRef(g)<<8 | ColorRef(b)<<16
}
