package win32

var ZP = Point{}

type Point struct {
	X int32
	Y int32
}

func Pt(x, y int) Point {
	return Point{int32(x), int32(y)}
}
