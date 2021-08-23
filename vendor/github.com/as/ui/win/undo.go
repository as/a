package win

import (
	"fmt"
	"time"
)

// EnableUndoExperiment toggles the state of the Undo/Redo
// features. This is set to off by default. Eventually the variable
// will be removed.
//
// See-also:
// 	ins.go:/EnableUndoExperiment/
// 	del.go:/EnableUndoExperiment/
//
var EnableUndoExperiment = false

func (w *Win) Undo() bool {
	if !EnableUndoExperiment {
		return false
	}
	return w.ops.Undo(w)
}

func (w *Win) Redo() bool {
	if !EnableUndoExperiment {
		return false
	}
	return w.ops.Redo(w)
}

type buffer struct {
	tm      time.Time
	scratch op
	Ops
}

func (o *buffer) commit(ins bool) {
	if ins {
		o.Ops.Insert(o.scratch.p, o.scratch.q0)
		o.scratch.q0 += o.scratch.q0 + int64(len(o.scratch.p))
	} else {
		o.Ops.Delete(o.scratch.q0, o.scratch.q1, o.scratch.p)
		o.scratch.q1 = o.scratch.q0
	}
}
func (o *buffer) chop(p []byte, q0 int64) bool {
	return o.timeout() || q0 != o.scratch.q0+1 || string(p) == " " || string(p) == "	"
}
func (o *buffer) Insert(p []byte, q0 int64) int {
	if o.chop(p, q0) {
		o.Ops.Insert(o.scratch.p, o.scratch.q0)
		o.scratch.q0 = q0
		o.scratch.p = []byte{}
	}
	o.scratch.q1 = o.scratch.q0 + int64(len(p))
	o.scratch.p = append(o.scratch.p, p...)
	return (len(p))
}
func (o *buffer) Delete(q0, q1 int64, p []byte) int {
	return o.Ops.Delete(q0, q1, p)
}
func (o *buffer) timeout() bool {
	to := time.Since(o.tm) > time.Second*3
	o.tm = time.Now()
	return to
}

type Ops struct {
	q  int
	Op []Op
}

func (o *Ops) Insert(p []byte, q0 int64) int {
	println(fmt.Sprintf("#%d i,%q,\n", q0, p))
	return o.commit(OpIns{q0: q0, q1: int64(len(p)) + q0, p: p})
}
func (o *Ops) Delete(q0, q1 int64, p []byte) int {
	println(fmt.Sprintf("#%d,#%d d\n", q0, q1))
	return o.commit(OpDel{q0: q0, q1: q1, p: []byte(string(p[q0:q1]))})
}
func (o *Ops) Redo(w *Win) bool {
	if o.q == len(o.Op) {
		return false
	}
	o.Op[o.q].Do(w)
	o.q++
	return true
}
func (o *Ops) Undo(w *Win) bool {
	if o.q == 0 {
		return false
	}
	o.q--
	o.Op[o.q].Un().Do(w)
	return true
}
func (o *Ops) commit(op Op) int {
	if o.q != len(o.Op) {
		o.Op = append([]Op{}, o.Op[:o.q]...)
	}
	o.Op = append(o.Op, op)
	o.q++
	return 0
}

type (
	Op interface {
		Do(w *Win) int
		Un() Op
	}
	OpIns op
	OpDel op
	op    struct {
		q0, q1 int64
		p      []byte
	}
)

func (o OpIns) Do(w *Win) int { return w.insert(o.p, o.q0) }
func (o OpDel) Do(w *Win) int { return w.delete(o.q0, o.q1) }
func (o OpIns) Un() Op        { return OpDel(o) }
func (o OpDel) Un() Op        { return OpIns(o) }
