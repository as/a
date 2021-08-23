package tag

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/as/text/find"
)

func (t *Tag) FileName() string {
	if t == nil || t.Label == nil {
		return ""
	}
	name, err := bufio.NewReader(bytes.NewReader(t.Label.Bytes())).ReadString('\t')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(name)
}

func (t *Tag) Open(basepath, title string) {
	println(title)
	t.Get(title)
}

func (t *Tag) fixtag(abs string) {
	l := t.Label
	maint := find.Find(l.Bytes(), 0, []byte{'|'})
	if maint == -1 {
		maint = l.Len()
	}
	l.Delete(0, maint+1)
	l.InsertString(abs+"\tPut Del |", 0)
	l.Refresh()
}

type GetEvent struct {
	Basedir string
	Name    string
	Addr    string
	IsDir   bool
}
