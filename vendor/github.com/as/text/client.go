package text

import (
	"errors"
	"io"
)

var ErrNilBuffer = errors.New("Nil Buffer")

// Open returns an Editor capable of managing a selection
// on b. The selection is maintained automatically as long
// as insertions and deletions happen through the returned
// editor.
func Open(b Buffer) (w Editor, err error) {
	if b == nil {
		return nil, ErrNilBuffer
	}
	return &client{b, 0, 0}, nil
}

// client is a non-graphical editor that maintains a coherent
// selection across insertion/deletions and maintains a local
// buffer.
type client struct {
	Buffer
	q0, q1 int64
}

func (c *client) WriteAt(p []byte, q0 int64) (n int, err error) {
	if t, ok := c.Buffer.(io.WriterAt); ok {
		n, err = t.WriteAt(p, q0)
		//q0, _ = c.clamp(q0, q0)
		//c.q0, c.q1 = Coherence(-1, q0, q0+int64(n), c.q0, c.q1)
		//c.q0, c.q1 = Coherence(1, q0, q0+int64(n), c.q0, c.q1)
		return n, err
		t = t
	}
	n = c.Insert(p, q0)
	q0 += int64(n)
	q1 := q0 + int64(n)
	c.Delete(q0, q1)
	return
}

func (c *client) clamp(q0, q1 int64) (int64, int64) {
	nr := c.Buffer.Len()
	return clamp(q0, 0, nr), clamp(q1, 0, nr)
}
func (c *client) Dot() (q0, q1 int64) {
	return c.clamp(c.q0, c.q1)
}
func (c *client) Select(q0, q1 int64) {
	c.q0, c.q1 = c.clamp(q0, q1)
}
func (c *client) Insert(s []byte, q0 int64) (n int) {
	n = c.Buffer.Insert(s, q0)
	q0, _ = c.clamp(q0, q0)
	c.q0, c.q1 = Coherence(1, q0, q0+int64(n), c.q0, c.q1)
	return n
}
func (c *client) Delete(q0, q1 int64) (n int) {
	n = c.Buffer.Delete(q0, q1)
	q0, q1 = c.clamp(q0, q1)
	c.q0, c.q1 = Coherence(-1, q0, q1, c.q0, c.q1)
	return n
}
func (c *client) Close() error {
	return c.Buffer.Close()
}
