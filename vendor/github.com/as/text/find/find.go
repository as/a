package find

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/as/io/rev"
	"github.com/as/text"
)

var (
	Lefts        = [...]byte{'(', '{', '[', '<', '\''}
	Rights       = [...]byte{')', '}', ']', '>', '\''}
	Free         = [...]byte{'"', '\'', '`', '\n'}
	AlphaNum     = []byte("*&!%-_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	AlphaNumFile = []byte(`*&!%+#-_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/:,\.`)
)

func Run(br io.ByteScanner, s []byte) (n int) {
	for {
		b, err := br.ReadByte()
		if err != nil && err != io.EOF {
			return
		}
		if bytes.IndexAny(s, string(b)) == -1 {
			br.UnreadByte()
			return
		}
		n++
		if err != nil {
			break
		}
	}
	return n
}

func Any(b byte, s []byte) int64 {
	for i, v := range s {
		if b == v {
			return int64(i)
		}
	}
	return -1
}

func Isany(b byte, s []byte) bool {
	for _, v := range s {
		if b == v {
			return true
		}
	}
	return false
}

func ExpandFile(p []byte, i int64) (int64, int64) {
	j := Accept(p, i, AlphaNumFile)
	i = Acceptback(p, i, AlphaNumFile)
	return i, j
}

func Expand(p []byte, i int64) (int64, int64) {
	j := Accept(p, i, AlphaNum)
	i = Acceptback(p, i, AlphaNum)
	return i, j
}

func FreeExpand(f text.Editor, i int64) (int64, int64) {
	p := f.Bytes()
	if len(p) == 0 || i < 0 || i > int64(len(p))-1 {
		return i, i
	}
	if i != 0 {
		if which := Any(p[i-1], Lefts[:]); which != -1 {
			br := bytes.NewReader(f.Bytes()[i:])
			d0 := Parity(br, Lefts[which], Rights[which], false)
			if d0 != -1 {
				return i, i + int64(d0) - 1
			}
		}
		if i < int64(len(p)) {
			if which := Any(p[i], Free[:]); which != -1 {
				j := Findback(p, i-1, []byte{Free[which]})
				if j != -1 {
					return i - (i - j), i
				}
			}
		}
		if which := Any(p[i-1], Free[:]); which != -1 {
			j := Find(p, i, []byte{Free[which]})
			if j != -1 {
				return i, j
			}
		}
	}
	return Expand(p, i)
}

func Acceptback(p []byte, i int64, sep []byte) int64 {
	q0 := i
	for ; q0-1 >= 0 && Isany(p[q0-1], sep); q0-- {
	}
	return q0
}
func Accept(p []byte, q1 int64, sep []byte) int64 {
	for ; q1 < int64(len(p)) && Isany(p[q1], sep); q1++ {
	}
	return q1
}

func Findback(p []byte, i int64, sep []byte) int64 {
	q0 := i
	for ; q0-1 >= 0 && !Isany(p[q0-1], sep); q0-- {
	}
	if q0 < 0 {
		return i
	}
	return q0
}
func Find(p []byte, j int64, sep []byte) int64 {
	q1 := j
	for ; q1 != int64(len(p)) && !Isany(p[q1], sep); q1++ {
	}
	if q1 == int64(len(p)) {
		return j
	}
	return q1
}
func Parity(br io.ByteReader, l byte, r byte, back bool) (n int) {
	if back {
		panic("unimplemented")
	}
	push := 1
	var err error
	for ; err == nil && push != 0; n++ {
		b, err := br.ReadByte()
		if err != nil && err != io.EOF {
			break
		}
		if b == l {
			push++
		} else if b == r {
			push--
		}
		if err == io.EOF {
			break
		}
	}
	if err != nil && err != io.EOF || push != 0 {
		return -1
	}
	return n
}
func findParity(br io.ByteReader, l byte, r byte, back bool) (n int) {
	if back {
		panic("unimplemented")
	}
	push := 1
	var err error
	for ; err == nil && push != 0; n++ {
		b, err := br.ReadByte()
		if err != nil && err != io.EOF {
			break
		}
		if b == l {
			push++
		} else if b == r {
			push--
		}
	}
	println("n is")
	println(n)
	if err != nil && err != io.EOF || push != 0 {
		return -1
	}
	return n
}

func FindNext(f text.Editor, i, j int64, text []byte) (q0, q1 int64) {
	//	i, j := f.Dot()
	p := f.Bytes()
	x := text
	q0 = int64(bytes.Index(p[j:], x))
	if q0 == -1 {
		q0 = int64(bytes.Index(p[:i], x))
		if q0 < 0 {
			return i, j
		}
	} else {
		q0 += j
	}
	q1 = q0 + int64(len(x))
	return q0, q1
}

func Findlinerev(p []byte, org, N int64) (q0, q1 int64) {
	N = -N + 1
	p0 := p
	p = p[:org]
	q0, q1 = Findline2(N, rev.NewReader(p)) // 0 = len(p)-1
	l := q1 - q0
	q0 = org - q1
	q1 = q0 + l
	q0 = q1 - l
	if q0 >= 0 && q0 < int64(len(p0)) && p0[q0] == '\n' {
		q0++
	}
	return
}
func Findline3(p []byte, org, N int64) (q0, q1 int64) {
	p = p[org:]
	q0, q1 = Findline2(N, bytes.NewReader(p))
	return q0 + org, q1 + org
}

// Put	Edit 354
func Findline2(N int64, r io.Reader) (q0, q1 int64) {
	br := bufio.NewReader(r)
	nl := int64(0)
	for nl != N {
		b, err := br.ReadByte()
		if err != nil {
			break
		}
		q1++
		if b == '\n' {
			nl++
			if nl == N {
				break
			}
			q0 = q1
		}
	}
	return
}

var (
	regexpNL     = regexp.MustCompile(`\n`)
	regexpLine   = regexp.MustCompile(`^.*$`)
	regexpSpaces = regexp.MustCompile(`\s+`)
	regexpWords  = regexp.MustCompile(`\w+`)
)

func Next(p []byte, i, j int64) (q0, q1 int64, err error) {
	defer func(r0, r1 int64) {
		fmt.Printf("Next: [%d:%d]->[%d:%d]\n", r0, r1, q0, q1)
	}(i, j)
	x := p[i:j]
	q0 = int64(bytes.Index(p[j:], x))
	if q0 == -1 {
		q0 = int64(bytes.Index(p[:i], x))
		if q0 < 0 {
			return i, j, io.EOF
		}
	} else {
		q0 += j
	}
	q1 = q0 + int64(len(x))
	return q0, q1, nil
}

func FindLineAt(r io.SectionReader, off int64, dy int) (q0, q1 int64) {
	//	q0, q1 = FindLine(io.NewSectionReader(r, off, r.Size()), dy)
	return q0 + off, q1 + off
}

func FindLine(r io.Reader, dy int) (q0, q1 int64) {
	br := bufio.NewReader(r)
	for n := 0; n != dy; {
		b, err := br.ReadByte()
		if err != nil {
			break
		}
		q1++
		if b == '\n' {
			n++
			if n == dy {
				break
			}
			q0 = q1
		}
	}
	return
}
func FindParity(f text.Editor) (q0, q1 int64, ok bool) {
	q0, q1 = f.Dot()
	for i := range Lefts {
		d0 := findParity(bytes.NewReader(f.Bytes()[q0:]), Lefts[i], Rights[i], false)
		if q0 != -1 {
			return q0, q0 + int64(d0) - 1, true
		}
	}
	return q0, q1, false
}
