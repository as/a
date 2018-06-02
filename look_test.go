package main

import (
	"image"
	"testing"

	"github.com/as/event"
	"github.com/as/ui"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
)

type LFlag int

const (
	LRes   LFlag = 1 << iota // Restore original selection
	LScav                    // Expand the nil selection to gather context
	LMatch                   // Expect the look to find a result
	LJump                    // Expect the look to jump to a window
	LNew                     // Expect the look to open a new window
	LWrap                    // Expect the look to wrap around
)

var testText = "the quick brown fox jumps over the lazy dog"

func TestLook(t *testing.T) {
	etch := ui.NewEtch()
	type a struct {
		s, e int64
	}
	for name, tc := range map[string]struct {
		pre   a // pre-sweep selection
		sweep a // selection after sweep complete
		post  a // selection after look
		flags LFlag

		data string
	}{
		"the0": {a{0, 2}, a{}, a{31, 33}, LScav | LMatch | LJump, testText},
		"the1": {a{31, 33}, a{}, a{0, 2}, LScav | LMatch | LJump | LWrap, testText},
		//	"the3": {a{0, 0}, a{0, 2}, a{31, 33}, LScav | LMatch | LJump, testText},
		// TODO(as): break it with selected look above
	} {
		t.Run(name, func(t *testing.T) {

			// TODO(as): wow, this is a lot of code to initialize one grid
			g := NewGrid(etch, GridConfig)
			c := col.New(etch, ColConfig)
			tag := tag.New(etch, nil)
			tag.Body.Insert([]byte(tc.data), 0)
			w := tag.Body
			w.Select(tc.pre.s, tc.pre.e)
			col.Attach(g, c, image.ZP)
			col.Attach(c, tag, image.ZP)

			// TODO(as): wow, this is a lot of code to express one look
			ev := event.Look{
				Name: "w", From: w,
				To: []event.Editor{w},
				Rec: event.Rec{
					Q0: tc.pre.s, Q1: tc.pre.e,
					P: w.Bytes()[tc.pre.s:tc.pre.e],
				},
			}

			g.Look(ev)

			q0, q1 := w.Dot()
			have := a{q0, q1}
			if have != tc.post {
				t.Fatalf("have %v, want %v", have, tc.post)
			}
		})

	}

}
