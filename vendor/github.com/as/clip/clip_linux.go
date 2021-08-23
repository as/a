package clip

import (
	"bytes"
	"os/exec"
)

type Clip struct{}

func New() (*Clip, error) {
	return &Clip{}, nil
}

func (c *Clip) Read(p []byte) (n int, err error) {
	cmd := exec.Command("xclip", "-o")
	b := new(bytes.Buffer)
	cmd.Stdout = b
	err = cmd.Run()
	return copy(p, b.Bytes()), err
}

func (c *Clip) Write(p []byte) (n int, err error) {
	n, err = c.writeTo(p, "clipboard")
	n, err = c.writeTo(p, "primary")
	return c.writeTo(p, "primary")
}

func (c *Clip) writeTo(p []byte, dst string) (n int, err error) {
	b := bytes.NewBuffer(p)
	cmd := exec.Command("xclip", "-i", "-selection", dst)
	cmd.Stdin = b
	return len(p) - (len(p) - b.Len()), cmd.Run()
}
