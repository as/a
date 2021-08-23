package text

import (
	"github.com/as/event"
	"github.com/as/worm"
)

func NewRecorder(l worm.Logger) *Recorder {
	return &Recorder{
		l: l,
	}
}

type Recorder struct {
	l worm.Logger
}

func (r *Recorder) Insert(p []byte, q0 int64) int {
	if len(p) == 0 {
		//return 0
	}
	ev := &event.Insert{event.Rec{Kind: 'i', P: p, Q0: q0, Q1: q0 + int64(len(p))}}
	r.l.Write(ev)
	return len(p)
}
func (r *Recorder) Delete(q0, q1 int64) int {
	ev := &event.Delete{event.Rec{Kind: 'd', Q0: q0, Q1: q1}}

	r.l.Write(ev)
	return int(q1 - q0)
}
func (r *Recorder) Write(p []byte) (int, error) {
	return r.Insert(p, (^int64(0))>>1), nil
}
func (r *Recorder) Select(q0, q1 int64) {
	r.l.Write(&event.Select{event.Rec{Kind: 's', Q0: q0, Q1: q1}})
}
