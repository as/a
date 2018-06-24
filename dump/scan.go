package dump

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		Reader: bufio.NewReader(r),
	}
}

type Scanner struct {
	*bufio.Reader
	err error
}

func (s *Scanner) Scan(v ...interface{}) bool {
	for _, v := range v {
		if !s.scan(v) {
			return false
		}
	}
	return s.err == nil
}
func (s *Scanner) scan(v interface{}) bool {
	if s.err != nil {
		return false
	}
	switch v := v.(type) {
	case *[]int:
		const maxelem = 100
		var max = maxelem
		if len(*v) != 0 {
			max = len(*v)
		}
		i := 0
		for ; i != max; i++ {
			n := s.parseInt(s.readn(11))
			delim := s.parseByte(s.readn(1))
			*v = append(*v, n)
			if delim == '\n' || s.err != nil {
				break
			}
		}
		if i == maxelem {
			s.err = fmt.Errorf("too many fields")
		}
	case *bool:
		var x int
		if s.scan(&x) && x != 0 {
			*v = true
		}
	case *int64:
		var x int
		if s.scan(&x) {
			*v = int64(x)
		}
	case *int:
		*v = s.parseInt(s.readn(11))
		if s.err != nil {
			return false
		}
		if delim := s.parseByte(s.readn(1)); delim != ' ' && delim != '\n' {
			s.err = fmt.Errorf("int: %v: delim: want space, have %q", *v, delim)
			return false
		}
	case *string:
		*v = s.readline()
	case *rune:
		*v = rune(s.parseByte(s.readn(1)))
	case *byte:
		*v = s.parseByte(s.readn(1))
	case rune:
		s.scan(byte(v))
	case byte:
		if b := s.readn(1); s.err != nil || b[0] != v {
			if s.err == nil {
				s.err = fmt.Errorf("bad byte: want %q, have %q", v, b[0])
			}
			return false
		}
	case interface{}:
		s.err = fmt.Errorf("bad type: %T", v)
	}
	return s.err == nil
}
func (s *Scanner) Err() error {
	return s.err
}
func (s *Scanner) readline() (v string) {
	if s.err != nil {
		return ""
	}
	v, s.err = s.ReadString('\n')
	return strings.TrimSuffix(v, "\n")
}
func (s *Scanner) readn(n int) string {
	if s.err != nil {
		return ""
	}
	tmp := make([]byte, n)
	n, s.err = io.ReadFull(s, tmp[:])
	Printf("read %d: %v: %q", n, s.err, tmp[:n])
	return string(tmp[:n])
}
func (s *Scanner) parseByte(v string) byte {
	if s.err != nil || len(v) == 0 {
		return 0
	}
	return v[0]
}
func (s *Scanner) parseInt(v string) int {
	if s.err != nil {
		return 0
	}
	var f float64
	f, s.err = strconv.ParseFloat(strings.TrimSpace(v), 64)
	return int(f)
}
