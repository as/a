package frame

import (
	"image"
	"image/color"
)

// copied from github.com/as/rgba to avoid a dependency

// hex converts a 32-bit RGBA quad to a color.RGBA
func hex(rgba uint32) color.RGBA {
	return color.RGBA{
		R: uint8(rgba >> 24),
		G: uint8(rgba << 8 >> 24),
		B: uint8(rgba << 16 >> 24),
		A: uint8(rgba << 24 >> 24),
	}
}

// uniform is short for image.NewUniform(hex(rgba))
func uniform(rgba uint32) *image.Uniform {
	return image.NewUniform(hex(rgba))
}
