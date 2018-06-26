package main

import (
	"image"
	"log"

	"github.com/as/ui/tag"
)

const debug = 1 ^ 1

func assert(where string, g *Grid) {
	if debug == 0 {
		return
	}
	log.Printf("%s: grid %s\n", where, g.Bounds())
	for i, c := range g.List {
		for j, c2 := range g.List {
			if c == c2 {
				continue
			}
			if c.Bounds().Intersect(c2.Bounds()) != image.ZR {
				log.Printf("%s: 	col %v %s intersects col %v %s\n", where, i, c.Bounds(), j, c2.Bounds())
				panic("suicide")
			}
		}
	}
	if g.Tag.Vis != tag.VisTag {
		log.Printf("%v\n", g.Tag.Vis)
		//		panic("grid tag not vistag") // Put
	}
	for i, c := range g.List {
		if c.(*Col).Tag.Vis != tag.VisTag {
			log.Printf("number %d = %v\n", i, g.Tag.Vis)
			//panic("col tag not vistag")
		}
		if !c.Bounds().In(g.Bounds()) {
			log.Printf("%s: col %v %s not in grid %s\n", where, i, c.Bounds(), g.Bounds())
			//			panic("suicide")
		}
		log.Printf("%s: col %v %s is good\n", where, i, c.Bounds())
		{
			c, _ := c.(*Col)
			if c == nil {
				continue
			}
			for j, t := range c.List {
				if !t.Bounds().In(c.Bounds()) {
					log.Printf("%s: 	tag %v %s not in col %v %s\n", where, j, t.Bounds(), i, c.Bounds())
					//panic("suicide")
				}
				if !t.Bounds().In(g.Bounds()) {
					log.Printf("%s: 	tag %v %s not in grid %v %s\n", where, j, t.Bounds(), i, g.Bounds())
					//panic("suicide")

				}
				log.Printf("%s: 	tag %v %s is good\n", where, j, t.Bounds())
			}
		}
	}
}
