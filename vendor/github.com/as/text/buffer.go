package text

type Buffer interface {
	Insert(p []byte, at int64) (n int)
	Delete(q0, q1 int64) (n int)
	Write(p []byte) (n int, err error)
	Len() int64
	Bytes() []byte
	Close() error
}

func NewBuffer() Buffer {
	return &buf{
		R: make([]byte, 0, 64*1024),
	}
}

func BufferFrom(b []byte) Buffer {
	return &buf{
		R: b,
	}
}

type buf struct {
	Q0, Q1 int64
	R      []byte
}

func (w *buf) Len() int64 {
	return int64(len(w.R))
}

func (w *buf) Select(q0, q1 int64) {
	if q0 < 0 {
		return
	}
	if q1 > w.Len() {
		q1 = w.Len()
	}
	w.Q0 = q0
	w.Q1 = q1
}

// Write appends the contents of p to the buffer.
func (w *buf) Write(p []byte) (int, error) {
	w.R = append(w.R, p...)
	return len(p), nil
}

func (w *buf) WriteAt(p []byte, offset int64) (n int, err error) {
	n = copy(w.R[offset:], p)
	return
}

func (w *buf) Insert(s []byte, q0 int64) (n int) {
	if n = len(s); n == 0 {
		return 0
	}
	if q0 < 0 {
		// Let's be precise and annoying
		// 0 is the real lower bound
		return 0
	}
	if q0 == 0 {
		w.R = append(s, w.R...) // append(s, w.R[q0:]...)...)
		return n
	}
	if q0 >= w.Len() { // Common case: append
		w.R = append(w.R, s...)
		return n
	}
	// Interpolate
	w.R = append(w.R[:q0], append(s, w.R[q0:]...)...)
	return n
}

func (w *buf) Delete(q0, q1 int64) (n int) {
	if q1 < q0 || q0 < 0 {
		return 0
	}
	if n = int(q1 - q0); n == 0 {
		return 0
	}
	nr := w.Len()
	copy(w.R[q0:], w.R[q1:][:nr-q1])
	w.R = w.R[:nr-int64(n)]
	return int(n)
}

func (w *buf) Close() error {
	w.R = nil
	return nil
}

func (w *buf) Dot() (q0, q1 int64) {
	q0 = w.Q0
	q1 = w.Q1
	return
}

func (w *buf) Dirty() bool {
	return false
}

func (w *buf) Bytes() []byte {
	return w.R
}

/*
	// TODO(as): Below is the correct implementation of WriteAt, which
	// appends to the buffer

	// WriteAt writes the contents of p at the given offset. It writes
	// through existing data in the buffer and appends to it if p overflows
	// the buffer at that offset.
	func (w *buf) WriteAt(p []byte, offset int64) (n int, err error) {
		if offset >= int64(len(p)) {
			w.R = append(w.R, p...)
			return len(p), nil
		}

		if n = copy(w.R[offset:], p); n == len(p){
			return n, nil
		}

		w.R = append(w.R, p[n:]...)
		return len(p), nil
	}
*/
