package main

import "flag"

func argparse() (list []string) {
	if len(flag.Args()) > 0 {
		list = append(list, flag.Args()...)
		showbanner = false
	} else {
		list = append(list, "guide")
		list = append(list, ".")
		showbanner = true
	}
	return
}
