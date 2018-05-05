package main

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type (
	PlumbAction func(msg *Plumbmsg) error
	PlumbRule   func(msg *Plumbmsg) PlumbAction
)

type Plumbmsg struct {
	src, dst string
	wdir     string
	kind     string
	Attr
	Data []byte
}
type Attr map[string]string

func (p *Plumbmsg) Arg() string {
	return strings.TrimSpace(string(p.Data))
}

var httpLink = NewRegexp("^http(s?)://", openBrowser)

func NewRegexp(expr string, action PlumbAction) *regexpRule {
	return &regexpRule{regexp.MustCompile(expr), action}
}

var openBrowser = PlumbAction(func(p *Plumbmsg) (err error) {
	for _, browser := range browsers() {
		full := append(strings.Fields(browser), p.Arg())
		cmd := exec.Command(full[0], full[1:]...)
		err = cmd.Start()
		if err == nil {
			break
		}
	}
	return err
})

type regexpRule struct {
	*regexp.Regexp
	action PlumbAction
}

func (r *regexpRule) Plumb(msg *Plumbmsg) (matched bool, error error) {
	if r.Match(msg.Data) {
		return true, r.action(msg)
	}
	return false, nil
}

func matches(match bool, err error) bool {
	if match && err == nil {
		return true
	}
	if err != nil {
		logf("plumber: %s", err)
	}
	return match
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
