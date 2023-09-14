/*
Copyright 2019 by ofunc

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package lmodutf8

import (
	"unicode/utf8"

	"github.com/ofunc/lua"
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
