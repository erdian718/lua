package lmodutf8

import (
	"strconv"
	"unicode/utf8"

	"ofunc/lua"
)

// Open opens the module.
func Open(l *lua.State) int {
	l.NewTable(0, 8)

	l.Push("char")
	l.Push(lchar)
	l.SetTableRaw(-3)

	l.Push("charpattern")
	l.Push("[\x00-\x7F\xC2-\xF4][\x80-\xBF]*")
	l.SetTableRaw(-3)

	l.Push("codes")
	l.Push(lcodes)
	l.SetTableRaw(-3)

	l.Push("codepoint")
	l.Push(lcodepoint)
	l.SetTableRaw(-3)

	l.Push("len")
	l.Push(llen)
	l.SetTableRaw(-3)

	l.Push("offset")
	l.Push(loffset)
	l.SetTableRaw(-3)

	return 1
}

func lchar(l *lua.State) int {
	n := l.AbsIndex(-1)
	xs := make([]rune, 0, n)
	for i := 1; i <= n; i++ {
		xs = append(xs, rune(l.ToInteger(i)))
	}
	l.Push(string(xs))
	return 1
}

func lcodes(l *lua.State) int {
	next := liter(l, l.ToString(1), 1, -1)
	l.Push(func(l *lua.State) int {
		if r, n, k := next(); n <= 0 {
			return 0
		} else if r == utf8.RuneError {
			panic("invalid utf8 code at " + strconv.Itoa(k))
			return 0
		} else {
			l.Push(k)
			l.Push(r)
			return 2
		}
	})
	return 1
}

func lcodepoint(l *lua.State) int {
	i := int(l.OptInteger(2, 1))
	j := int(l.OptInteger(3, int64(i)))
	next := liter(l, l.ToString(1), i, j)
	m := 0
	for {
		if r, n, k := next(); n <= 0 {
			break
		} else if r == utf8.RuneError {
			panic("invalid utf8 code at " + strconv.Itoa(k))
		} else {
			l.Push(r)
			m += 1
		}
	}
	return m
}

func llen(l *lua.State) int {
	i := int(l.OptInteger(2, 1))
	j := int(l.OptInteger(3, -1))
	next := liter(l, l.ToString(1), i, j)
	m := 0
	for {
		if r, n, k := next(); n <= 0 {
			break
		} else if r == utf8.RuneError {
			l.Push(nil)
			l.Push(k)
			return 2
		} else {
			m += 1
		}
	}
	l.Push(m)
	return 1
}

func loffset(l *lua.State) int {
	s := l.ToString(1)
	n := int(l.ToInteger(2))
	d := 0
	if n >= 0 {
		d = 1
	} else {
		d = len(s) + 1
	}
	i := int(l.OptInteger(3, int64(d)))
	if i < 0 {
		i = len(s) + i + 1
	}

	var next func() (rune, int, int)
	var ok func() bool
	r, m, k := rune(0), 0, i+1
	if n > 0 {
		next = liter(l, s, i, -1)
		ok = func() bool {
			n -= 1
			return n >= 0
		}
	} else if n < 0 {
		next = riter(l, s, 1, i-1)
		ok = func() bool {
			n += 1
			return n <= 0
		}
	} else {
		next = riter(l, s, 1, -1)
		ok = func() bool {
			return k > i
		}
	}
	for ok() {
		if r, m, k = next(); m <= 0 {
			return 0
		} else if r == utf8.RuneError {
			panic("invalid utf8 code at " + strconv.Itoa(k))
		}
	}
	l.Push(k)
	return 1
}
