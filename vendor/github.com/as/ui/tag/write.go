package tag

type Address interface {
	Dot() (q0, q1 int64)
}
type Range [2]int64

func (r Range) Dot() (int64, int64) {
	return r[0], r[1]
}

/*
func (t *Tag) Write(p []byte) (n int, err error) {
	return t.Body.Write(p)
}
func (t *Tag) Insert(p []byte, q0 int64) int {
	return t.Body.Insert(p, q0)
}
func (t *Tag) Delete(q0,q1 int64) int {
	return t.Body.Delete(q0,q1)
}
func (t *Tag) Dot() (q0,q1 int64) {
	return t.Body.Dot()
}
func (t *Tag) Bytes() []byte {
	return t.Body.Bytes()
}
func (t *Tag) Len() int64 {
	return t.Body.Len()
}
func (t *Tag) Select() int64 {
	return t.Body.Len()
}
*/