package dump

import "log"

const Debug = 0

func Printf(fm string, v ...interface{}) {
	if Debug == 0 {
		return
	}
	log.Printf(fm, v...)
}
