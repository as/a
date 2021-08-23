// +build darwin

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
	cmd := exec.Command("pbpaste")
	b := new(bytes.Buffer)
	cmd.Stdout = b
	err = cmd.Run()
	return copy(p, b.Bytes()), err
}

func (c *Clip) Write(p []byte) (n int, err error) {
	cmd := exec.Command("pbcopy")
	b := bytes.NewBuffer(p)
	cmd.Stdin = b
	return len(p) - b.Len(), cmd.Run()
}
