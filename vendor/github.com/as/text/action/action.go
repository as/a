package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/as/edit"
	"github.com/as/event"
	"github.com/as/text"
	"github.com/as/text/find"
)

func clean(dir string) string {
	dir = filepath.ToSlash(dir)
	dir = filepath.FromSlash(dir)
	return filepath.Clean(dir)
}
func LookPath(dir, name string) interface{} {
	name, addr := SplitPath(name)
	if name == "" && addr != "" {
		return edit.MustCompile(addr)
	}
	name = filepath.Clean(name)
	if !filepath.IsAbs(name) {
		pref := filepath.Clean(dir)
		if !isdir(pref) {
			pref = filepath.Dir(pref)
		}
		name = filepath.Clean(filepath.Join(pref, name))
	}
	if isdir(name) {
		if addr != "" {
			// A directory with an address doesn't make sense
			// user probably refers to a file on another system
			// with the same name as the dir, so look should fail
			return nil
		}
		return event.Get{Path: name, IsDir: true}
	} else if isfile(name) {
		return event.Get{Path: name, Addr: addr}
	}
	return nil
}

func Look(w text.Editor, q0, q1 int64) interface{} {
	if q0 == q1 {
		q1 = find.Accept(w.Bytes(), q1, []byte(string(find.AlphaNum)+`\/.:`))
		q0 = find.Acceptback(w.Bytes(), q0, []byte(string(find.AlphaNum)+`\/.:`))
	}
	return event.Select{event.Rec{Q0: q0, Q1: q1}}
}
func Dirof(file string) string {
	file = clean(file)
	return filepath.Dir(file)
}
func IsDir(path string) bool  { return isdir(path) }
func IsFile(path string) bool { return isfile(path) }
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
func isfile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func SplitPath(path string) (name, addr string) {
	name = path
	x := strings.Index(name, ":")
	if x == -1 {
		// No colon, so we're done
		return name, ""
	}

	if x == 0 {
		if len(name) == 1 {
			// An invalid empty file name or address
			// ambiguous.
			return ":", ""
		}
		// Empty name
		return "", name[1:]
	}

	// Dealing with Windows and their "drive letters"
	// Default behavior is to do what Unix, Linux, and Plan9 expect
	// Windows should be last in line
	if x == 1 && strings.IndexAny(name, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") == 0 {
		if isdir(name[:2]) {
			n, a := SplitPath(name[2:])
			return name[:2] + n, a
		}
	}
	return name[:x], name[x+1:]

}
