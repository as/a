package tag

import (
	"image"

	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/srv/fs"
	"github.com/as/ui/scroll"
	"github.com/as/ui/win"
)

type Config struct {
	Facer      func(int) font.Face
	FaceHeight int
	Margin     image.Point
	Color      [3]frame.Color
	Tag        win.Config
	Body       win.Config
	Ctl        chan interface{}
	Filesystem fs.Fs
	Image      bool // Image makes win.Win ->img.Img instead
}

func validConfig(c *Config) *Config {
	if c == nil {
		return &DefaultConfig
	}
	if c.Ctl == nil {
		panic("ctl cant be nil")
	}
	if c.Filesystem == nil {
		c.Filesystem = &fs.Local{}
	}
	if c.FaceHeight == 0 {
		c.FaceHeight = DefaultConfig.FaceHeight
	}
	return c
}

func (c *Config) TagHeight() int {
	return Height(c.FaceHeight)
}

func (c *Config) TagConfig() *win.Config {
	return &win.Config{
		Ctl:    c.Ctl,
		Facer:  c.Facer,
		Margin: c.Margin,
		Scroll: scroll.Config{Enable: false},
		Frame: frame.Config{
			Color: c.Color[0],
			Face:  c.Facer(c.FaceHeight),
		},
	}
}
func (c *Config) WinConfig() *win.Config {
	return &win.Config{
		Ctl:    c.Ctl,
		Facer:  c.Facer,
		Margin: c.Margin,
		Scroll: scroll.Config{
			Enable: true,
			Color: [2]image.Image{
				c.Color[2].Text,
				c.Color[2].Back,
			},
		},
		Frame: frame.Config{
			Color: c.Color[1],
			Face:  c.Facer(c.FaceHeight),
			Flag:  c.Body.Frame.Flag,
		},
	}
}
