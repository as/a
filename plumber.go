package main

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type Plumbmsg struct {
	src, dst string
	wdir     string
	kind     string
	Attr
	Data []byte
}

func (p *Plumbmsg) Arg() string {
	return strings.TrimSpace(string(p.Data))
}

type Attr map[string]string

var isHTTP = regexp.MustCompile("^http(s?)://")
var match = isHTTP
var actionFunc = func(p *Plumbmsg) (err error) {
	for _, browser := range browsers() {
		full := append(strings.Fields(browser), p.Arg())
		cmd := exec.Command(full[0], full[1:]...)
		err = cmd.Start()
		if err == nil {
			break
		}
	}
	return err
}

// TODO(as): This is an overly-concrete implementation equivalent
// to one plumber rule.
func PlumberExp(msg *Plumbmsg) bool {
	if match.Match(msg.Data) {
		actionFunc(msg)
		return true
	}
	return false
}

// browsers returns a list of commands to attempt for web visualization.
func browsers() []string {
	cmds := []string{"chrome", "google-chrome", "firefox"}
	switch runtime.GOOS {
	case "darwin":
		return append(cmds, "/usr/bin/open")
	case "windows":
		return append(cmds, "cmd /c start")
	default:
		userBrowser := os.Getenv("BROWSER")
		if userBrowser != "" {
			cmds = append([]string{userBrowser, "sensible-browser"}, cmds...)
		} else {
			cmds = append([]string{"sensible-browser"}, cmds...)
		}
		return append(cmds, "xdg-open")
	}
}
