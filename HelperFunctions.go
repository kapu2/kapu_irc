package main

func SpliceIsSame[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TwoDSpliceIsSame[T comparable](a [][]T, b [][]T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !SpliceIsSame(a[i], b[i]) {
			return false
		}
	}
	return true
}

func FindRunePos(line []rune, toFind rune) int {
	pos := -1
	for i, r := range line {
		if toFind == r {
			pos = i
			break
		}
	}
	if pos != -1 {
		return pos
	} else {
		return pos
	}
}
