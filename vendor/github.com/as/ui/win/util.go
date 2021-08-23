package win

import (
	"fmt"
)

var Db = new(Dbg)

type Dbg struct {
	indent int
}

func Trace(p *Dbg, msg string) *Dbg {
	p.Trace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func Un(p *Dbg) {
	p.indent--
	p.Trace(")")
}
func (p *Dbg) Trace(a ...interface{}) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

func (w *Win) InsertString(s string, q0 int64) int {
	return w.Insert([]byte(s), q0)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func clamp(v, l, h int64) int64 {
	if v < l {
		return l
	}
	if v > h {
		return h
	}
	return v
}
