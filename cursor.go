package main

import "image"

type Cursor struct {
	sweep    bool
	sweepCol bool
	sr       image.Rectangle
	srcCol   *Col
	src      Plane
	dst      Plane
	detach   func()
}
