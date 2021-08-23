package tag

import (
	"bytes"
	"errors"
	"image"
	"io"
	"runtime"

	"github.com/as/path"
	"github.com/as/shiny/event/mouse"
	"github.com/as/text"
	"github.com/as/text/action"
	"github.com/as/ui/win"
)

var ErrNoFile = errors.New("no file")

func Pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

// This is actually the version that is preferrable to use. In the interest
// of time--I'm deferring its use to a future date.
func (t *Tag) zGet(name string) error {
	data, err := t.Fs.Get(name)
	if err != nil {
		return err
	}
	w := t
	w.Delete(0, t.Len())
	w.Insert(data, 0)
	w.SetOrigin(0, true)
	w.Select(0, 0)
	t.fixtag(name)
	return nil
}

// Get retrieves the named resource. It uses the local disk only
// TODO(as): fix the semantics of this method. It should use the
// fs resolver and properly handle paths to arbitrary resources
func (t *Tag) Get(name string) {
	if name == "" {
		t.fixtag("")
		return
	}
	name, addr := action.SplitPath(name)
	if IsAbs(name) && path.Exists(name) {
		t.fixtag(name)
		t.getbody(name, addr)
		return
	}
	t.fixtag(name)
	t.getbody(name, addr)
}

func (t *Tag) Put() (err error) {
	return t.Fs.Put(t.FileName(), t.Bytes())
}

func (t *Tag) getbody(abs, addr string) {
	w := t
	w.Delete(0, w.Len())
	w.Insert(t.readfile(abs), 0)
	w.Select(0, 0)
	w.SetOrigin(0, true)
	if addr != "" {
		t.ctl <- mustCompile(addr)
		//w.SendFirst(mustCompile(addr)) //TODO
	}
}

func Visible(w *win.Win, q0, q1 int64) bool {
	if q0 < w.Origin() {
		return false
	}
	if q1 > w.Origin()+w.Frame.Nchars {
		return false
	}
	return true
}

// Snarf -- TODO(as): move snarf out of here, it doesn't belong in this package
func Snarf(w text.Editor, e mouse.Event) {
	var n int
	if runtime.GOOS == "windows" {
		n = copy(ClipBuf, toUTF16([]byte(Rdsel(w))))
	} else {
		n = copy(ClipBuf, []byte(Rdsel(w)))
	}
	io.Copy(Clip, bytes.NewReader(ClipBuf[:n]))
	q0, q1 := w.Dot()
	w.Delete(q0, q1)
	w.Select(q0, q0)
}

// Paste -- TODO(as): move paste out of here, it doesn't belong in this package
func Paste(w text.Editor, e mouse.Event) (int64, int64) {
	n, _ := Clip.Read(ClipBuf)
	s := ClipBuf[:n]
	if runtime.GOOS == "windows" {
		s = fromUTF16(s)
	} else {
		s = append([]byte{}, s...)
	}
	q0, q1 := w.Dot()
	if q0 != q1 {
		w.Delete(q0, q1)
		q1 = q0
	}
	w.Insert(s, q0)
	w.Select(q0, q0+int64(len(s)))
	return w.Dot()
}

func Rdsel(w text.Editor) string {
	q0, q1 := w.Dot()
	return string(w.Bytes()[q0:q1])
}
