package col

import (
	"image"

	"github.com/as/font"
	"github.com/as/ui"
	"github.com/as/ui/tag"
)

type Table2 struct {
	dev  ui.Dev
	ft   font.Face
	sp   image.Point
	size image.Point
	tdy  int

	Config *tag.Config
	Table
}

func NewTable2(dev ui.Dev, conf *tag.Config) Table2 {
	t := tag.New(dev, conf)
	return Table2{
		dev:    dev,
		ft:     conf.Facer(conf.FaceHeight),
		tdy:    conf.TagHeight(),
		Table:  Table{Tag: t},
		Config: &t.Config,
	}
}

func (co *Table2) Bounds() image.Rectangle {
	if co == nil {
		return image.ZR
	}
	return image.Rectangle{co.sp, co.sp.Add(co.size)}
}

func (co *Table2) Move(sp image.Point) {
	delta := sp.Sub(co.sp)
	co.Tag.Move(co.Tag.Bounds().Min.Add(delta))
	for _, t := range co.List {
		t.Move(t.Bounds().Min.Add(delta))
	}
	co.sp = sp
}

func (co *Table2) Upload() {
	type Uploader interface {
		Upload()
	}
	for _, w := range co.List {
		w, _ := w.(Uploader)
		if w != nil {
			defer w.Upload()
		}
	}
	co.Tag.Upload()
}

func (c *Table2) Dev() ui.Dev                { return c.dev }
func (c *Table2) Face() font.Face            { return c.ft }
func (c *Table2) ForceSize(size image.Point) { c.size = size }

// SetFont sets the font of all applicable nodes in the
// table. Because font.Face is not safe to use concurrently,
// f must also be a font.Resizer, otherwise the call is a
// no-op.
func (c *Table2) SetFont(f font.Face) {
	if f, ok := f.(font.Cache); ok {
		for _, c := range c.List {
			if c, ok := c.(facer); ok {
				c.SetFont(f)
			}
		}
		return
	}

	res, ok := f.(font.Resizer)
	if !ok {
		return
	}

	dy := res.Height()
	for _, c := range c.List {
		if c, ok := c.(facer); ok {
			c.SetFont(res.New(dy))
		}
	}
}
