package col

import (
	"bytes"
	"fmt"

	"github.com/as/ui/win"
)

var noteDelim = []byte("	|")

func notesize(k interface{}) {
	Note(k.(Labeled), "%s", k.(Plane).Bounds())
}

func Note(l Labeled, fm string, v ...interface{}) {
	if l.Label() == nil {
		return
	}
	w := l.Label()
	q0 := int64(bytes.Index(w.Bytes(), noteDelim))
	if q0 == -1 {
		return
	}
	q0 += 2
	q1 := int64(len(w.Bytes()))
	if q0 < q1 {
		w.Delete(q0, q1)
	}

	w.InsertString(fmt.Sprintf(fm, v...), q1)
}

type Labeled interface {
	Label() *win.Win
}

func (co *Col) badID(id int) bool {
	return id < 0
}

func (col *Col) PrintList() {
	for i, v := range col.List {
		fmt.Printf("%d: %#v\n", i, v)
	}
}
