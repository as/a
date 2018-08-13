package main

import (
	"os/exec"
	"strings"
)

var guruModes = "callees callers callstack definition describe freevars implements peers pointsto referrers what whicherrs"

func (g *Grid) Selection() string {
	return string(g.Tag.Label.Rdsel())
}

func (g *Grid) guru(label string, q0, q1 int64) (advance bool, err error) {
	if !strings.HasSuffix(label, ".go") {
		return true, nil
	}
	asel := g.Selection()
	mode := ""
	//	scope := "."
	for _, v := range strings.Fields(guruModes) {
		if asel == v {
			mode = v
			break
		}
	}
	if mode == "" {
		return true, nil
	}

	data, err := exec.Command(
		"guru",
		//		"-scope",
		//		scope,
		mode,
		addrfmt(label, q0, q1),
	).CombinedOutput()

	if err != nil {
		g.aerr("guru: %s", err)
	}
	if len(data) != 0 {
		g.aguru("%s", data)
	}
	return false, err
}
