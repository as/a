package box

import "unicode/utf8"

//import "log"

func (r *Run) ensure(nb int) {
	if nb == r.Nalloc {
		r.Grow(r.delta)
		if r.delta < 32768 {
			r.delta *= 2
		}
	}
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *Run) zRunescan(s []byte, ymax int) {
	r.Nbox = 0
	r.Nchars = 0
	r.Nchars += int64(len(s))
	i := 0
	nb := 0

	adv := 0
	for nl := 0; nl <= ymax; nb++ {
		if nb == r.Nalloc {
			r.Grow(r.delta)
			if r.delta < 32768 {
				r.delta *= 2
			}
		}
		i += adv
		if i == len(s) {
			break
		}
		c := s[i]
		switch c {
		default:
			for _, c := range string(s[i:min(len(s), MaxBytes)]) {
				if c == '\t' || c == '\n' {
					break
				}
				adv = utf8.RuneLen(c)
			}
			r.Box[nb] = Box{
				Nrune: i,
				Ptr:   s[i : i+adv],
				Width: r.Face.Dx(s[i : i+adv]),
			}
		case '\t':
			adv = 1
			r.Box[nb] = Box{
				Nrune:    -1,
				Ptr:      s[i : i+adv],
				Width:    r.minDx,
				Minwidth: r.minDx,
			}
		case '\n':
			adv = 1
			r.Box[nb] = Box{
				Nrune: -1,
				Ptr:   s[i : i+adv],
				Width: r.maxDx,
			}
			nl++
		}
	}
	r.Nchars -= int64(len(s))
	r.Nbox += nb
}

func (r *Run) Runescan(s []byte, ymax int) {
	r.Boxscan(s, ymax)
}
func (r *Run) Boxscan(s []byte, ymax int) {
	r.Nbox = 0
	r.Nchars = 0
	r.Nchars += int64(len(s))
	i := 0
	nb := 0

	for nl := 0; nl <= ymax; nb++ {
		if nb == r.Nalloc {
			r.Grow(r.delta)
			if r.delta < 32768 {
				r.delta *= 2
			}
		}
		if i == len(s) {
			break
		}
		i++
		c := s[i-1]
		switch c {
		default:
			for _, c = range s[i:min(len(s), MaxBytes)] {
				if special(c) {
					break
				}
				i++
			}
			r.Box[nb] = Box{
				Nrune: i,
				Ptr:   s[:i],
				Width: r.Face.Dx(s[:i]),
			}
		case '\t':
			r.Box[nb] = Box{
				Nrune:    -1,
				Ptr:      s[:i],
				Width:    r.minDx,
				Minwidth: r.minDx,
			}
		case '\n':
			r.Box[nb] = Box{
				Nrune: -1,
				Ptr:   s[:i],
				Width: r.maxDx,
			}
			nl++
		}
		s = s[i:]
		i = 0
	}
	r.Nchars -= int64(len(s))
	r.Nbox += nb
}

func special(c byte) bool {
	return c == '\t' || c == '\n'
}
