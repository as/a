package tag

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/as/frame"
	"github.com/as/text"
)

// Open returns an Editor capable of managing a selection
// on b. The selection is maintained automatically as long
// as insertions and deletions happen through the returned
// editor.
func Open(b text.Buffer) (w text.Editor, err error) {
	if b == nil {
		return nil, fmt.Errorf("bad buffer")
	}
	return &client{b, 0, 0, 0, nil}, nil
}

type client struct {
	text.Buffer
	q0, q1, s int64
	col       *frame.Color
}

func (c *client) Dot() (q0, q1 int64) {
	return c.q0, c.q1
}

func (c *client) Mark(s int64) {
	c.s = s
}

func (c *client) Select(q0, q1 int64) {
	if q0 > q1 {
		q1, q0 = q0, q1
	}
	if c.col != nil {
		if t, ok := c.Buffer.(text.Scroller); ok {
			org := t.Origin()
			if t, ok := c.Buffer.(painter); ok {
				framepaint(t, c.col, c.s-org, q0-org, q1-org, c.q0-org, c.q1-org)
			}
		}
	}
	c.q0, c.q1 = q0, q1
}
func (c *client) Insert(s []byte, q0 int64) (n int) {
	return n
}
func (c *client) Delete(q0, q1 int64) (n int) {
	return n
}
func (c *client) Read(p []byte) (n int, err error) {
	q0, q1 := c.Dot()
	data := c.Bytes()[q0:q1]
	if len(data) == 0 {
		return 0, io.EOF
	}
	n = copy(p, data)
	if n > 0 {
		c.Select(q0+int64(n), q1)
	}
	return
}
func (c *client) String() string {
	data, _ := ioutil.ReadAll(c)
	return string(data)
}
