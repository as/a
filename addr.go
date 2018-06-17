package main

import (
	"github.com/as/text"
	"github.com/as/text/find"
)

type q int64

func (a q) In(a1 Addr) bool { return text.Region3(int64(a), int64(a1.s-1), int64(a1.e)) == 0 }

type Addr struct {
	s q
	e q
}

type a = Addr

func d2a(q0, q1 int64) Addr { return Addr{q(q0), q(q1)} }

func (a Addr) Empty() bool { return a.s == a.e }
func (a Addr) Len() int    { return int(a.e - a.s) }
func (a Addr) In(a1 Addr) bool {
	if a.Empty() {
		return a.e.In(a1)
	}
	return a.Len() <= a1.Len() && text.Region5(int64(a.s), int64(a.e), int64(a1.s), int64(a1.e)) == 0
}
func (a Addr) Dot() (q0, q1 int64) { return int64(a.s), int64(a.e) }

func expandAddr(a Addr, ed text.Editor) Addr {
	if !a.Empty() {
		return a
	}
	a0 := d2a(ed.Dot())
	if a.In(a0) {
		return a0
	}
	return a
}

func expandFile(a Addr, ed text.Editor) Addr {
	if !a.Empty() {
		return a
	}
	return d2a(find.ExpandFile(ed.Bytes(), int64(a.s)))
}
