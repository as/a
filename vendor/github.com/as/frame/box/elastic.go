package box

func (f *Run) Stretch2(nb int, bx []Box) (pb int) {
	panic("unimplemented")
}

// Elastic tabstop experiment section. Not ready for general use by any means
// the text/tabwriter package implements elastic tabstops, but that package
// assumes that all chars are the same width and that text needs to be rescanned.
//
// The frame already distinguishes between tabs, newlines, and plain text characters
// by encapsulating them in measured boxes. A direct copy of the tabwriter code would
// ignore the datastructures in the frame and their sublinear runtime cost.
//
// The current elastic algorithm has suboptimal runtime performance
//
func (f *Run) Stretch(nb int) (pb int) {
	if nb <= 0 {
		return 0
	}
	//	fmt.Println()
	//	fmt.Printf("\n\ntrace bn=%d\n", nb)
	nc := 0
	nl := 0
	dx := 0

	cmax := make(map[int]int)
	cbox := make(map[int][]int)

	nb = f.StartCell(nb)
	pb = nb - 1
	if nb == f.Nbox {
		return 0
	}
Loop:
	for ; nb < f.Nbox; nb++ {
		b := &f.Box[nb]
		//		fmt.Printf("switch box: %#v\n", b)
		switch b.Break() {
		case '\t':
			dx += b.Width
			cbox[nc] = append(cbox[nc], nb)
			max := cmax[nc]
			if dx > max {
				cmax[nc] = dx
			}
			nc++
			//			fmt.Printf("	tab: dx=%d ncol=%d\n", dx, nc)
			dx = 0
		case '\n':
			nl++
			dx = 0
			if nc == 0 {
				// A line with no tabs; end of cell
				//fmt.Printf("	nl (no cols): dx=%d nl=%d\n", dx, nl-1)
				break Loop
			}
			//fmt.Printf("	nl : dx=%d nl=%d nc=%d\n", dx, nl-1, nc)
			nc = 0
		default:
			dx += b.Width
			//fmt.Printf("	plain : dx=%d wid=%d nc=%d\n", dx, b.Width, nc)
		}
	}
	for c, bns := range cbox {
		max := cmax[c]
		for _, bn := range bns {
			b := &f.Box[bn]
			b.Width = max
			if bn == 0 {

			} else {
				pb := f.Box[bn-1]
				if pb.Break() != '\n' && pb.Break() != '\t' {
					b.Width -= f.Box[bn-1].Width
				}
			}
			if b.Width < b.Minwidth {
				b.Width = b.Minwidth
			}
		}
	}
	return pb
}

func (f *Run) Findcol(bn int, coln int) (cbn int, xmax int) {
	c := 0
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		if b.Break() == '\t' {
			c++
		}
		if b.Break() != '\n' {
			xmax += b.Width
		}
		if c == coln {
			break
		}
		bn++
	}
	if c != coln {
		return -1, 0
	}
	return bn, xmax

}

func (f *Run) Colof(bn int) (coln, xmax int) {
	if bn == 0 {
		return 0, 0
	}
	bs := f.StartLine(bn)
	for {
		b := &f.Box[bs]
		if b.Break() == '\t' {
			coln++
		}
		if b.Break() != '\n' {
			xmax += b.Width
		}
		if bn == bs {
			break
		}
		bs++
	}
	if xmax != 0 {
		coln++
	}
	return coln, xmax
}

// EndCell returns the first box beyond the end of the
// current cell under bn
func (f *Run) EndCell(bn int) int {
	oldbn := bn
	bn = f.StartLine(bn)
	ltb := 0
	ncol := 0
Loop:
	for ; bn != f.Nbox; bn++ {
		b := &f.Box[bn]
		switch b.Break() {
		case '\n':
			if ncol == 0 {
				bn = ltb
				break Loop
			}
			ncol = 0
		case '\t':
			ncol++
			ltb = bn
		}
	}
	if oldbn > f.Nbox {
		return oldbn
	}
	if bn >= f.Nbox {
		return f.Nbox
	}
	if bn <= oldbn {
		return oldbn
	}
	return bn + 1
}

// StartCell returns the first box in the cell
func (f *Run) StartCell(bn int) int {
	//	println(bn)
	if bn == 0 {
		return 0
	}
	ncols := 0
	nrows := 0
	//	oldbn := bn
	bn = f.EndLine(bn)
	lsb := bn
	if bn == f.Nbox {
		//nrows++
	}
	var b *Box
Loop:
	for bn-1 != 0 {
		b = &f.Box[bn-1]
		switch b.Break() {
		case '\n':
			if ncols == 0 {
				if nrows == 0 {
					return 0
				}
				break Loop
			}
			lsb = bn
			nrows++
			ncols = 0
		case '\t':
			ncols++
		default:
		}
		bn--
	}
	if ncols == 0 {
		return lsb
	}
	if bn-1 == 0 {
		if f.Box[bn-1].Break() == '\n' {
			return bn
		}
		return bn - 1
	}
	return lsb
	//	println("bn-1", bn-1)
	//	f.DumpBoxes()
	if bn-1 == 0 && f.Box[bn-1].Break() != '\n' {
		//return 0
	}
	if f.Box[bn].Break() == '\t' && f.Box[bn-1].Break() != '\n' {
		return bn
	}
	//	println("return", bn)
	return bn
}

// NextCell is like EndCell, except it doesn't assume bn
// is part of a cell. It skips past the current cell under
// bn and any non-cellular boxes afterward, returning the
// starting box of the next cell or f.Nbox
func (f *Run) NextCell(bn int) int {
	bn = f.EndCell(bn)
	oldbn := bn
	for ; bn != f.Nbox; bn++ {
		b := &f.Box[bn]
		if b.Break() == '\t' {
			bn = f.StartCell(bn)
			break
		}
	}
	if bn <= oldbn {
		return oldbn
	}
	return bn
}

func (f *Run) StartLine(bn int) int {
	for ; bn-1 >= 0; bn-- {
		b := &f.Box[bn-1]
		if b.Break() == '\n' {
			break
		}
	}
	return bn
}

func (f *Run) EndLine(bn int) int {
	for bn < f.Nbox {
		b := &f.Box[bn]
		if b.Break() == '\n' {
			break
		}
		bn++
	}
	return bn
}

func (f *Run) NextLine(bn int) int {
	bn = f.EndLine(bn)
	if bn < f.Nbox {
		return bn + 1
	}
	return bn
}

func (f *Run) PrevLine(bn int) int {
	for ; bn >= 0; bn-- {
		b := &f.Box[bn]
		if b.Break() == '\n' {
			break
		}
	}
	if bn == -1 && f.Box[0].Break() == '\n' {
		return 0
	}
	for bn-1 >= 0 {
		b := &f.Box[bn-1]
		if b.Break() == '\n' {
			break
		}
		bn--
	}
	return bn
}
