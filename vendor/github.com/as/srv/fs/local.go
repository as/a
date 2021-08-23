package fs

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Local struct {
}

func (f *Local) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
func (f *Local) Get(name string) (data []byte, err error) {
	// Don't rely on the system to read a directory as a file. We may be on a legacy OS.
	if isdir(name) {
		return formatdir(name)
	}
	return ioutil.ReadFile(name)
}

func (f *Local) Put(name string, data []byte) (err error) {
	return ioutil.WriteFile(name, data, 0666)
}
func (f *Local) Cmd(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	return exec.CommandContext(ctx, name, arg...), nil
}

// formatdir returns the files in the directory as a line of tab
// seperated names
func formatdir(path string) ([]byte, error) {
	fi, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(fi, func(i, j int) bool {
		if fi[i].IsDir() && !fi[j].IsDir() {
			return true
		}
		ni, nj := fi[i].Name(), fi[j].Name()
		return strings.Compare(ni, nj) < 0
	})
	b := new(bytes.Buffer)
	sep := ""
	for _, v := range fi {
		nm := v.Name()
		if v.IsDir() {
			nm += string(os.PathSeparator)
		}
		fmt.Fprintf(b, "%s%s", sep, nm)
		sep = "\t"
	}
	return b.Bytes(), nil
}

func isdir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		if err == os.ErrNotExist {
			return false
		}
		fmt.Println(err)
		return false
	}
	return fi.IsDir()
}
