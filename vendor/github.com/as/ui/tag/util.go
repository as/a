package tag

import (
	"bytes"
	"fmt"

	"github.com/as/clip"
	"github.com/as/frame"
	"github.com/as/ui/win"
)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func (t *Tag) dirfmt(p []byte) []byte {
	win, _ := t.Window.(*win.Win)
	if win == nil {
		return nil
	}
	x, dx, maxx := 0, 8, 600
	win.Config.Frame.Flag |= frame.FrElastic
	if win.Graphical() {
		dx = win.Face().Dx([]byte{'e'}) // common lowercase rune
		maxx = win.Frame.Bounds().Dx()
		fl := win.Frame.Flags() | frame.FrElastic
		t.Config.Body.Frame.Flag = fl
		win.Frame.SetFlags(fl)
	}
	w := new(bytes.Buffer)
	for _, nm := range bytes.Split(p, []byte{'\t'}) {
		word := fmt.Sprintf("\t%s", nm)
		wordlen := len(word) - 1
		wordpix := wordlen * dx
		advance := max(wordpix, 8*x)
		if x+advance > maxx {
			fmt.Fprintf(w, "\t\n")
			x = -advance
		}
		fmt.Fprintf(w, word)
		x += advance
	}
	return w.Bytes()
}

func (t *Tag) readfile(s string) (p []byte) {
	fi, err := t.Fs.Stat(s)
	dir := false
	if err == nil && fi.IsDir() {
		dir = true
	}

	p, err = t.Fs.Get(s)
	if err != nil {
		return []byte{}
	}
	if dir {
		p = t.dirfmt(p)
	}
	return p
}
func (t *Tag) writefile(s string, p []byte) {
	err := t.Fs.Put(s, p)
	if err != nil {
		t.ctl <- err
	}
}

func init() {
	var err error
	Clip, err = clip.New()
	if err != nil {
		panic(err)
	}
}
