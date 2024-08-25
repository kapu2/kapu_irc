package main

import (
	"unicode/utf8"
)

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

func TwoDStringAreSame(a []string, b []string) bool {
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

func Split[K comparable](toSplit []K, delimiter K) [][]K {
	var ret [][]K
	start := 0
	i := 0
	var v K
	for i, v = range toSplit {
		if v == delimiter {
			ret = append(ret, toSplit[start:i])
			start = i + 1
		}
	}
	if len(toSplit) > 0 {
		remaining := toSplit[start : i+1]
		if len(remaining) > 0 {
			ret = append(ret, remaining)
		}
	}

	return ret
}

func RemoveFirstRuneFromString(s string) string {
	if utf8.RuneCountInString(s) > 0 {
		if utf8.RuneCountInString(s) == 1 {
			s = ""
		} else {
			s = string([]rune(s)[1:])
		}
	}
	return s
}
