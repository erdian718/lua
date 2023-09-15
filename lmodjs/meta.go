// +build js,wasm

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

package lmodjs

import (
	"syscall/js"

	"github.com/ofunc/lua"
)

func jsmeta(l *lua.State) int {
	l.NewTable(0, 8)
	idx := l.AbsIndex(-1)

	l.Push("__len")
	l.Push(lLen)
	l.SetTableRaw(idx)

	l.Push("__index")
	l.Push(lIndex)
	l.SetTableRaw(idx)

	l.Push("__newindex")
	l.Push(lNewIndex)
	l.SetTableRaw(idx)

	l.Push("__pairs")
	l.Push(lPairs)
	l.SetTableRaw(idx)

	l.Push("__call")
	l.Push(lCall)
	l.SetTableRaw(idx)

	return idx
}

func lLen(l *lua.State) int {
	l.Push(js.ValueOf(l.GetRaw(1)).Length())
	return 1
}

func lIndex(l *lua.State) int {
	v := js.ValueOf(l.GetRaw(1))
	if i, err := l.TryInteger(2); err == nil && int(i) > 0 {
		wrap(l, v.Index(int(i)-1))
	} else {
		wrap(l, v.Get(l.ToString(2)))
	}
	return 1
}

func lNewIndex(l *lua.State) int {
	v := js.ValueOf(l.GetRaw(1))
	if i, err := l.TryInteger(2); err == nil && int(i) > 0 {
		v.SetIndex(int(i)-1, value(l, 3))
	} else {
		v.Set(l.ToString(2), value(l, 3))
	}
	return 0
}

func lPairs(l *lua.State) int {
	iter := object.Call("entries", js.ValueOf(l.GetRaw(1))).Call("entries")
	l.Push(func(l *lua.State) int {
		x := iter.Call("next")
		if x.Get("done").Bool() {
			return 0
		}
		y := x.Get("value").Index(1)
		wrap(l, y.Index(0))
		wrap(l, y.Index(1))
		return 2
	})
	return 1
}

func lCall(l *lua.State) int {
	v := js.ValueOf(l.GetRaw(1))
	n := l.AbsIndex(-1)
	args := make([]interface{}, 0, n-1)
	for i := 2; i <= n; i++ {
		args = append(args, value(l, i))
	}
	wrap(l, v.Call("call", args...))
	return 1
}
