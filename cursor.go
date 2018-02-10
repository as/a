package main


type Cursor struct {
	sweep    bool
	sweepCol bool
	srcCol   *Col
	src      Plane
}
