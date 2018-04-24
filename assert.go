package main

import (
	"image"
	"log"

	"github.com/as/ui/tag"
)

func assert(where string, g *Grid) {
	log.Printf("%s: grid %s\n", where, g.Loc())
	for i, c := range g.List {
		for j, c2 := range g.List {
			if c == c2 {
				continue
			}
			if c.Loc().Intersect(c2.Loc()) != image.ZR {
				log.Printf("%s: 	col %v %s intersects col %v %s\n", where, i, c.Loc(), j, c2.Loc())
				panic("suicide")
			}
		}
	}
	if g.Tag.Vis != tag.VisTag {
		log.Printf("%v\n", g.Tag.Vis)
		panic("grid tag not vistag") // Put
	}
	for i, c := range g.List {
		if c.(*Col).Tag.Vis != tag.VisTag {
			log.Printf("number %d = %v\n", i, g.Tag.Vis)
			//panic("col tag not vistag")
		}
		if !c.Loc().In(g.Loc()) {
			log.Printf("%s: col %v %s not in grid %s\n", where, i, c.Loc(), g.Loc())
			//			panic("suicide")
		}
		log.Printf("%s: col %v %s is good\n", where, i, c.Loc())
		{
			c, _ := c.(*Col)
			if c == nil {
				continue
			}
			for j, t := range c.List {
				if !t.Loc().In(c.Loc()) {
					log.Printf("%s: 	tag %v %s not in col %v %s\n", where, j, t.Loc(), i, c.Loc())
					//panic("suicide")
				}
				if !t.Loc().In(g.Loc()) {
					log.Printf("%s: 	tag %v %s not in grid %v %s\n", where, j, t.Loc(), i, g.Loc())
					//panic("suicide")

				}
				log.Printf("%s: 	tag %v %s is good\n", where, j, t.Loc())
			}
		}
	}
}
