package edit

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type (
	Append    struct{ Data []byte }
	Insert    struct{ Data []byte }
	Delete    struct{}
	Change    struct{ To []byte }
	ReadFile  struct{ Name string }
	WriteFile struct{ Name string }
	Pipe      struct{ To string }
	Trade     struct{ Address }
	S         struct {
		*regexp.Regexp
		ReplaceAmp
		Repl  string
		Limit int64
	}
)

func (c Append) Apply(ed Editor) {
	_, q1 := ed.Dot()
	ed.Insert(c.Data, q1)
}
func (c Insert) Apply(ed Editor) {
	q0, _ := ed.Dot()
	ed.Insert(c.Data, q0)
}

func (c Delete) Apply(ed Editor) {
	ed.Delete(ed.Dot())
}

func (c Change) Apply(ed Editor) {
	q0, q1 := ed.Dot()
	del := q1 - q0
	ins := int64(len(c.To))
	if del < ins {
		// write del bytes through insert the rest
		ed.(io.WriterAt).WriteAt(c.To[:del], q0)
		ed.Insert(c.To[del:], q0+del)
	} else if del > ins {
		// delete del-ins bytes write the rest through
		ed.Delete(q0+ins, q1)
		ed.(io.WriterAt).WriteAt(c.To, q0)
	} else {
		ed.(io.WriterAt).WriteAt(c.To, q0)
	}
}

func (c ReadFile) Apply(ed Editor) {
	data, err := ioutil.ReadFile(c.Name)
	if err != nil {
		eprint(err)
		return
	}
	q0, q1 := ed.Dot()
	if q0 != q1 {
		ed.Delete(q0, q1)
	}
	ed.Insert(data, q0)
}

func (c WriteFile) Apply(ed Editor) {
	fd, err := os.Create(c.Name)
	if err != nil {
		eprint(err)
		return
	}
	defer fd.Close()
	q0, q1 := ed.Dot()
	_, err = io.Copy(fd, bytes.NewReader(ed.Bytes()[q0:q1]))
	if err != nil {
		eprint(err)
	}
}

func (c Pipe) Apply(ed Editor) {
	x := strings.Fields(c.To)
	if len(x) == 0 || x[0] == "" {
		eprint("|: nothing on rhs")
	}
	n := x[0]
	var a []string
	if len(x) > 1 {
		a = x[1:]
	}
	q0, q1 := ed.Dot()
	cmd := exec.Command(n, a...)
	cmd.Stdin = bytes.NewReader(append([]byte{}, ed.Bytes()[q0:q1]...))
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		eprint(err)
	}
	Change{To: buf.Bytes()}.Apply(ed)
}
func (c S) Apply(ed Editor) {
	sp, ep := ed.Dot()
	buf := bytes.NewReader(ed.Bytes()[sp:ep])
	q0 := int64(0)
	for i := int64(1); q0+sp != ep; i++ {
		loc := c.FindReaderIndex(buf)
		if loc == nil {
			break
		}
		q1 := q0 + int64(loc[1])
		q0 += int64(loc[0])
		ed.Select(sp+q0, sp+q1)

		if i == c.Limit || c.Limit == -1 {
			buf := c.ReplaceAmp.Gen(ed.Bytes()[q0:q1])
			Change{buf}.Apply(ed)
			if i == 500000 {
				break
			}
		}
		q0 = q1
		buf.Seek(q0, 0)
	}
	ed.Select(ep, ep)
}
