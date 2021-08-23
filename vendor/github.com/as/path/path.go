// Package path simulates a current working directory in a text
// editing environment. It provides a way to keep track of a
// relative display path that represents a larger absolute path.
//
// The primary operation is 'Look', which inspects a third junction
// to the managed path. The goal of this package is to separate
// the file system specific nature of relative/absolute path and
// provide a straightforward structure that a text editor can consume.
//
// Please be aware that none of the methods for path.Path have
// pointer references. Their state will not change, and they
// return a new path if necessary.
package path

import "path/filepath"
import "os"
import "log"

// New returns a path string representing the base path
// to begin traversal. If the provided path is relative
// the returned path is computed from the process's
// current working directory
func New(path string) (t Path) {
	if path == "" || path == "." || !filepath.IsAbs(path) {
		// An assumption is made here that os.Getwd()
		// will never returns a non-directory without
		// error.
		//
		// Remember the absolute base path but display
		// the relative path given
		wd := getawd()
		return Path{base: wd, disp: Clean(path)}
	}
	return Path{base: path}
}

func getwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return Clean(wd)
}
func getawd() string {
	awd, err := filepath.Abs(getwd())
	if err != nil {
		panic(err)
	}
	return awd
}

// Path consists of a base path and a display path. The
// display path is what you see in a text editor. The base
// path holds the prefix to that display path so that combining
// them creates a valid absolute path to the file.
//
type Path struct {
	base string
	disp string
}

// Path returns itself it it's pointing to a file, otherwise
// it returns the path of the directory containing the file
func (t Path) Dir() Path {
	if t.IsDir() {
		return t
	}
	return t.Look("..")
}

// Name returns the display name of the path
func (t Path) Name() string {
	return Clean(t.disp)
}

// Blank returns a copy of Path without a display name set
// The path points to the base path
func (t Path) Blank() Path {
	t.disp = ""
	return t
}

func (t Path) Exists() bool {
	return Exists(t.Abs())
}

// Base returns the base path
func (t Path) Base() string {
	return t.base
}

// IsDir returns true if the path is a directory
func (t Path) IsDir() bool {
	return IsDir(t.Abs())
}

// Path returns an absolute path of the current state.
func (t Path) Abs() string {
	if filepath.IsAbs(t.disp) && t.disp == t.base {
		return Clean(t.base)
	}
	return Clean(filepath.Join(t.base, t.disp))
}

// Look returns a new state. An absolute path returns a state with a new base
// and display path set to that path. A relative path adds on to the existing
// display path unless the path consists of enough double-dots to erase the
// display path. In that case, the state of both base and display path is
// set to the join of the base path and the double-dot path.
func (t Path) Look(dir string) (p Path) {
	defer func() {
		log.Println(p)
	}()
	wasdot := dir == "."
	dir = Clean(dir)
	if filepath.IsAbs(dir) {
		return Path{disp: dir, base: dir}
	}
	abs := filepath.Join(t.base, t.disp)
	if abs0 := filepath.Join(filepath.Join(abs, dir), ".."); len(abs0) < len(t.base) {
		t.base = abs0
		t.disp, _ = filepath.Rel(abs0, abs)
	}
	t.disp = filepath.Join(t.disp, dir)
	if t.disp == "." && !wasdot {
		t.disp = t.base
	}
	return Path{base: Clean(t.base), disp: Clean(t.disp)}
}

/*
func main() {
	fmt.Println(NewPath("."))
	fmt.Println(Path{`/windows/system32/`, `drivers/etc/`}.Look("../hosts").Look("..").Look("/"))
}
*/
