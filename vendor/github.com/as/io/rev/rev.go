package rev

import "io"

type Reader struct {
	b []byte
	i int
}

func NewReader(b []byte) *Reader {
	return &Reader{
		b: b,
		i: len(b),
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	pl := len(p)
	for {
		r.i--
		if r.i < 0 {
			if n == 0 {
				return 0, io.EOF
			}
			return n, nil
		}
		if n >= pl {
			return n, err
		}
		p[n] = r.b[r.i]
		n++
	}
}
