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

package lmodtable

import (
	"sort"
	"strings"

	"ofunc/lua"
)

// Open opens the module.
func Open(l *lua.State) int {
	l.NewTable(0, 8)

	l.Push("concat")
	l.Push(lconcat)
	l.SetTableRaw(-3)

	l.Push("insert")
	l.Push(linsert)
	l.SetTableRaw(-3)

	l.Push("move")
	l.Push(lmove)
	l.SetTableRaw(-3)

	l.Push("pack")
	l.Push(lpack)
	l.SetTableRaw(-3)

	l.Push("remove")
	l.Push(lremove)
	l.SetTableRaw(-3)

	l.Push("sort")
	l.Push(lsort)
	l.SetTableRaw(-3)

	l.Push("unpack")
	l.Push(lunpack)
	l.SetTableRaw(-3)

	return 1
}

func lconcat(l *lua.State) int {
	n := l.Length(1)
	s := l.OptString(2, "")
	i := l.OptInteger(3, 1)
	j := l.OptInteger(4, int64(n))

	xs := make([]string, 0, n)
	for k := i; k <= j; k++ {
		l.Push(k)
		l.GetTable(1)
		xs = append(xs, l.ToString(-1))
		l.Pop(1)
	}
	l.Push(strings.Join(xs, s))
	return 1
}

func linsert(l *lua.State) int {
	n := int64(l.Length(1))
	p := n + 1
	v := 2

	if l.AbsIndex(-1) > 2 {
		p = l.ToInteger(2)
		v = 3
	}
	for i := n; i >= p; i-- {
		l.Push(i + 1)
		l.Push(i)
		l.GetTable(1)
		l.SetTable(1)
	}
	l.Push(p)
	l.PushIndex(v)
	l.SetTable(1)
	return 0
}

func lmove(l *lua.State) int {
	f := l.ToInteger(2)
	e := l.ToInteger(3)
	t := l.ToInteger(4)
	a := 5
	if l.IsNil(5) {
		a = 1
	}

	if e >= f {
		n := e - f
		if t > e || t <= f || a != 1 {
			for i := int64(0); i <= n; i++ {
				l.Push(t + i)
				l.Push(f + i)
				l.GetTable(1)
				l.SetTable(a)
			}
		} else {
			for i := n; i >= 0; i-- {
				l.Push(t + i)
				l.Push(f + i)
				l.GetTable(1)
				l.SetTable(a)
			}
		}
	}
	l.PushIndex(a)
	return 1
}

func lpack(l *lua.State) int {
	n := l.AbsIndex(-1)
	l.NewTable(n, 2)
	idx := l.AbsIndex(-1)
	l.Push("n")
	l.Push(n)
	l.SetTableRaw(idx)
	for i := 1; i <= n; i++ {
		l.Push(i)
		l.PushIndex(i)
		l.SetTableRaw(idx)
	}
	return 1
}

func lremove(l *lua.State) int {
	n := int64(l.Length(1))
	p := l.OptInteger(2, n)

	if p != n && (p < 1 || p > n) {
		panic("index out of range")
	}
	l.Push(p)
	l.GetTable(1)

	for i := p; i < n; i++ {
		l.Push(i)
		l.Push(i + 1)
		l.GetTable(1)
		l.SetTable(1)
	}
	l.Push(n)
	l.Push(nil)
	l.SetTable(1)
	return 1
}

func lsort(l *lua.State) int {
	if l.AbsIndex(-1) == 1 {
		l.Push(func(l *lua.State) int {
			l.Push(l.Compare(1, 2, lua.OpLessThan))
			return 1
		})
	}
	sort.Sort((*sorter)(l))
	return 0
}

func lunpack(l *lua.State) int {
	i := l.OptInteger(2, 1)
	j := l.OptInteger(3, int64(l.Length(1)))

	if i > j {
		return 0
	}
	n := int(j - i + 1)
	for ; i <= j; i++ {
		l.Push(i)
		l.GetTable(1)
	}
	return n
}
