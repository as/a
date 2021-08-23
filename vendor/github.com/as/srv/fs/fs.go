package fs

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type Command interface {
	CombinedOutput() ([]byte, error)
	Output() ([]byte, error)
	Run() error
	Start() error
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Wait() error
}

type Fs interface {
	Stat(name string) (os.FileInfo, error)
	Get(name string) (data []byte, err error)
	Put(name string, data []byte) (err error)
	Cmd(ctx context.Context, name string, arg ...string) (*exec.Cmd, error)
}
