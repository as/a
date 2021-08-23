package font

import (
	"unicode"

	"github.com/as/font/gomedium"
	"github.com/as/font/gomono"
	"github.com/as/font/goregular"
	"github.com/golang/freetype/truetype"
)

func notGraphic(r rune) bool { return !unicode.IsGraphic(r) }

// New returns a GoMedium font that is cached and replaces
// non-printable characters with their hex equivalent encodings
func NewFace(size int) Face {
	return NewCache(Replacer(NewGoMedium(size), NewHex(size), nil))
	// return NewCache(Replacer(NewGoMedium(size), NewHex(size), nil))
	//	return NewCache(NewRune(NewGoMedium(size)))
	//	var fn func(size int) Face
	//	fn = func(size int) Face{
	//		return &Resizer{
	//			Face: NewCache(Replacer(NewGoMedium(size), NewHex(size), nil)),
	//			New: fn,
	//		}
	//	}
	//	return fn(size)
}

func NewGoMedium(size int) Face {
	return Open(truetype.NewFace(gomedium.Font, &truetype.Options{
		SubPixelsX: 64,
		SubPixelsY: 64,
		Size:       float64(size),
	}))
}

func NewGoMono(size int) Face {
	return Open(truetype.NewFace(gomono.Font, &truetype.Options{
		SubPixelsX: 64,
		SubPixelsY: 64,
		Size:       float64(size),
	}))
}

func NewGoRegular(size int) Face {
	return Open(truetype.NewFace(goregular.Font, &truetype.Options{
		SubPixelsX: 64,
		SubPixelsY: 64,
		Size:       float64(size),
	}))
}
