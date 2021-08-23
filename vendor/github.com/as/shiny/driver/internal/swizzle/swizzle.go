package swizzle

var (
	swizzler = pureBGRA
)

func Swizzle(p, q []byte) {
	if len(p) < 4 {
		return
	}
	swizzler(p, q)
}

func pureBGRA(p, q []byte) {
	if len(p)%4 != 0 {
		return
	}
	for i := 0; i < len(p); i += 4 {
		q[i+0], q[i+1], q[i+2], q[i+3] = p[i+2], p[i+1], p[i+0], p[i+3]
	}
}
