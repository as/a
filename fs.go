package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/as/srv/fs"
)

type fileresolver struct {
	fs.Fs
	err error
}

type pathinfo struct {
	root string
	tag  string
	pred string
}

type fileinfo struct {
	os.FileInfo
	path string
	dir  string
}

func (r *fileresolver) set(f *fileinfo, name string) (ok bool) {
	f.path = name
	f.dir = name
	if f.FileInfo != nil && !f.FileInfo.IsDir() {
		f.dir = filepath.Dir(f.path)
	}
	return r.err == nil
}

func (f *fileinfo) String() string {
	if f == nil {
		return "<nil>"
	}
	return fmt.Sprint(f.FileInfo)
}

func (r *fileresolver) isAbs(name string) bool {
	const Letters = "cCABDEFGHIJKLMNOPQRSTUVWXYZabdefghijklmnopqrstuvwxyz"
	if len(name) == 0 {
		return false
	}
	if path.IsAbs(name) || filepath.IsAbs(name) || name[0] == '\\' || name[0] == '/' {
		return true
	}
	return len(name) > 1 && name[1] == ':' && strings.ContainsAny(name[:1], Letters)
}

// joindir resolves a to a file or directory and runs filepath.Join on
// the directory of a. If a is already a directory the operation is
// join(a,b), if it's a file the operation is join(a/.., b).
func (r *fileresolver) joindir(a, b string) string {
	info, err := r.Stat(a)
	if err != nil || !info.IsDir() {
		a = filepath.Dir(a)
	}
	return filepath.Join(a, b)
}

func (r *fileresolver) look(pi pathinfo) (f fileinfo, ok bool) {
	if r.isAbs(pi.pred) {
		// last element is absolute;
		f.FileInfo, r.err = r.Stat(pi.pred)
		return f, r.set(&f, filepath.Clean(pi.pred))
	}

	if r.isAbs(pi.tag) {
		file := r.joindir(pi.tag, pi.pred)
		f.FileInfo, r.err = r.Stat(file)
		return f, r.set(&f, filepath.Clean(file))
	}

	// need to know if
	pi.tag = filepath.Join(pi.root, pi.tag)
	file := r.joindir(pi.tag, pi.pred)
	f.FileInfo, r.err = r.Stat(file)
	return f, r.set(&f, filepath.Clean(file))
}
