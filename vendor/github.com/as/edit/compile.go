package edit

import (
	"errors"
	"fmt"
	"io"

	"github.com/as/event"
	"github.com/as/text"
	"github.com/as/worm"
)

var (
	ErrNilFunc   = errors.New("empty program")
	ErrNilEditor = errors.New("nil editor")
)

var (
	noop = func(ed Editor) {}
)

type Options struct {
	Sender Sender
	Origin string
}

type Command struct {
	fn       func(Editor)
	s        string
	args     string
	next     *Command
	Emit     *Emitted
	modified bool
}

func MustCompile(s string) (cmd *Command) {
	cmd, err := Compile(s)
	if err != nil {
		panic(fmt.Sprintf("MustCompile: %s\n", err))
	}
	return cmd
}

// Compile runs the build steps on the input string and returns
// a runnable command.
func Compile(s string, opts ...*Options) (cmd *Command, err error) {
	_, itemc := lex("cmd", s)
	p := parse(itemc, opts...)
	err = <-p.stop
	return compile(p), err
}

// Modified returns true if the last call to c.Run() modified the contents
// of the editor
func (c *Command) Modified() bool {
	return c.modified
}

// Func returns a function entry point that operates on a Editor
func (c *Command) Func() func(Editor) {
	return c.fn
}

func net(hist worm.Logger) (ins, del int64) {
	for i := int64(0); i < hist.Len(); i++ {
		t, _ := hist.ReadAt(int64(i))
		switch t := t.(type) {
		case *event.Insert:
			ins += t.Q1 - t.Q0
		case *event.Delete:
			del += t.Q1 - t.Q0
		}
	}
	return
}

// Commit plays back the history onto ed, starting from
// the last event to the first in reverse order. This is
// useful only when hist contains a set of independent events
// applied as a transaction where shifts in the address offsets
// are not observed.
//
// If the command is
//		a,abc,
//		x,.,a,Q,
// The result is:
//		abc -> aQbQcQ
// The log should contain
// 		i 1 Q
// 		i 2 Q	(not i 3 Q)
// 		i 3 Q (not i 5 Q)
//
// Commit will only reallocate ed's size once. If ed implements
// io.WriterAt, a write-through fast path is used to commit the
// transaction.
func Commit(ed Editor, hist worm.Logger) (err error) {
	//	log.Printf("commit: content: %q", ed.Bytes())
	_, del := net(hist)
	for i := int64(hist.Len()) - 1; i >= 0; i-- {
		e, err := hist.ReadAt(i)
		if err != nil {
			return err
		}
		switch t := e.(type) {
		case *event.Write:
			ed.(io.WriterAt).WriteAt(t.P, t.Q0)
			//			log.Printf("event[%d]: %#v\n", i, e)
		case *event.Insert:
			ed.Insert(t.P, t.Q0+del)
			//			log.Printf("event[%d]: %#v\n", i, e)
		case *event.Delete:
			del -= int64(t.Q1 - t.Q0)
			ed.Delete(t.Q0, t.Q1)
			//			log.Printf("event[%d]: %#v\n", i, e)
		}
	}
	return err
}

func (c *Command) ck(ed Editor) error {
	c.modified = false
	if ed == nil {
		return ErrNilEditor
	}
	if c.fn == nil {
		return ErrNilFunc
	}
	return nil
}

// Transcribe runs the compiled program on ed
func (c *Command) Transcribe(ed Editor) (log worm.Logger, err error) {
	if err = c.ck(ed); err != nil {
		return nil, err
	}
	log = worm.NewLogger()
	hist := text.NewHistory(&Recorder{ed}, log)
	c.Emit.Dot = c.Emit.Dot[:0]
	c.fn(hist)
	return log, nil
}

// Run runs the compiled program on ed
func (c *Command) RunTransaction(ed Editor) (err error) {
	hist, err := c.Transcribe(ed)
	if err != nil {
		return err
	}
	c.modified = hist.Len() > 0
	return Commit(ed, hist)
}

// Run runs the compiled program on ed
func (c *Command) Run(ed Editor) (err error) {
	return c.RunTransaction(ed)
}

func (c *Command) oldRun(ed Editor) (err error) {
	if err = c.ck(ed); err != nil {
		return err
	}
	c.Emit.Dot = c.Emit.Dot[:0]
	c.fn(ed)
	return nil
}

// Next returns the next instruction for the compiled program. This
// effectively steps through x,..., and y,...,
func (c *Command) Next() *Command {
	return c.next
}

func (c *Command) nextFn() func(f Editor) {
	if c.next == nil {
		return nil
	}
	return c.next.fn
}

func compileAddr(a Address) func(f Editor) {
	if a == nil {
		return noop
	}
	return a.Set
}

func compile(p *parser) (cmd *Command) {
	for i := range p.cmd {
		if i+1 == len(p.cmd) {
			break
		}
		p.cmd[i].next = p.cmd[i+1]
	}
	fn := func(f Editor) {
		addr := compileAddr(p.addr)
		if addr != nil {
			addr(f)
		}
		if p.cmd != nil && p.cmd[0] != nil && p.cmd[0].fn != nil {
			p.cmd[0].fn(f)
		}
	}
	return &Command{fn: fn, Emit: p.Emit}
}
