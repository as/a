package text

// Copy-on Reader
// This bad name will be changed

import (
	"time"

	"github.com/as/worm"
)

// COR processes a stream of event records and writes them to a coalescing logger.
// The event stream is saved to the provided logger for later use.
type COR struct {
	ed     Editor
	rec    *Recorder       // Transcodes function calls to event Records
	co     *worm.Coalescer // Batches those records
	hist   worm.Logger     // Stores those records
	noi    int
	nod    int
	q0, q1 int64
}

func NewCOR(ed Editor, hist worm.Logger) *COR {
	c := &COR{
		ed:   ed,
		hist: hist,
		co:   worm.NewCoalescer(hist, time.Second*3),
	}
	c.rec = NewRecorder(c.co)
	return c
}

func (c *COR) Stats() (int64, int64) {
	return int64(c.noi), int64(c.nod)
}
func (c *COR) Len() int64 {
	return c.ed.Len()
}
func (c *COR) Bytes() []byte {
	return c.ed.Bytes()
}
func (c *COR) Select(q0, q1 int64) {
	c.q0 = q0
	c.q1 = q1
}
func (c *COR) Dot() (q0, q1 int64) {
	q0 = clamp(c.q0, 0, c.ed.Len())
	q1 = clamp(c.q1, 0, c.ed.Len())
	return
}
func (c *COR) Insert(p []byte, q0 int64) int {
	n := c.rec.Insert(p, q0)
	c.noi += n
	return n
}
func (c *COR) Delete(q0, q1 int64) int {
	n := c.rec.Delete(q0, q1)
	c.nod += n
	return n
}
func (c *COR) Close() (err error) {
	c.ed = nil
	return err
}
func (c *COR) Flush() (err error) {
	return c.co.Flush()
}
