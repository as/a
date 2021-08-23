package cursor

import (
	"fmt"
	"strconv"
	"strings"
)

func parse12(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s[:11]))
}

func ReadString(s string) (x, y, btn, ms int, err error) {
	if len(s) != 49 {
		err = fmt.Errorf("s != 49 (%d)\n", len(s))
		return
	}
	if s[0] != 'm' {
		err = fmt.Errorf("bad header: %c", s[0])
		return
	}
	s = s[1:]
	if x, err = parse12(s); err != nil {
		err = fmt.Errorf("x: %d", x)
		return
	}
	s = s[12:]
	if y, err = parse12(s); err != nil {
		err = fmt.Errorf("y: %d", y)
		return
	}
	s = s[12:]
	if btn, err = parse12(s); err != nil {
		err = fmt.Errorf("btn: %x", btn)
		return
	}
	s = s[12:]
	if ms, err = parse12(s); err != nil {
		err = fmt.Errorf("ms: %d", ms)
		return
	}
	return
}
func WriteString(x, y, btn, ms int) (s string) {
	return fmt.Sprintf("m%11d %11d %11d %11d ", x, y, btn, ms)
}
