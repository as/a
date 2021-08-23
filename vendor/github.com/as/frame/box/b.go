package box

type B interface {
	Seek(bn int, whence int) int
	Next() bool
	Prev() bool
	Box() *Box
}
type boxes struct {
	bn int
	b  []Box
}

func clamp(v, l, h int) int {
	if v < l {
		return l
	}
	if v > h {
		return h
	}
	return v
}
func (b *boxes) Prev() bool {
	if b.bn == 0 {
		return false
	}
	b.bn--
	return true
}
func (b *boxes) Next() bool {
	if b.bn+1 == len(b.b) {
		return false
	}
	b.bn++
	return true
}
func (b *boxes) Box() *Box {
	return &b.b[b.bn]
}

func (b *boxes) Seek(bn int, whence int) int {
	oldbn := b.bn
	switch whence {
	case 0:
		b.bn = clamp(bn, 0, len(b.b))
	case 1:
		b.bn = clamp(b.bn+bn, 0, len(b.b))
	case 2:
		b.bn = clamp(b.bn+bn, 0, len(b.b))
	}
	return oldbn
}

func PrevLine(bx B) bool {
	if bx.Seek(0, 1) == 0 {
		return false
	}
	for bx.Box().Break() != '\n' && bx.Prev() {
	}
	if bx.Seek(0, 1) == 0 && bx.Box().Break() == '\n' {
		return true
	}
	for bx.Prev() && bx.Box().Break() != '\n' {
	}
	if bx.Box().Break() == '\n' {
		return bx.Next()
	}
	return true
}
