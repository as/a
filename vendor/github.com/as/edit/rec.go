package edit

// Recorder is an overlay over an Editor that prevents mutable
// changes from occuring.
type Recorder struct {
	Editor
}

func (r *Recorder) Write(p []byte) (int, error) {
	return len(p), nil
}
func (r *Recorder) WriteAt(p []byte, at int64) (int, error) {
	return len(p), nil
}
func (r *Recorder) Insert(p []byte, at int64) (n int) {
	return len(p)
}
func (r *Recorder) Delete(q0, q1 int64) (n int) {
	return int(q1 - q0)
}
