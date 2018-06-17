package main

import (
	"flag"
	"runtime"
)

func defaultFaceSize() int {
	switch runtime.GOOS {
	case "darwin":
		return 13
	default:
		return 11
	}
}

var (
	utf8       = flag.Bool("u", false, "enable utf8 experiment")
	elastic    = flag.Bool("elastic", false, "enable elastic tabstops")
	oled       = flag.Bool("b", false, "OLED display mode (black)")
	ftsize     = flag.Int("ftsize", defaultFaceSize(), "font size")
	srvaddr    = flag.String("l", "", "(dangerous) announce and serve file system clients on given endpoint")
	clientaddr = flag.String("d", "", "dial to a remote file system on endpoint")
	quiet      = flag.Bool("q", false, "dont interact with the graphical subsystem (use with -l)")
)

var _ int = parse()

func parse() int {
	flag.Parse()
	return 0
}
