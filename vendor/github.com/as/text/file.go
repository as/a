package text

/*
type ReadWriteSeekCloser interface{
	io.ReadWriteSeeker
	io.Closer
}

func (f *file) Insert(p []byte, at int64) (n int){
	at, err = f.f.Seek(at, io.SeekStart)
	if err==nil{return}
	dx := len(p)
	dy := 4096
	tmp := make([]byte, dy)
	sp := f.f.Len()-1
	for {
		_, _ = f.f.Seek(sp, io.SeekCurrent)
		if sp-dy < at+dx{
			sp -= dy+(sp-dy) (at+dx)
		}

	}


}

type File struct{
	f ReadWriteSeekCloser
	Insert(p []byte, at int64) (n int)
	Delete(q0, q1 int64) (n int)
	Len() int64
	ReadAt(p []byte, at int64) (n int, err error)
	WriteAt(p []byte, at int64) (n int, err error)
}


*/
