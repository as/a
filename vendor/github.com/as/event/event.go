package event

import (
	"fmt"
	"unicode"
)

type Editor interface {
	Insert(p []byte, at int64) (n int)
	Delete(q0, q1 int64) (n int)
	Len() int64
	Bytes() []byte
	Select(q0, q1 int64)
	Dot() (q0, q1 int64)
	Write(p []byte) (n int, err error)
	Close() error
}

type Record interface {
	String() string
	//	Equal(Record) bool
	Coalesce(Record) Record
	//	Invert(Record) Record
}

type Event interface {
	isevent()
}

type Rec struct {
	ID     int
	Kind   byte
	Q0, Q1 int64
	N      int64
	P      []byte
}

func (r Rec) String() string {
	return fmt.Sprintf("%d	%x	%d	%d	%d	%q\n", r.ID, r.Kind, r.Q0, r.Q1, r.N, r.P)
}
func (r Rec) Record() (p []byte, err error) {
	return []byte(r.String()), nil
}

func (r Insert) Record() (p []byte, err error) {
	r.Kind = 'i'
	r.N = int64(len(p))
	return r.Rec.Record()
}
func (r Delete) Record() (p []byte, err error) {
	r.Kind = 'd'
	return r.Rec.Record()
}
func (r Select) Record() (p []byte, err error) {
	r.Kind = 's'
	return r.Rec.Record()
}

type Insert struct {
	Rec
}

func space(b byte) int {
	if unicode.IsSpace(rune(b)) {
		return 1
	}
	return 0
}

func (e *Insert) Coalesce(v Record) Record {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case *Insert:
		if len(v.P) == 0 {
			return v
		}
		if v.ID != e.ID {
			return nil
		}
		if space(v.P[0]) != space(e.P[0]) {
			return nil
		}
		if v.Q0 == e.Q0 {
			e.Q1 = v.Q1
			e.P = append(e.P, v.P...)
			return e
		}
	case *Delete:
		if v.ID != e.ID {
			return nil
		}
		if e.Q0 == v.Q0 && e.Q1 == v.Q1 {
			// EXPERIMENT
			e.Kind = 'w'
			//e.P = v.P
			return &Write{Rec: e.Rec}
		}
		if len(v.P) == 1 && len(e.P) == 1 && (space(v.P[0]) != space(e.P[0])) {
			return nil
		}
		if e.Q0 >= v.Q0 && v.Q1 == e.Q1-1 {
			// 0      3        3       4
			e.Q1 -= v.Q1 - v.Q0
			e.P = e.P[:min(int64(len(e.P)), e.Q1)]
			return e
		}
	case *Select:
		return nil
	}
	return nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func (e *Delete) Coalesce(v Record) Record {
	switch v := v.(type) {
	case *Insert:
		if e.Q0 == v.Q0 && e.Q1 == v.Q1 {
			e.Kind = 'w'
			e.P = v.P
			return &Write{Rec: e.Rec}
		} else if e.Q0 == v.Q0 && e.Q1 < v.Q1 {
			e.Kind = 'w'
			split := int64(len(v.P)) + (e.Q1 - v.Q1)
			e.P = v.P[:split]
			return &Write{Rec: e.Rec,
				Residue: &Insert{
					Rec: Rec{
						Q0:   v.Q0 + (e.Q1 - e.Q0),
						Q1:   v.Q1,
						Kind: 'i',
						P:    v.P[split:],
					},
				},
			}
		} else {
			d0, d1 := e.Q1, v.Q1
			e.Kind = 'w'
			e.P = v.P
			e.Q1 = v.Q1
			return &Write{
				Rec: e.Rec,
				Residue: &Delete{
					Rec: Rec{
						Kind: 'd',
						Q0:   d1,
						Q1:   d0,
					},
				},
			}
		}
	case *Delete:
		if v.ID != e.ID {
			return nil
		}
		if e.Q1 != v.Q0 {
			return nil
		}
		e.Q1 = v.Q1
		e.P = append(v.P, e.P...)
		return e
	}
	return nil
}
func (e *Select) Coalesce(v Record) Record {
	switch v := v.(type) {
	case *Select:
		if v.Q0 == e.Q0 && v.Q1 == e.Q1 {
			return e
		}
	}
	return nil
}

type Delete struct {
	Rec
}
type Select struct {
	Rec
}
type SetOrigin struct {
	ID    int
	Q0    int64
	Exact bool
}
type Fill struct {
}
type Scroll struct {
}
type Redraw struct {
}
type Sweep struct {
}
type Move struct {
}
type Cmd struct {
	Rec
	From    Editor
	To      []Editor
	Basedir string
	Name    string
}
type Look struct {
	Rec
	From    Editor
	To      []Editor
	Basedir string
	Name    string
}

type Get struct {
	ID    int
	Name  string
	Path  string
	Addr  string
	IsDir bool
}
type Put struct {
	ID int
}

func (Look) isevent()      {}
func (Cmd) isevent()       {}
func (Insert) isevent()    {}
func (Write) isevent()     {}
func (Delete) isevent()    {}
func (Select) isevent()    {}
func (SetOrigin) isevent() {}
func (Fill) isevent()      {}
func (Scroll) isevent()    {}
func (Redraw) isevent()    {}
func (Sweep) isevent()     {}
func (Move) isevent()      {}
func (Get) isevent()       {}
func (Put) isevent()       {}
