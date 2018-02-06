package main

import (
	"image"
	"time"

	mus "github.com/as/text/mouse"
	"github.com/as/ui/win"
	"golang.org/x/mobile/event/paint"
)

var (
	// Records the last time a scroll event was recieved so
	// a proper acceleration factor can be applied.
	lastscroll   time.Time
	lastscrolldy int
	scrollmul    = 1
	smoothness   float64
	hurryup      float64
)

func smoothscroll(act *win.Win, e mus.ScrollEvent) {
}

func zsmoothscroll(act *win.Win, e mus.ScrollEvent) {
	wind := act.Window()
	dp := act.Loc().Min
	sr := act.Buffer().Bounds()
	sr.Min.X += 10
	sr.Min.Y += 15
	dp.X += 10
	dp.Y += 15
	h := act.Frame.Font.Dy()
	if false && e.Button == -1 {
		for i := 0; i != -h*e.Dy; i -= e.Dy {
			wind.Upload(dp, act.Buffer(), sr.Add(image.Pt(0, i)))
			wind.Publish()
			act.Flush()
		}
		e.Dy = -e.Dy
	} else if false {
		for i := 0; i != h*e.Dy; i += e.Dy {
			wind.Upload(dp, act.Buffer(), sr.Add(image.Pt(0, i)))
			wind.Publish()
			act.Flush()
		}
	}
	if e.Button == -1{
		e.Dy = -e.Dy
	}
	actTag.Body.Scroll(e.Dy)
	//act.Refresh()
}

func doScrollEvent(act *win.Win, e mus.ScrollEvent) {
	if e.Button == -1{
		e.Dy = -e.Dy
	}
	actTag.Body.Scroll(e.Dy)
	act.Window().Send(paint.Event{})
}
func zdoScrollEvent(act *win.Win, e mus.ScrollEvent) {
	tm := time.Now()
	wind := act.Window()

	if e.Dy == 0 {
		lastscroll = tm
		lastscrolldy = 0
		return
	}
	if tm.Sub(lastscroll) > time.Second/3 {
		smoothness = 1.0
		hurryup = 0.33
	}
	smoothness += float64(tm.Sub(lastscroll) / time.Second)
	if smoothness > 4.0 {
		hurryup = 0.33
		smoothness = 4.0
	}
	if smoothness > 0.0 {
		smoothness -= float64(e.Dy) * 0.20
		if e.Dy != 3 {
			e.Dy = 3
		}
		smoothscroll(act, e)
		wind.SendFirst(mus.Drain{})
		wind.Send(mus.DrainStop{})
		ck()
	} else {
		hurryup *= 1.015
		if e.Button == -1 {
			e.Dy = -e.Dy
		}
		actTag.Body.Scroll(int(float64(e.Dy) * hurryup))
		ck()
	}
	lastscroll = tm
	lastscrolldy = e.Dy
}
