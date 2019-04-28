package lmodutf8

import (
	"unicode/utf8"

	"ofunc/lua"
)

func index(l *lua.State, n, i, j int) (int, int) {
	if i < 0 {
		i = n + i + 1
	}
	if j < 0 {
		j = n + j + 1
	}
	if i <= 0 || j > n {
		panic("index out of range")
	}
	return i, j
}

func liter(l *lua.State, s string, i, j int) func() (rune, int, int) {
	xs := []byte(s)
	i, j = index(l, len(xs), i, j)
	return func() (rune, int, int) {
		if i > j {
			return utf8.RuneError, 0, 0
		}
		r, n := utf8.DecodeRune(xs[i-1:])
		i += n
		return r, n, i - n
	}
}

func riter(l *lua.State, s string, i, j int) func() (rune, int, int) {
	xs := []byte(s)
	i, j = index(l, len(xs), i, j)
	return func() (rune, int, int) {
		if i > j {
			return utf8.RuneError, 0, 0
		}
		r, n := utf8.DecodeLastRune(xs[:j])
		j -= n
		return r, n, j + 1
	}
}
