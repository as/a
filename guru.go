package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/as/ui/tag"
)

var guruModes = "callees callers callstack definition describe freevars implements peers pointsto referrers what whicherrs"

func addrfmt(label string, q0, q1 int64) string {
	if q0 == q1 {
		return fmt.Sprintf("%s:#%d", label, q1)
	}
	return fmt.Sprintf("%s:#%d,#%d", label, q0, q1)
}

func (g *Grid) afindguru(wd string, name string) *tag.Tag {
	if !strings.HasSuffix(name, "+Guru") {
		name += "+Guru"
	}
	t := g.FindName(name)
	if t == nil {
		c := g.List[len(g.List)-1].(*Col)
		t = New(c, "", name, SizeThirdOf).(*tag.Tag)
		if t == nil {
			panic("cant create tag")
		}
	}
	return t
}

func (g *Grid) aguru(fm string, i ...interface{}) {
	t := g.afindguru(".", "")
	q1 := t.Body.Len()
	t.Body.Select(q1, q1)
	n := int64(t.Body.Insert([]byte(time.Now().Format(timefmt)+": "+fmt.Sprintf(fm, i...)), q1))
	t.Body.Select(q1+n, q1+n)
	ajump(t.Body, cursorNop)
}

func (g *Grid) Selection() string {
	return string(g.List[0].(*tag.Tag).Win.Rdsel())
}

func (g *Grid) Label() string {
	return string(g.List[0].(*tag.Tag).Win.Bytes())
}

func (g *Grid) guru(label string, q0, q1 int64) error {
	if !strings.HasSuffix(label, ".go") {
		return nil
	}
	asel := g.Selection()
	mode := ""
	//	scope := "."
	for _, v := range strings.Fields(guruModes) {
		if asel == v {
			mode = v
			break
		}
	}
	if mode == "" {
		return nil
	}

	data, err := exec.Command(
		"guru",
		//		"-scope",
		//		scope,
		mode,
		addrfmt(label, q0, q1),
	).CombinedOutput()

	if err != nil {
		g.aerr("guru: %s", err)
	}
	if len(data) != 0 {
		g.aguru("%s", data)
	}
	return err
}
