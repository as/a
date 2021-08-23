package edit

import (
	"bytes"
	"regexp"

	"github.com/as/io/rev"
	"github.com/as/text/find"
)

type Editor interface {
	Insert(p []byte, at int64) (n int)
	Delete(q0, q1 int64) (n int)
	Select(q0, q1 int64)
	Dot() (q0, q1 int64)
	Len() int64
	Bytes() []byte
	Close() error
}

// Address implements Set on the Editor. Possibly selecting
// some range of text (a dot).
type Address interface {
	// Set computes and sets the address on the provided Editor
	Set(f Editor)
	// Back returns true if the address semantics should be executed in reverse
	Back() bool
}

// Regexp is an address computed by a regexp
type Regexp struct {
	re   *regexp.Regexp
	back bool
	rel  int
}

// Bytes is an address computed by a relative or absolute byte offset
type Byte struct {
	Q   int64
	rel int
}

// Line is an address computed by a relative or absolute byte offset
type Line struct {
	Q   int64
	rel int
}

// Dot is the current dot address
type Dot struct {
	Q0, Q1 int64
}

type File struct {
	Name string
	Dot
}

// Compound combines two address values with an operator
type Compound struct {
	a0, a1 Address
	op     byte
}

func (r Regexp) Back() bool   { return r.rel == -1 }
func (b Byte) Back() bool     { return b.rel == -1 }
func (l Line) Back() bool     { return l.rel == -1 }
func (d Dot) Back() bool      { return false }
func (c Compound) Back() bool { return c.a1.Back() }

func (c *Compound) Set(f Editor) {
	if c.a0 == nil {
		return
	}
	c.a0.Set(f)
	q0, _ := f.Dot()

	if c.a1 == nil {
		return
	}
	c.a1.Set(f)
	_, r1 := f.Dot()
	if c.Back() {
		return
	}
	f.Select(q0, r1)
}

func (b *Byte) Set(f Editor) {
	q0, q1 := f.Dot()
	q := b.Q
	if b.rel == -1 {
		f.Select(q+q0, q+q0)
	} else if b.rel == 1 {
		f.Select(q+q1, q+q1)
	} else {
		f.Select(q, q)
	}
}
func (r *Regexp) Set(f Editor) {
	_, q1 := f.Dot()
	org := q1
	buf := bytes.NewReader(f.Bytes()[q1:])
	loc := r.re.FindReaderIndex(buf)
	if loc == nil {
		return
	}
	r0, r1 := int64(loc[0])+org, int64(loc[1])+org
	if r.rel == 1 {
		//r0 = r1
	}
	f.Select(r0, r1)
}

func (r *Line) Set(f Editor) {
	p := f.Bytes()
	switch r.rel {
	case 0:
		q0, q1 := find.Findline2(r.Q, bytes.NewReader(p))
		f.Select(q0, q1)
	case 1:
		_, org := f.Dot()
		r.Q++
		if org == 0 || p[org-1] == '\n' {
			r.Q--
		}
		p = p[org:]
		q0, q1 := find.Findline2(r.Q, bytes.NewReader(p))
		f.Select(q0+org, q1+org)
	case -1:
		org, _ := f.Dot()
		r.Q = -r.Q + 1
		if org == 0 || p[org-1] == '\n' {
			//r.Q--
		}
		p = p[:org]
		q0, q1 := find.Findline2(r.Q, rev.NewReader(p)) // 0 = len(p)-1
		//fmt.Printf("Line.Set 1: %d:%d\n", q0, q1)
		l := q1 - q0
		q0 = org - q1
		q1 = q0 + l
		q0 = q1 - l
		if q0 >= 0 && q0 < int64(len(f.Bytes())) && f.Bytes()[q0] == '\n' {
			q0++
		}
		//fmt.Printf("Line.Set 2: %d:%d\n", q0, q1)
		f.Select(q0, q1)
	}
}

func (d Dot) Set(f Editor) {
	//f.Select(d.Q0, d.Q1)
}
