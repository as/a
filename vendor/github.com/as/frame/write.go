package frame

import "io"

// Write implements io.Writer. The write operation appends to the current
// selection given by Dot(). It returns io.EOF when the entire message doesn't
// fit on the frame.
func (f *Frame) Write(p []byte) (n int, err error) {
	_, p1 := f.Dot()
	if n = f.Insert(p, p1); n == 0 {
		return 0, io.EOF
	}
	return n, nil
}
