package main

func SizeThirdOf(size int) int {
	return size - (size / 5 * 2)
}
func SizeSmall(size int) int {
	if size-200 < 0 {
		return 100
	}
	return size - 200
}
