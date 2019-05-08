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

package lmodstring

import (
	"fmt"
	"strings"

	"ofunc/lua"
	"ofunc/lua/lmodstring/pm"
)

// Open opens the module.
func Open(l *lua.State) int {
	l.NewTable(0, 16)

	l.Push("byte")
	l.Push(lbyte)
	l.SetTableRaw(-3)

	l.Push("char")
	l.Push(lchar)
	l.SetTableRaw(-3)

	l.Push("find")
	l.Push(lfind)
	l.SetTableRaw(-3)

	l.Push("format")
	l.Push(lformat)
	l.SetTableRaw(-3)

	l.Push("gmatch")
	l.Push(lgmatch)
	l.SetTableRaw(-3)

	l.Push("gsub")
	l.Push(lgsub)
	l.SetTableRaw(-3)

	l.Push("len")
	l.Push(llen)
	l.SetTableRaw(-3)

	l.Push("lower")
	l.Push(llower)
	l.SetTableRaw(-3)

	l.Push("match")
	l.Push(lmatch)
	l.SetTableRaw(-3)

	l.Push("rep")
	l.Push(lrep)
	l.SetTableRaw(-3)

	l.Push("reverse")
	l.Push(lreverse)
	l.SetTableRaw(-3)

	l.Push("sub")
	l.Push(lsub)
	l.SetTableRaw(-3)

	l.Push("upper")
	l.Push(lupper)
	l.SetTableRaw(-3)

	l.Push("__index")
	l.PushIndex(-2)
	l.SetTableRaw(-3)

	l.Push("")
	l.PushIndex(-2)
	l.SetMetaTable(-2)
	l.Pop(1)
	return 1
}

func lbyte(l *lua.State) int {
	s := l.OptString(1, "")
	n := len(s)
	i := int(l.OptInteger(2, 1))
	j := int(l.OptInteger(3, int64(i)))
	if i < 0 {
		i = n + i + 1
	}
	if j < 0 {
		j = n + j + 1
	}
	if i < 1 {
		i = 1
	}
	if j > n {
		j = n
	}

	m := 0
	for k := i - 1; k < j; k++ {
		l.Push(int64(s[k]))
		m += 1
	}
	return m
}

func lchar(l *lua.State) int {
	n := l.AbsIndex(-1)
	xs := make([]byte, 0, n)
	for i := 1; i <= n; i++ {
		xs = append(xs, byte(l.ToInteger(i)))
	}
	l.Push(string(xs))
	return 1
}

func lfind(l *lua.State) int {
	s := l.OptString(1, "")
	n := len(s)
	pattern := l.OptString(2, "")
	init := int(l.OptInteger(3, 1))
	if init < 0 {
		init = n + init + 1
	}
	if init < 1 {
		init = 1
	}
	init -= 1

	if l.ToBoolean(4) {
		if pos := strings.Index(s[init:], pattern); pos < 0 {
			l.Push(nil)
			return 1
		} else {
			l.Push(init + pos + 1)
			l.Push(init + pos + len(pattern))
			return 2
		}
	}

	mds, err := pm.Find(pattern, []byte(s), init, 1)
	if err != nil {
		panic(err)
	}
	if len(mds) == 0 {
		l.Push(nil)
		return 1
	}
	md := mds[0]
	l.Push(md.Capture(0) + 1)
	l.Push(md.Capture(1))
	for i := 2; i < md.CaptureLength(); i += 2 {
		if md.IsPosCapture(i) {
			l.Push(md.Capture(i))
		} else {
			l.Push(s[md.Capture(i):md.Capture(i+1)])
		}
	}
	return md.CaptureLength()/2 + 1
}

func lformat(l *lua.State) int {
	n := l.AbsIndex(-1)
	xs := make([]interface{}, 0, n-1)
	for i := 2; i <= n; i++ {
		xs = append(xs, l.GetRaw(i))
	}
	l.Push(fmt.Sprintf(l.OptString(1, ""), xs...))
	return 1
}

func lgmatch(l *lua.State) int {
	s := l.OptString(1, "")
	pattern := l.OptString(2, "")
	mds, err := pm.Find(pattern, []byte(s), 0, -1)
	if err != nil {
		panic(err)
	}

	idx, n := 0, len(mds)
	l.Push(func(l *lua.State) int {
		if idx >= n {
			return 0
		}
		md := mds[idx]
		idx += 1
		if md.CaptureLength() == 2 {
			l.Push(s[md.Capture(0):md.Capture(1)])
			return 1
		}
		for i := 2; i < md.CaptureLength(); i += 2 {
			if md.IsPosCapture(i) {
				l.Push(md.Capture(i))
			} else {
				l.Push(s[md.Capture(i):md.Capture(i+1)])
			}
		}
		return md.CaptureLength()/2 - 1
	})
	return 1
}

func lgsub(l *lua.State) int {
	s := l.OptString(1, "")
	pattern := l.OptString(2, "")
	limit := int(l.OptInteger(4, -1))
	mds, err := pm.Find(pattern, []byte(s), 0, limit)
	if err != nil {
		panic(err)
	}
	if len(mds) == 0 {
		l.PushIndex(1)
		l.Push(0)
		return 2
	}
	switch l.TypeOf(3) {
	case lua.TypeString:
		l.Push(gsubstr(l, s, mds))
	case lua.TypeTable:
		l.Push(gsubtable(l, s, mds))
	case lua.TypeFunction:
		l.Push(gsubfunction(l, s, mds))
	default:
		panic("string.gsub: invalid argument #3: " + l.ToString(3))
	}
	l.Push(len(mds))
	return 2
}

func llen(l *lua.State) int {
	l.Push(l.Length(1))
	return 1
}

func llower(l *lua.State) int {
	l.Push(strings.ToLower(l.OptString(1, "")))
	return 1
}

func lmatch(l *lua.State) int {
	s := l.OptString(1, "")
	n := len(s)
	pattern := l.OptString(2, "")
	offset := int(l.OptInteger(3, 1))
	if offset < 0 {
		offset = offset + n + 1
	}
	if offset < 1 {
		offset = 1
	}
	offset -= 1

	mds, err := pm.Find(pattern, []byte(s), offset, 1)
	if err != nil {
		panic(err)
	}
	if len(mds) == 0 {
		return 0
	}
	md := mds[0]
	switch nsubs := md.CaptureLength() / 2; nsubs {
	case 1:
		l.Push(s[md.Capture(0):md.Capture(1)])
		return 1
	default:
		for i := 2; i < md.CaptureLength(); i += 2 {
			if md.IsPosCapture(i) {
				l.Push(md.Capture(i))
			} else {
				l.Push(s[md.Capture(i):md.Capture(i+1)])
			}
		}
		return nsubs - 1
	}
}

func lrep(l *lua.State) int {
	s := l.OptString(1, "")
	n := int(l.OptInteger(2, 1))
	sep := l.OptString(3, "")
	xs := make([]string, n)
	for i := 0; i < n; i++ {
		xs[i] = s
	}
	l.Push(strings.Join(xs, sep))
	return 1
}

func lreverse(l *lua.State) int {
	s := l.OptString(1, "")
	n := len(s)
	xs := make([]byte, 0, n)
	for i := n - 1; i >= 0; i-- {
		xs = append(xs, s[i])
	}
	l.Push(string(xs))
	return 1
}

func lsub(l *lua.State) int {
	s := l.OptString(1, "")
	n := len(s)
	i := int(l.OptInteger(2, 1))
	j := int(l.OptInteger(3, -1))
	if i < 0 {
		i = n + i + 1
	}
	if j < 0 {
		j = n + j + 1
	}
	if i < 1 {
		i = 1
	}
	if j > n {
		j = n
	}
	if i > j {
		return 0
	}
	l.Push(s[i-1 : j])
	return 1
}

func lupper(l *lua.State) int {
	l.Push(strings.ToUpper(l.OptString(1, "")))
	return 1
}
