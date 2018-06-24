package dump

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Tag struct {
	Kind byte // 'F', 'f', 'e', 'x', 'i', 'v', 'h'

	Row, Col int64
	Q0, Q1   int64
	Percent  int64
	NBody    int

	Wid   int
	NTag  int
	Dir   bool
	Dirty bool

	Label []byte
	Body  []byte
}

const Supported = "Ff"

func (k Tag) Encode(w io.Writer) (err error) {
	switch k.Kind {
	case 'F':
		_, err = fmt.Fprintf(w,
			"%c%11d %11d %11d %11d %11d %11d \n%11d %11d %11d %11d %11d %s\n%s",
			k.Kind, k.Row, k.Col, k.Q0, k.Q1, k.Percent, len(k.Body),
			k.Wid, len(k.Label), len(k.Body), b2i(k.Dir), b2i(k.Dirty), k.Label, k.Body,
		)
	case 'f':
		_, err = fmt.Fprintf(w,
			"%c%11d %11d %11d %11d %11d \n%11d %11d %11d %11d %11d %s\n",
			k.Kind, k.Row, k.Col, k.Q0, k.Q1, k.Percent,
			k.Wid, len(k.Label), len(k.Body), b2i(k.Dir), b2i(k.Dirty), k.Label,
		)
	default:
		return fmt.Errorf("encode: not supported: %q", k.Kind)
	}
	return err
}

func (k *Tag) Decode(b *bufio.Reader) error {
	var label, body bytes.Buffer
	if err := k.DecodeTo(b, &label, &body); err != nil {
		return err
	}
	k.Label = label.Bytes()
	k.Body = body.Bytes()
	return nil
}
func (k *Tag) DecodeTo(b *bufio.Reader, label, body io.Writer) error {
	r := NewScanner(b)
	var nbody, ntag int64

	if !r.Scan(&k.Kind) {
		return r.Err()
	}

	switch k.Kind {
	case 'F':
		r.Scan(&k.Row, &k.Col, &k.Q0, &k.Q1, &k.Percent, &nbody, '\n', &k.Wid, &ntag, &nbody, &k.Dir, &k.Dirty)
	case 'f':
		r.Scan(&k.Row, &k.Col, &k.Q0, &k.Q1, &k.Percent, '\n', &k.Wid, &ntag, &nbody, &k.Dir, &k.Dirty)
		nbody = 0
	default:
		return fmt.Errorf("decode: not supported: %q", k.Kind)
	}

	if r.Err() != nil {
		return fmt.Errorf("decode: %q: %v", k.Kind, r.Err())
	}
	if _, err := io.Copy(label, io.LimitReader(r, ntag)); err != nil {
		return err
	}
	if !r.Scan('\n') {
		return fmt.Errorf("decode: %v", r.Err())
	}
	_, err := io.Copy(body, io.LimitReader(r, nbody))
	return err
}
func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
