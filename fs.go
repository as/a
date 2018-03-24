package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/as/srv/fs"
)

type tagresolver struct {
}

type fileresolver struct {
	fs.Fs
	err error
}

type pathinfo struct {
	wd   string
	tag  string
	name string
	src  interface{}
}

type fileinfo struct {
	os.FileInfo
	abspath, visible string
	err              error
}

func (r *fileresolver) set(f *fileinfo, name string) (ok bool) {
	if r.err == nil {
		f.abspath = name
		f.visible = name
		return true
	}
	return false
}

func (f *fileinfo) String() string {
	if f == nil {
		return "<nil>"
	}
	return fmt.Sprintf("abspath: %s\nvisible: %s\nfi: %#v\n", f.abspath, f.visible, f.FileInfo)
}

func (r *fileresolver) isAbs(name string) bool {
	const Letters = "cCABDEFGHIJKLMNOPQRSTUVWXYZabdefghijklmnopqrstuvwxyz"
	if path.IsAbs(name) || filepath.IsAbs(name) || name[0] == '\\' {
		return true
	}
	return len(name) > 1 && name[1] == ':' && strings.ContainsAny(name[:1], Letters)
}

func (r *fileresolver) look(pi pathinfo) (f fileinfo, ok bool) {
	if r.isAbs(pi.name) {
		f.FileInfo, r.err = r.Stat(pi.name)
		return f, r.set(&f, filepath.Clean(pi.name))
	}

	if !r.isAbs(pi.tag) {
		pi.tag = filepath.Join(pi.wd, pi.tag)
	}

	f.FileInfo, r.err = r.Fs.Stat(pi.tag)
	if r.err != nil {
		return f, false
	}

	if !f.FileInfo.IsDir() {
		pi.tag = filepath.Join(pi.tag, "..")
	}
	f.abspath = filepath.Join(pi.tag, pi.name)
	f.visible = f.abspath
	return f, r.err == nil
}
