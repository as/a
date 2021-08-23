package col

import (
	"image"

	"github.com/as/font"
	"github.com/as/frame"
	"github.com/as/rgba"
	"github.com/as/ui"
	"github.com/as/ui/tag"
)

type Col struct {
	Table2
}

var DefaultConfig = &tag.Config{
	Margin:     image.Pt(15, 0),
	Facer:      font.NewFace,
	FaceHeight: 11,
	Color: [3]frame.Color{
		frame.Theme(rgba.Gray, rgba.Strata, rgba.White, rgba.Mauve),
	},
	Ctl: make(chan interface{}, 10),
}

// New creates a new column with the device, font, source point
// and size.
func New(dev ui.Dev, conf *tag.Config) *Col {
	if conf == nil {
		conf = DefaultConfig
	}
	return &Col{
		Table2: NewTable2(dev, conf),
	}
}

func (co *Col) Resize(size image.Point) {
	if size.X < 0 || size.Y < 0 {
		return
	}
	co.size = size
	size.Y = co.tdy
	co.Tag.Resize(size)
	Fill(co)
}

func (co *Col) RollDown(id int, dy int) {
}
func (co *Col) MoveWin(id int, y int) {
	if id >= len(co.List) {
		return
	}
	//FLAG
	maxy := co.Bounds().Max.Y - co.tdy
	if y >= maxy {
		return
	}
	Attach(co, co.detach(id), image.Pt(0, y))
}

func (co *Col) Grow(id int, dy int) {
	a, b := id-1, id
	if co.badID(a) || co.badID(b) {
		return
	}
	ra, rb := co.List[a].Bounds(), co.List[b].Bounds()
	ra.Max.Y -= dy
	if dy := ra.Dy() - co.tdy; dy < 0 {
		co.Grow(a, -dy)
	}
	co.MoveWin(b, rb.Min.Y-dy)
}

func (co *Col) RollUp(id int, dy int) {
	if id <= 0 || id >= len(co.List) {
		return
	}
	pt := co.Tag.Bounds().Min
	pt.Y += co.tdy
	for x := 1; x <= id; x++ {
		pt.Y += co.tdy
		co.List[x].Move(pt)
	}
	Fill(co)
	co.Refresh()
}
