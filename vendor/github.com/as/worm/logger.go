// Package worm provides Write-Once Read-Many (WORM) logging semantics for
// sequential log read access and append only write access. It also provides
// a coalescer to compact records in transit to the underlying log.
//
// A log is an record of data defined in github.com/as/event that implements
// the event.Record interface
//
package worm

import (
	"fmt"
	"github.com/as/event"
)

type Logger interface{
	// Write appends the record to the log
	Write(event.Record) (err error)

	// ReadAt reads and returns log record n
	ReadAt(at int64) (event.Record, error)

	// Len returns the number of records
	Len() int64

}

// NewLogger returns a Write-Once Read-Many (WORM) logger capable of
// serializing an ordered stream of event.Records.
func NewLogger() Logger{
	return &logWORM{}
}

type logWORM struct {
	rec []event.Record
}

// ReadAt reads and returns log record n
func (l *logWORM) ReadAt(n int64) (event.Record,  error){
	if n < 0 || n >= int64(len(l.rec)){
		return nil, fmt.Errorf("bad read offset: %d\n", n)
	}
	return l.rec[n], nil
}

// Write writes v to the tail of the log
func (l *logWORM) Write(v event.Record) (err error){
	l.rec = append(l.rec, v)
	return nil
}

// Len returns the number of records stored the log
func (l *logWORM) Len() int64{
	return int64(len(l.rec))
}
