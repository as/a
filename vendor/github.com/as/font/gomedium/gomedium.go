package gomedium

import (
	"log"

	"github.com/golang/freetype/truetype"
	. "golang.org/x/image/font/gofont/gomedium"
)

var Font, err = truetype.Parse(TTF)

func init() {
	if err != nil {
		log.Fatalln("gomedium", err)
	}
}
