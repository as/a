package main

import "flag"

func init() {
	flag.Parse()
}

func argparse() (list []string) {
	if len(flag.Args()) > 0 {
		list = append(list, flag.Args()...)
	} else {
		list = append(list, "guide")
		list = append(list, ".")
	}
	return
}
