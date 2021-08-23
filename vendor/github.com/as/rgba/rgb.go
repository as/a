// Package rgba provides a conversion between 32-bit RGBA quad values
// and the standard library's image and color packages. This package
// defines no types of its own and is used primarily to translate from
// a common web notation to the standard library
package rgba

import (
	"image"
	"image/color"
)

// Hex converts a 32-bit RGBA quad to a color.RGBA
func Hex(rgba uint32) color.RGBA {
	return color.RGBA{
		R: uint8(rgba >> 24),
		G: uint8(rgba << 8 >> 24),
		B: uint8(rgba << 16 >> 24),
		A: uint8(rgba << 24 >> 24),
	}
}

func Plan9(c color.Color) color.Color {
	return c
	//	return palette.Plan9[color.Palette(palette.Plan9).Index(c)]
}

// Uniform is short for image.NewUniform(Hex(rgba))
func Uniform(rgba uint32) *image.Uniform {
	return image.NewUniform(Plan9(Hex(rgba)))
}

// Uint32 converts a color.RGBA to a uint32
func Uint32(c color.RGBA) uint32 {
	return uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
}
