package dump

import (
	"bufio"
	"fmt"
	"os"
)

func Create(path string) (*writer, error) {
	fd, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &writer{
		fd:     fd,
		Writer: bufio.NewWriter(fd),
	}, nil
}

type writer struct {
	fd *os.File
	*bufio.Writer
	err error
}

func (w *writer) Line(s string) {
	if w.err != nil {
		return
	}
	fmt.Fprintln(w, s)
}

func (w *writer) Int(n int, delim byte) {
	fmt.Fprintf(w, "%11d%c", n, delim)
}

func (w *writer) Ints(n ...int) {
	if w.err != nil {
		return
	}
	sep := byte(' ')
	for i, v := range n {
		if i+1 == len(n) {
			sep = '\n'
		}
		w.Int(v, sep)
	}
}

func (w *writer) Close() error {
	w.Writer.Flush()
	return w.fd.Close()
}
