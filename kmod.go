package main

const (
	KShift = 1 << iota
	KCtrl
	KAlt
	KMeta
)

/*

type KFlag uint64

// kmod should work well as long as there is only one writer
// to kmod; which is true at the time of writing. future changes
// may introduce bizzare inconsistencies undetectable by the race
// detector
var kmod = new(uint64)

func KMod() KFlag {
	return KFlag(atomic.LoadUint64(kmod))
}
func SetKMod(k KFlag) {
	atomic.StoreUint64(kmod, uint64(k))
}
func KModOn(s KFlag) {
	SetKMod(KMod() | s)
}
func KModOff(s KFlag) {
	SetKMod(KMod() &^ s)
}

*/
