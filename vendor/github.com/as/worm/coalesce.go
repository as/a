package worm

import (
	"time"
	"github.com/as/event"
)

// Coalescer coalesces logs written to it until the deadband expires. After
// expiration, the coalesced log is flushed to the underlying logger upon
// the next call to Write().
type Coalescer struct {
	Logger

	last event.Record

	// period during which coalesced writes can be
	// buffered without flush to the underlying logger
	deadband time.Duration
	flushc chan chan error
	writec chan event.Record
	timer *time.Timer
}

// NewCoalescer wraps the given logger and returns a coalescer
func NewCoalescer(lg Logger, deadband time.Duration) *Coalescer {
	c := &Coalescer{
		Logger:   lg,
		last:     nil,
		deadband: deadband,
		timer: time.NewTimer(deadband),
	}
	c.run()
	return c
}

// ReadAt reads and returns log record n
func (l *Coalescer) ReadAt(n int64) (event.Record, error) {
	return l.Logger.ReadAt(n)
}

func (l *Coalescer) combine(v event.Record) bool{
	next := l.last.Coalesce(v)
	// log.Printf("result \n\t\t%#v\n", next)
	if next == nil{
		return false
	}
	l.last=next
	return true
}

func (l *Coalescer) run(){
	l.flushc = make(chan chan error)
	l.writec = make(chan event.Record)
	l.reclock()
	go func(){
	for{
		select{
		case <- l.timer.C:
			// deadline expired, flush what we have now
			l.flush()
			l.reclock()
		case v := <- l.writec:
			if l.last == nil{
				// keep going
				l.last = v
				continue
			}
			fused := l.combine(v)
			if !fused{
				l.flush()
				l.last = v
			}
			l.reclock() // deadline extended 
		case donec := <- l.flushc:
			// the user did this with a public function
			l.flush()
			l.reclock()
			donec <- nil
			return
		}
	}
	}()
}

func (l *Coalescer) reclock(){
			if !l.timer.Stop() {
				<-l.timer.C
			}
			l.timer.Reset(l.deadband)
}

// Write writes v to the tail of the log
func (l *Coalescer) Write(v event.Record) (err error) {
	l.writec <- v
	return nil 
}

// Flush flushes the last unwritten log to the underlying logger
func (l *Coalescer) Flush() error{
	donec := make(chan error)
	l.flushc <- donec
	return <- donec
}
func (l *Coalescer) flush() error {
	if l.last == nil{
		return nil
	}
	l.Logger.Write(l.last)
	switch e := l.last.(type){
	case *event.Write:
		if e.Residue != nil{
			l.last = e.Residue
			l.flush()
		}
	}
	l.last = nil
	return nil
}
