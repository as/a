package font

// Resizer knows how to return an enlarged version of its face
type Resizer struct {
	Face

	// New returns a new face with the given size. The current
	// size is not closed and is still usable by callers
	New func(int) Face
}
