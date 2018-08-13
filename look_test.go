package main

import (
	"fmt"
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
	LLabel                   // The user looked from the label
)

func (f LFlag) On(bit LFlag) bool {
	return f&LLabel != 0
}

type ent struct {
	pre   Addr
	sweep Addr
	post  Addr
	flags LFlag
	data  string
}

// Window Resolver
// Filesystem Resolver
// Look Resolver

func TestLook(t *testing.T) {
	etch := ui.NewEtch()

	sentence := "the quick brown fox jumps over the lazy dog"
	norm := LScav | LMatch | LJump

	//exampleURL :=  "http://example.com"
	type a = Addr
	for name, tc := range map[string]struct {
		label ent
		ent
		onlabel bool
	}{
		// Null selection in body, select letters of "the"
		"the0/0": {ent: ent{a{0, 0}, a{}, a{31, 34}, norm, sentence}},
		"the0/1": {ent: ent{a{0, 0}, a{0, 1}, a{31, 32}, norm, sentence}},
		"the0/2": {ent: ent{a{0, 0}, a{0, 2}, a{31, 33}, norm, sentence}},

		// Partial selection of "th"
		"the2/0": {ent: ent{a{0, 2}, a{}, a{31, 33}, norm, sentence}},
		"the2/1": {ent: ent{a{0, 2}, a{0, 1}, a{31, 32}, norm, sentence}},
		"the2/2": {ent: ent{a{0, 2}, a{0, 2}, a{31, 33}, norm, sentence}},

		//
		"h/0": {ent: ent{a{1, 2}, a{1, 1}, a{32, 33}, norm, sentence}},

		// Wrap-around test
		"the1": {ent: ent{a{31, 33}, a{5, 5}, a{31, 33}, norm | LWrap, sentence}},

		// Clicked on the tag label instead
		"dog":  {ent: ent{a{0, 0}, a{}, a{40, 43}, norm | LLabel, sentence}},
		"dog1": {ent: ent{a{0, 0}, a{0, 1}, a{40, 41}, norm | LLabel, sentence}},
		"dog2": {ent: ent{a{0, 0}, a{0, 2}, a{40, 42}, norm | LLabel, sentence}},

		// If it's clicking in the label
		"label/the": {
			label:   ent{a{0, 0}, a{0, 3}, a{0, 3}, 0, "the"},
			ent:     ent{a{0, 0}, a{0, 0}, a{0, 3}, norm | LLabel, sentence},
			onlabel: true,
		},
	} {
		t.Run(name, func(t *testing.T) {

			// TODO(as): wow, this is a lot of code to initialize one grid
			g := NewGrid(etch, GridConfig)
			c := col.New(etch, ColConfig)

			if tc.label.data == "" {
				tc.label.data = "dog"
			}

			tag := tag.New(etch, nil)
			tag.Label.Insert([]byte(tc.label.data), 0)

			w := tag.Window
			w.Insert([]byte(tc.data), 0)
			w.Select(int64(tc.pre.s), int64(tc.pre.e))

			col.Attach(g, c, image.ZP)
			col.Attach(c, tag, image.ZP)

			from := w
			if tc.flags.On(LLabel) {
				from = tag.Label
			} else if name == "dog" {
				t.Fatal("bad test")
			}

			a1 := tc.sweep
			if a1.Empty() {
				a1 = expandAddr(a1, from)
			}
			if a1.Empty() {
				a1 = expandFile(a1, from)
			}

			// TODO(as): wow, this is a lot of code to express one look
			ev := event.Look{
				Name: "w", From: from,
				To: []event.Editor{w},
				Rec: event.Rec{
					Q0: int64(a1.s), Q1: int64(a1.e),
					P: from.Bytes(),
				},
			}
			t.Logf("looking for %s", from.Bytes()[int64(a1.s):int64(a1.e)])
			g.Look(ev)

			q0, q1 := w.Dot()
			have := a{q(q0), q(q1)}
			if have != tc.post {
				t.Fatalf("have %v, want %v", have, tc.post)
			}
		})
	}

}

func TestExpand(t *testing.T) {
	etch := ui.NewEtch()

	sentence := "stall install reinstall installing stalling"
	g := NewGrid(etch, GridConfig)
	c := col.New(etch, ColConfig)
	tag := tag.New(etch, nil)
	w := tag.Window
	w.Insert([]byte(sentence), 0)
	w.Select(0, 0)
	col.Attach(g, c, image.ZP)
	col.Attach(c, tag, image.ZP)

	a1 := a{0, q(len("stall"))}

	for i := 0; i < 10; i++ {
		if a1.Empty() {
			a1 = expandAddr(a1, w)
		}
		if a1.Empty() {
			a1 = expandFile(a1, w)
		}
		ev := event.Look{
			Name: "w", From: w,
			To: []event.Editor{w},
			Rec: event.Rec{
				Q0: int64(a1.s), Q1: int64(a1.e),
				P: w.Bytes(),
			},
		}
		g.Look(ev)
		q0, q1 := w.Dot()
		if have := string(w.Bytes()[q0:q1]); have != "stall" {
			t.Fatalf("have %v, want %v", have, "stall")
		}
	}
}

func TestLookAddrFileLine(t0 *testing.T) {
	t0.Skip("fix this bug")
	etch := ui.NewEtch()

	label := "TestLookAddrFileLine"

	g := NewGrid(etch, GridConfig)
	c := col.New(etch, ColConfig)
	t := tag.New(etch, nil)
	w := t.Window

	fmt.Fprint(t.Label, label+"\t:2")
	fmt.Fprint(w, label+":2\ntwo\nthree")

	if t := g.FindName(label); t != nil {
		t0.Fatal("FindName: cant find window")
	}

	col.Attach(g, c, image.ZP)
	col.Attach(c, t, image.ZP)

	ev := event.Look{
		Name: "w", From: w,
		To: []event.Editor{w},
		Rec: event.Rec{
			Q0: 0, Q1: int64(len(label)),
			P: w.Bytes(),
		},
	}
	g.Look(ev)
	q0, q1 := w.Dot()

	have := string(w.Bytes()[q0:q1])
	want := "two\n"
	if have != want {
		t0.Fatalf("have %q, want %q", have, want)
	}
}
