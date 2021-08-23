package fs

import (
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
)

var (
	ErrStrlen = errors.New("string too long")
)

// readBytes reads a 64bit length-prefixed byte slice
// in the form: [8]n [n]byte
func readBytes(r io.Reader, max int64) ([]byte, error) {
	n := int64(0)
	err := binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return nil, err
	}

	if n > max {
		return nil, ErrStrlen
	}
	data, err := ioutil.ReadAll(io.LimitReader(r, n))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// writeString writes a 64bit length-prefixed  byte slice
// in the form: [8]n [n]byte
func writeBytes(w io.Writer, s []byte) (err error) {
	if err = binary.Write(w, binary.BigEndian, int64(len(s))); err != nil {
		return err
	}
	n, err := w.Write(s)
	if err != nil {
		return err
	}
	if n != len(s) {
		panic("writeString: err != nil && short write")
	}
	return nil
}

func readString(r io.Reader, max int64) (string, error) {
	s, err := readBytes(r, max)
	return string(s), err
}
func writeString(w io.Writer, s string) error {
	return writeBytes(w, []byte(s))
}
