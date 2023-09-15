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

// Access to the WebAssembly host environment.
package lmodjs

import (
	"syscall/js"

	"github.com/ofunc/lua"
)

// Open opens the module.
func Open(l *lua.State) int {
	m := jsmeta(l)
	l.Push(RegistryKey)
	l.NewTable(0, 0)
	l.Push("m")
	l.PushIndex(m)
	l.SetTableRaw(-3)
	l.SetTableRaw(lua.RegistryIndex)
	l.NewTable(0, 8)

	l.Push("version")
	l.Push("0.0.1")
	l.SetTableRaw(-3)

	l.Push("global")
	l.Push(global)
	l.PushIndex(m)
	l.SetMetaTable(-2)
	l.SetTableRaw(-3)

	l.Push("type")
	l.Push(lType)
	l.SetTableRaw(-3)

	l.Push("value")
	l.Push(lValue)
	l.SetTableRaw(-3)

	l.Push("new")
	l.Push(lNew)
	l.SetTableRaw(-3)

	l.Push("free")
	l.Push(lFree)
	l.SetTableRaw(-3)

	return 1
}

func lType(l *lua.State) int {
	l.Push(js.ValueOf(l.GetRaw(1)).Type().String())
	return 1
}

func lValue(l *lua.State) int {
	wrap(l, value(l, 1))
	return 1
}

func lNew(l *lua.State) int {
	n := l.AbsIndex(-1)
	args := make([]interface{}, 0, n-1)
	for i := 2; i <= n; i++ {
		args = append(args, value(l, i))
	}
	wrap(l, js.ValueOf(l.GetRaw(1)).New(args...))
	return 1
}

func lFree(l *lua.State) int {
	if v, ok := l.GetRaw(1).(js.Value); ok {
		if x := v.Get(LuaKey); !(x.IsUndefined() || x.IsNull()) {
			id := x.Int()
			l.Push(RegistryKey)
			l.GetTableRaw(lua.RegistryIndex)
			l.Push(id)
			l.GetTableRaw(-2)
			if !l.IsNil(-1) {
				l.Push(2)
				l.GetTableRaw(-2)
				l.GetRaw(-1).(js.Func).Release()
				l.Pop(2)
				l.Push(id)
				l.Push(nil)
				l.SetTableRaw(-3)
			}
		}
	}
	return 0
}
