package main

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/as/a/dump"
	"github.com/as/ui/col"
	"github.com/as/ui/tag"
)

var ErrBadDump = errors.New("bad dump file")

func Load(g *Grid, file string) (e error) {
	defer func() {
		dump.Printf("load: err: %v", e)
	}()

	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	r := dump.NewScanner(fd)

	var wd, f0, f1 string
	if !r.Scan(&wd, &f0, &f1) {
		return fmt.Errorf("load: wd: %v", r.Err())
	}
	if err = os.Chdir(wd); err != nil {
		dump.Printf("cd %q: %s", wd, err)
	}
	f0, f1 = f0, f1

	var colcents []int
	if !r.Scan(&colcents) {
		return fmt.Errorf("load: colcents: %v", r.Err())
	}
	dump.Printf("colcents: %#v\n", colcents)

	glabel := ""
	if !r.Scan('w', ' ', &glabel) {
		return fmt.Errorf("load: glabel: %v", r.Err())
	}
	dump.Printf("grid: label: %q\n", glabel)

	g.Tag.Label.Delete(0, g.Tag.Label.Len())
	g.Tag.Label.InsertString(glabel, 0)
	g.Tag.Label.Select(0, 0)

	var cols []*col.Col
	x := g.Bounds().Min.X
	dx := g.Bounds().Dx()
	for i, cent := range colcents {
		var (
			j      int
			tlabel string
		)
		if !r.Scan('c', &j, &tlabel) {
			return fmt.Errorf("load: tlabel: %v", r.Err())
		}
		if i != j {
			dump.Printf("load: col: mismatch col ids: %d != %d", i, j)
		}
		dump.Printf("load: col: %d: tlabel: %q", i, tlabel)
		c := col.New(g.Dev(), nil)
		c.Tag.Label.InsertString(tlabel, 0)
		c.Tag.Label.Select(0, 0)
		col.Attach(g, c, image.Pt(x+cent*dx/100, 0))
		cols = append(cols, c)
	}

Loop:
	for {
		var i, j, q0, q1, cent, nbody, wid, ntag, _, dir, dirty int
		kind := byte('?')
		if !r.Scan(&kind) {
			break
		}

		switch kind {
		case 'F':
			r.Scan(&i, &j, &q0, &q1, &cent, &nbody, '\n', &wid, &ntag, &nbody, &dir, &dirty)
		case 'f':
			r.Scan(&i, &j, &q0, &q1, &cent, '\n', &wid, &ntag, &nbody, &dir, &dirty)
			nbody = 0
		default:
			dump.Printf("load: win: c%dr%d: scan: kind %q: can't parse this type: ", i, j, kind)
			break Loop
		}

		if r.Err() != nil {
			dump.Printf("load: win: scan: %c: col(%d) row(%d): %v", kind, i, j, r.Err())
			break
		}

		t := tag.New(g.Dev(), nil)
		n, err := io.Copy(t.Label, io.LimitReader(r, int64(ntag)))
		dump.Printf("tag: copy label: %d bytes: err (%v)", n, err)
		t.Label.Select(0, 0)

		if !r.Scan('\n') {
			logf("tag label delim: %v", err)
		}

		n, err = io.Copy(t, io.LimitReader(r, int64(nbody)))
		dump.Printf("tag: copy body: %d bytes: err (%v)", n, err)

		t.Select(int64(q0), int64(q1))
		c := cols[i]
		col.Attach(c, t, image.Pt(0, c.Bounds().Min.Y+cent*c.Bounds().Dy()/100))
	}

	return nil
}

func Dump(g *Grid, wdir string, font0, font1 string) {
	d, _ := dump.Create("a.dump")
	defer d.Close()

	d.Line(wdir)
	d.Line(font0)
	d.Line(font1)

	// Each column's % offset from the start of the X-axis
	sep := byte(' ')
	for i, c := range g.List {
		if i+1 == len(g.List) {
			sep = '\n'
		}
		d.Int(flattenX(c, g), sep)
	}

	// The label for the grid itself
	fmt.Fprintf(d, "w %s\n", g.Tag.Label.Bytes())

	// The label for each column along with the column's number
	for i, c := range g.List {
		c, _ := c.(*col.Col)
		if c == nil {
			continue
		}
		fmt.Fprintf(d, "c%11d %s\n", i, c.Tag.Label.Bytes())
	}

	wid := 1

	// The windows
	for i, c := range g.List {
		c, _ := c.(*col.Col)
		if c == nil {
			continue
		}

		for j, t := range c.List {
			t, _ := t.(*tag.Tag)
			if c == nil {
				continue
			}
			// x, e, f, F

			q0, q1 := t.Dot()
			cent := flattenY(t, c)
			ntag := t.Label.Len()
			nbody := t.Len()
			dir := 0
			dirty := 1

			fmt.Fprintf(d, "F%11d %11d %11d %11d %11d %11d \n", i, j, q0, q1, cent, nbody)
			fmt.Fprintf(d, "%11d %11d %11d %11d %11d %s\n%s", wid, ntag, nbody, dir, dirty, t.Label.Bytes(), t.Bytes())
			wid++
		}
	}
}

func flattenY(w, c Plane) int {
	return 100 * (w.Bounds().Min.Y - c.Bounds().Min.Y) / c.Bounds().Dy()
}
func flattenX(c, g Plane) int {
	return 100 * (c.Bounds().Min.X - g.Bounds().Min.X) / g.Bounds().Dx()
}
