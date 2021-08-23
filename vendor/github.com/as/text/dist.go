package text

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func min3(a, b, c int) int {
	return min(c, min(a, b))
}

func Distance(s, t []byte) int {
	dtab := make([]int, len(s)+1)

	for i := range dtab {
		dtab[i] = i
	}
	for j := 1; j <= len(t); j++ {
		last := dtab[0]
		dtab[0] = j
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				dtab[i], last = last, dtab[i]
			} else {
				dtab[i], last = min3(last, dtab[i], dtab[i-1])+1, dtab[i]
			}
		}
	}
	return dtab[len(s)]
}
