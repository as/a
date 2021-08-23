package col

import "image"

type Axis interface {
	Major(image.Point) image.Point
	Minor(pt image.Point) image.Point
	Area() image.Rectangle
}

// Area returns the bounds in which the list elements can reside
func (c *Col) Area() image.Rectangle {
	dy := c.Tag.Bounds().Dy()
	r := c.Bounds()
	r.Min.Y += dy
	return r
}

// Minor returns the a point in Col where X is left-aligned
// and Y is clamped between min.Y and max.Y
func (c *Col) Minor(pt image.Point) image.Point {
	pt.X = c.Area().Min.X
	pt.Y = clamp(pt.Y, c.Area().Min.Y, c.Area().Max.Y)
	return pt
}

// Major returns the a point in Col where X is right-aligned
// and Y is clamped between min.Y and max.Y
func (c *Col) Major(pt image.Point) image.Point {
	pt.X = c.Area().Max.X
	pt.Y = clamp(pt.Y, c.Area().Min.Y, c.Area().Max.Y)
	return pt
}

func clamp(v, l, h int) int {
	if v < l {
		return l
	}
	if v > h {
		return h
	}
	return v
}
