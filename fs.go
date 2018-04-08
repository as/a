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
	wd   string
	tag  string
	name string
}

type fileinfo struct {
	os.FileInfo
	abspath, visible string
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
	if len(name) == 0 {
		return false
	}
	if path.IsAbs(name) || filepath.IsAbs(name) || name[0] == '\\' {
		return true
	}
	return len(name) > 1 && name[1] == ':' && strings.ContainsAny(name[:1], Letters)
}

func (r *fileresolver) look(pi pathinfo) (f fileinfo, ok bool) {
	logf("resolver: look %v", pi)
	if r.isAbs(pi.name) {
		logf("resolver: A: absolute path")
		f.FileInfo, r.err = r.Stat(pi.name)
		logf("resolver: A: %v", f.FileInfo)
		return f, r.set(&f, filepath.Clean(pi.name))
	}

	if !r.isAbs(pi.tag) {
		logf("resolver: B: nottag: %v", pi.tag)
		pi.tag = filepath.Join(pi.wd, pi.tag)
		logf("resolver: B: joined: %v", pi.tag)
	}

	f.FileInfo, r.err = r.Fs.Stat(pi.tag)
	if r.err != nil {
		fi2, err := r.Fs.Stat(filepath.Dir(pi.tag))
		if err != nil {
			logf("resolver: C: fileinfo: %v err=%s", f, r.err)
			return f, false
		}
		f.FileInfo = fi2
		r.err = nil
	}

	if !f.FileInfo.IsDir() {
		logf("resolver: D: not dir: %v", f)
		pi.tag = filepath.Join(pi.tag, "..")
	}
	f.abspath = filepath.Join(pi.tag, pi.name)
	f.visible = f.abspath
	logf("resolver: E: abspath: %v", f.abspath)
	return f, r.err == nil
}
