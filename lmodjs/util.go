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

	"ofunc/lua"
)

var global = js.Global()
var undefined = js.Undefined()
var object = global.Get("Object")
var array = global.Get("Array")
var regid int64

func wrap(l *lua.State, v js.Value) {
	switch v.Type() {
	case js.TypeUndefined, js.TypeNull:
		l.Push(nil)
	case js.TypeBoolean:
		l.Push(v.Bool())
	case js.TypeNumber:
		l.Push(v.Float())
	case js.TypeString:
		l.Push(v.String())
	default:
		l.Push(v)
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
	}
}

func value(l *lua.State, idx int) js.Value {
	idx = l.AbsIndex(idx)
	typ := l.TypeOf(idx)
	switch typ {
	case lua.TypeBoolean:
		return js.ValueOf(l.ToBoolean(idx))
	case lua.TypeNumber:
		return js.ValueOf(l.ToFloat(idx))
	case lua.TypeString:
		return js.ValueOf(l.ToString(idx))
	case lua.TypeFunction:
		return vfunction(l, idx)
	case lua.TypeUserData:
		if v, ok := l.GetRaw(idx).(js.Value); ok {
			return v
		}
	}
	if l.GetMetaField(idx, "__call") != lua.TypeNil {
		l.Pop(1)
		return vfunction(l, idx)
	}
	if l.GetMetaField(idx, "__pairs") != lua.TypeNil {
		l.Pop(1)
		return vobject(l, idx)
	}
	if l.GetMetaField(idx, "__len") != lua.TypeNil {
		l.Pop(1)
		return varray(l, idx)
	}
	if typ == lua.TypeTable {
		n := l.Count(idx)
		if l.LengthRaw(idx) == n && n > 0 {
			return varray(l, idx)
		} else {
			return vobject(l, idx)
		}
	}
	return undefined
}

func vfunction(l *lua.State, idx int) js.Value {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		wrap(l, this) // TODO metatable
		for _, arg := range args {
			wrap(l, arg)
		}
		l.Call(len(args)+1, 1)
		return value(l, -1)
	}).Value
}

func vobject(l *lua.State, idx int) js.Value {
	v := object.New()
	l.ForEach(idx, func() bool {
		v.Set(l.ToString(-2), value(l, -1))
		return true
	})
	return v
}

func varray(l *lua.State, idx int) js.Value {
	n := l.Length(idx)
	v := array.New(n)
	for i := 0; i < n; i++ {
		l.Push(i + 1)
		l.GetTable(idx)
		v.SetIndex(i, value(l, -1))
		l.Pop(1)
	}
	return v
}
