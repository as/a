# WORM

Package worm provides Write-Once Read-Many (WORM) log storage semantics for
ordered read access and append only write access to log records. It also
provides a coalescer for compacting them.

A `log` is an record of data defined in `github.com/as/event` that
implements the `event.Record` interface

# TYPES

```
type Coalescer struct {
    Logger
    // contains filtered or unexported fields
}
    Coalescer coalesces logs written to it until the deadband expires. After
    expiration, the coalesced log is flushed to the underlying logger upon
    the next call to Write().

func NewCoalescer(lg Logger, deadband time.Duration) *Coalescer
    NewCoalescer wraps the given logger and returns a coalescer

func (l *Coalescer) Flush() error
    Flush flushes the last unwritten log to the underlying logger

func (l *Coalescer) ReadAt(n int64) (event.Record, error)
    ReadAt reads and returns log record n

func (l *Coalescer) Write(v event.Record) (err error)
    Write writes v to the tail of the log

type Logger interface {
    // Write appends the record to the log
    Write(event.Record) (err error)

    // ReadAt reads and returns log record n
    ReadAt(at int64) (event.Record, error)

    // Len returns the number of records
    Len() int64
}

func NewLogger() Logger
    NewLogger returns a Write-Once Read-Many (WORM) logger capable of
    serializing an ordered stream of event.Records.
```
