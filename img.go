// Image rendering
package main

import (
	"path/filepath"
	"strings"

	"github.com/as/ui/img"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"

	_ "image/gif"
	_ "image/jpeg"

	_ "golang.org/x/image/bmp"
)

var imageExt = map[string]bool{
	".png":  true,
	".jpeg": true,
	".jpg":  true,
	".gif":  true,
	".bmp":  true,
}

func tryImage(name string) bool {
	return imageExt[strings.ToLower(filepath.Ext(name))]
}

func rendercmd(t *tag.Tag) {
	switch t.Window.(type) {
	case *win.Win:
		// Currently text; convert to image
		render(t)
	case interface{}:
		// Currently some other format, convert to text
		unrender(t)
	}
}

func renderimage(t *tag.Tag) {
	switch t.Window.(type) {
	case *win.Win:
		var oldb = t.Window
		t.Config.Image = true
		t.Window = img.New(t.Label.Dev, nil) // TODO(as): how to get rid of Dev here
		t.Insert(oldb.Bytes(), 0)
		t.Window.Move(oldb.Bounds().Min)
		t.Window.Resize(oldb.Bounds().Size())

		if !t.Window.Graphical() && t.Window != oldb {
			t.Window = oldb
			t.Config.Image = false
		}
		return
	}

	unrender(t)
}

func render(t *tag.Tag) {
	var oldb = t.Window
	if tryImage(t.FileName()) {
		t.Config.Image = true
		t.Window = img.New(t.Label.Dev, nil)
		t.Get(t.FileName())
		t.Window.Move(oldb.Bounds().Min)
		t.Window.Resize(oldb.Bounds().Size())
	}

	// If it didn't render correctly; restore
	// the text buffer
	if !t.Window.Graphical() && t.Window != oldb {
		t.Window = oldb
		t.Config.Image = false
		t.Get(t.FileName())
	}
}

func unrender(t *tag.Tag) {
	var oldb = t.Window
	t.Config.Image = false
	t.Window = win.New(t.Label.Dev, &t.Config.Body)
	t.Get(t.FileName())
	t.Window.Move(oldb.Bounds().Min)
	t.Window.Resize(oldb.Bounds().Size())
}
