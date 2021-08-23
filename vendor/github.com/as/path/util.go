package path

import (
	"os"
	"path/filepath"
)

var sep = string(filepath.Separator)

// Clean returns a clean path that contains the file
// separator if that path ends in a file name.
func Clean(path string) string {
	return filepath.Clean(path)
	path = filepath.Clean(path)
	if IsDir(path) {
		println("not a dir", path)
		if len(path) > 0 && path[len(path)-1] != filepath.Separator {
			path += string(filepath.Separator)
		}
	}
	return path
}

// DirOf returns the directory of the path. If the
// path ends in a directory, it returns the path itself
func DirOf(path string) string {
	if IsDir(path) {
		return path
	}
	return filepath.Dir(path)
}

// FileOf is like DirOf, except for non-directories
func FileOf(path string) string {
	if !Exists(path) || IsDir(path) {
		return ""
	}
	return filepath.Base(path)
}

// IsDir returns true if path ends in a directory
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		if err == os.ErrNotExist {
			return false
		}
		return false
	}
	return fi.IsDir()
}

// Exists returns true if path is reachable
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
