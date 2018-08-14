package main

import "github.com/as/ui/win"

type Track struct {
	win    *win.Win
	q0, q1 int64
}

var track = Track{}

func (p *Track) esc() {
	if act != track.win {
		return
	}
	if q0, q1 := track.win.Dot(); q0 != q1 {
		track.win.Delete(q0, q1)
	} else {
		q0 = track.q1
		track.win.Select(q0, q1)
	}
}

func (p *Track) set(force bool) {
	if act != track.win || force {
		if act, ok := act.(*win.Win); ok {
			track.q0, track.q1 = act.Dot()
			track.win = act
		}
	}
}
