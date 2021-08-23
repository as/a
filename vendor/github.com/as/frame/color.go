package frame

import (
	"image"
)

var (
	// Common uniform colors found in Acme
	Black  = image.Black
	White  = image.White
	Yellow = uniform(0xfffffdff)
	Red    = uniform(0xffe8efff)
	Green  = uniform(0xefffe8ff)
	Blue   = uniform(0xe8efffff)

	// Other colors
	Gray  = uniform(0x1c1f26ff)
	Peach = uniform(0xfff8e8ff)
	Mauve = uniform(0x9090C0ff)
)

var (
	// Acme is the color scheme found in the Acme text editor
	Acme = Theme(Gray, Yellow, White, Blue)
	Mono = Theme(Black, White, White, Black)
	A    = Theme(Gray, Peach, White, Mauve)
)

// Color is constructed from a Palette pair. The Hi Palette describes
// the appearance of highlighted text.
type Color struct {
	Palette
	Hi Palette
}

// Pallete contains two images used to paint text and backgrounds
// on the frame.
type Palette struct {
	Text, Back image.Image
}

// Theme returns a Color for the given foreground and background
// images. Two extra colors may be provided to set the highlighted
// foreground and background image palette.
func Theme(fg, bg image.Image, hi ...image.Image) Color {
	c := Color{Palette: Palette{Text: fg, Back: bg}}
	if len(hi) > 0 {
		c.Hi.Text = hi[0]
	}
	if len(hi) > 1 {
		c.Hi.Back = hi[1]
	}
	return c
}
