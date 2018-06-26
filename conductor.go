package main

/*
import (
	"errors"
	"image"
	"strings"

	"github.com/as/ui/col"
	"github.com/as/ui/tag"
)

var (
	ErrPlacement = errors.New("placement error")
	ErrUnknown   = errors.New("unknown placement request")
)

type evp struct {
	dst, src Plane
	root     Plane
}

var (
	evpI = make(chan evp)
	evpO = make(chan error)
)

func conduct() {
	for {
		select {
		case evp := <-evpI:
			select{
			case evpO <- place(evp):
			}
		}
	}
}

func place(e evp) error {
	switch t := e.src.(type) {
	case *tag.Tag:
		g, _ := e.root.(*Grid)
		c, _ := e.dst.(*col.Col)
		if c == nil {
			return ErrPlacement
		}
		if g != nil && strings.HasSuffix(t.FileName(), "+Errors") && len(g.List) > 1 {
			c0 := g.List[len(g.List)-1].(*Col)
			if c0 != nil {
				c = c0
			}
		}

		dDY := c.Area().Dy()
		if c.Len() > 0 {
			last := c.List[len(c.List)-1]
			if c.Area().Dy() > last.Bounds().Dy()*3 {
				logf("should roll up windows--area too small")
			}
			dDY = last.Bounds().Dy() / 2
		}
		col.Attach(c, t, image.Pt(dDY, dDY))
		return nil
	}
	return ErrUnknown
}
*/
