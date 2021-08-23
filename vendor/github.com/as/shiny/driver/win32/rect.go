package win32

type Rectangle struct {
	Min, Max Point
}

func (r Rectangle) Dx() int32 {
	return r.Max.X - r.Min.X
}
func (r Rectangle) Dy() int32 {
	return r.Max.Y - r.Min.Y
}

func Rect(x, y, xx, yy int) Rectangle {
	return Rectangle{
		Min: Point{
			int32(x),
			int32(y),
		},
		Max: Point{
			int32(xx),
			int32(yy),
		},
	}
}
