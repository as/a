package main

import "image"

// Area returns the bounds in which the list elements can reside
func (c *Grid) Area() image.Rectangle {
	dy := c.Tag.Bounds().Dy()
	r := c.Bounds()
	r.Min.Y += dy
	return r
}

// Minor returns the a point in Grid where Y is top-aligned
// and X is clamped between min.X and max.X
func (c *Grid) Minor(pt image.Point) image.Point {
	pt.Y = c.Area().Min.Y
	pt.X = clampx(pt.X, c.Area().Min.X, c.Area().Max.X)
	return pt
}

// Major returns the a point in Col where X is right-aligned
// and X is clamped between min.X and max.X
func (c *Grid) Major(pt image.Point) image.Point {
	pt.Y = c.Area().Max.Y                               //-c.Area().Min.Y
	pt.X = clampx(pt.X, c.Area().Min.X, c.Area().Max.X) //-c.Area().Min.X
	return pt
}

func clampx(v, l, h int) int {
	if v < l {
		return l
	}
	if v > h {
		return h
	}
	return v
}
