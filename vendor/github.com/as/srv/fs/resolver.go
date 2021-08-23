package fs

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Resolver struct {
	Fs
	err error
}

type Path struct {
	Root string
	Tag  string
	Pred string
}

type FileInfo struct {
	os.FileInfo
	Path string
	Dir  string
}

func (f *FileInfo) String() string {
	if f == nil {
		return "<nil>"
	}
	return fmt.Sprint(f.FileInfo)
}

func (r *Resolver) set(f *FileInfo, name string) (ok bool) {
	f.Path = name
	f.Dir = name
	if f.FileInfo != nil && !f.FileInfo.IsDir() {
		f.Dir = filepath.Dir(f.Path)
	}
	return r.err == nil
}

func (r *Resolver) IsAbs(name string) bool {
	const Letters = "cCABDEFGHIJKLMNOPQRSTUVWXYZabdefghijklmnopqrstuvwxyz"
	if len(name) == 0 {
		return false
	}
	if path.IsAbs(name) || filepath.IsAbs(name) || name[0] == '\\' || name[0] == '/' {
		return true
	}
	return len(name) > 1 && name[1] == ':' && strings.ContainsAny(name[:1], Letters)
}

// Join resolves a to a file or directory and runs filepath.Join on
// the directory of a. If a is already a directory the operation is
// join(a,b), if it's a file the operation is join(a/.., b).
func (r *Resolver) Join(a, b string) string {
	info, err := r.Stat(a)
	if err != nil || !info.IsDir() {
		a = filepath.Dir(a)
	}
	return filepath.Join(a, b)
}

func (r *Resolver) Look(pi Path) (f FileInfo, ok bool) {
	if r.IsAbs(pi.Pred) {
		// last element is absolute;
		f.FileInfo, r.err = r.Stat(pi.Pred)
		return f, r.set(&f, filepath.Clean(pi.Pred))
	}

	if r.IsAbs(pi.Tag) {
		file := r.Join(pi.Tag, pi.Pred)
		f.FileInfo, r.err = r.Stat(file)
		return f, r.set(&f, filepath.Clean(file))
	}

	// need to know if
	pi.Tag = filepath.Join(pi.Root, pi.Tag)
	file := r.Join(pi.Tag, pi.Pred)
	f.FileInfo, r.err = r.Stat(file)
	return f, r.set(&f, filepath.Clean(file))
}
