package frame

import (
	"github.com/as/font"
)

type Config struct {
	Flag   int
	Scroll func(int)
	Color  Color
	Face   font.Face
	Drawer Drawer
}

func (c *Config) check() *Config {
	if c.Color == zc {
		c.Color = A
	}
	if c.Face == nil {
		c.Face = font.NewFace(11)
	}
	if c.Drawer == nil {
		c.Drawer = &defaultDrawer{}
	}
	return c
}

func getflag(flag ...int) (fl int) {
	if len(flag) != 0 {
		fl = flag[0]
	}
	if ForceElastic {
		fl |= FrElastic
	}
	if ForceUTF8 {
		fl |= FrUTF8
	}
	return fl
}
