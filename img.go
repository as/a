// Image rendering
package main

import (
	"path/filepath"
	"strings"

	"github.com/as/ui/img"
	"github.com/as/ui/tag"
	"github.com/as/ui/win"

	_ "golang.org/x/image/bmp"
	_ "image/gif"
	_ "image/jpeg"
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
	switch t.Body.(type) {
	case *win.Win:
		// Currently text; convert to image
		render(t)
	case interface{}:
		// Currently some other format, convert to text
		unrender(t)
	}
}

func renderimage(t *tag.Tag) {
	switch t.Body.(type) {
	case *win.Win:
		var oldb = t.Body
		t.Config.Image = true
		t.Body = img.New(t.Dev, nil)
		t.Body.Insert(oldb.Bytes(), 0)
		t.Body.Move(oldb.Bounds().Min)
		t.Body.Resize(oldb.Bounds().Size())

		if !t.Body.Graphical() && t.Body != oldb {
			t.Body = oldb
			t.Config.Image = false
		}
		return
	}

	unrender(t)
}

func render(t *tag.Tag) {
	var oldb = t.Body
	if tryImage(t.FileName()) {
		t.Config.Image = true
		t.Body = img.New(t.Dev, nil)
		t.Get(t.FileName())
		t.Body.Move(oldb.Bounds().Min)
		t.Body.Resize(oldb.Bounds().Size())
	}

	// If it didn't render correctly; restore
	// the text buffer
	if !t.Body.Graphical() && t.Body != oldb {
		t.Body = oldb
		t.Config.Image = false
		t.Get(t.FileName())
	}
}

func unrender(t *tag.Tag) {
	var oldb = t.Body
	t.Config.Image = false
	t.Body = win.New(t.Dev, &t.Config.Body)
	t.Get(t.FileName())
	t.Body.Move(oldb.Bounds().Min)
	t.Body.Resize(oldb.Bounds().Size())
}
