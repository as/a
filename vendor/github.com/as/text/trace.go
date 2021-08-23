package text

import "fmt"

// Trace wraps ed with methods that logs the arguments passed to each
// method call of the text.Editor interface for debugging purposes.
func Trace(ed Editor) Editor {
	return &tracer{ed}
}

type tracer struct {
	Editor
}

/*
func (x *tracer) WriteAt(p []byte, at int64) (int, error){
	fmt.Printf("WriteAt: %q @ %d\n", p, at)
	return x.Editor.(io.WriterAt).WriteAt(p, at)
}
*/
func (x *tracer) Insert(p []byte, at int64) int {
	fmt.Printf("Insert: %q @ %d\n", p, at)
	return x.Editor.Insert(p, at)
}
func (x *tracer) Delete(q0, q1 int64) int {
	fmt.Printf("Delete: %d:%d\n", q0, q1)
	return x.Editor.Delete(q0, q1)
}
func (x *tracer) Dot() (q0, q1 int64) {
	fmt.Printf("Dot: %d:%d\n", q0, q1)
	return x.Editor.Dot()
}
func (x *tracer) Select(q0, q1 int64) {
	fmt.Printf("Select: %d:%d\n\n", q0, q1)
	x.Editor.Select(q0, q1)
}
func (x *tracer) Len() int64 {
	fmt.Printf("Len: %d\n\n", x.Editor.Len())
	return x.Editor.Len()
}
func (x *tracer) Bytes() []byte {
	fmt.Printf("Bytes: %q\n\n", x.Editor.Bytes())
	return x.Editor.Bytes()
}
