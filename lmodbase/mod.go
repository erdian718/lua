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

package lmodbase

import (
	"fmt"
	"os"
	"path/filepath"

	"ofunc/lua"
)

const (
	// Version number
	Version = "1.0.0"
)

// The searchpaths for "require".
var Paths []string

// Open opens the module.
func Open(l *lua.State) int {
	l.Push("_LOADED")
	if l.GetTableRaw(lua.RegistryIndex) != lua.TypeTable {
		l.Push("_LOADED")
		l.NewTable(0, 64)
		l.SetTableRaw(lua.RegistryIndex)
	}
	l.Push("_PRELOAD")
	if l.GetTableRaw(lua.RegistryIndex) != lua.TypeTable {
		l.Push("_PRELOAD")
		l.NewTable(0, 16)
		l.SetTableRaw(lua.RegistryIndex)
	}

	l.Push("_VERSION")
	l.Push(Version)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("require")
	l.Push(lrequire)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("assert")
	l.Push(lassert)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("error")
	l.Push(lerror)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("getmetatable")
	l.Push(lgetmetatable)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("ipairs")
	l.Push(lipairs)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("pairs")
	l.Push(lpairs)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("pcall")
	l.Push(lpcall)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("print")
	l.Push(lprint)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("rawequal")
	l.Push(lrawequal)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("rawget")
	l.Push(lrawget)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("rawlen")
	l.Push(lrawlen)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("rawset")
	l.Push(lrawset)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("select")
	l.Push(lselect)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("setmetatable")
	l.Push(lsetmetatable)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("tonumber")
	l.Push(ltonumber)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("tostring")
	l.Push(ltostring)
	l.SetTableRaw(lua.GlobalsIndex)

	l.Push("type")
	l.Push(ltype)
	l.SetTableRaw(lua.GlobalsIndex)

	return 0
}

func lrequire(l *lua.State) int {
	l.Push("_LOADED")
	l.GetTableRaw(lua.RegistryIndex)
	loaded := l.AbsIndex(-1)
	l.PushIndex(1)
	if l.GetTableRaw(loaded) != lua.TypeNil {
		return 1
	}

	l.Push("_PRELOAD")
	l.GetTableRaw(lua.RegistryIndex)
	l.PushIndex(1)
	if l.GetTableRaw(-2) == lua.TypeFunction {
		l.PushIndex(1)
		l.Call(1, 1)
		if l.IsNil(-1) {
			l.Push(true)
		} else {
			l.PushIndex(-1)
		}
		l.PushIndex(1)
		l.PushIndex(-2)
		l.SetTableRaw(loaded)
		return 1
	}

	name := l.ToString(1)
	for _, p := range Paths {
		p = filepath.Join(p, name)
		f, err := os.Open(p + ".lua")
		if err != nil {
			f, err = os.Open(p + "/init.lua")
			if err != nil {
				continue
			}
		}
		if err := l.LoadText(f, name, 0); err != nil {
			panic(err)
			return 0
		}
		l.Call(0, 1)
		f.Close()
		if l.IsNil(-1) {
			l.Push(true)
		} else {
			l.PushIndex(-1)
		}
		l.PushIndex(1)
		l.PushIndex(-2)
		l.SetTableRaw(loaded)
		return 1
	}
	panic("require: module '" + name + "' not found")
}

func lassert(l *lua.State) int {
	if l.ToBoolean(1) {
		return l.AbsIndex(-1)
	} else {
		panic(l.OptString(2, "assertion failed!"))
	}
}

func lerror(l *lua.State) int {
	l.Error()
	return 0
}

func lgetmetatable(l *lua.State) int {
	if l.GetMetaField(1, "__metatable") == lua.TypeNil {
		if l.GetMetaTable(1) {
			return 1
		} else {
			return 0
		}
	} else {
		return 1
	}
}

func lipairs(l *lua.State) int {
	l.Push(func(l *lua.State) int {
		i := l.ToInteger(2) + 1
		l.Push(i)
		l.Push(i)
		if l.GetTable(1) == lua.TypeNil {
			return 1
		} else {
			return 2
		}
	})
	l.PushIndex(1)
	l.Push(0)
	return 3
}

func lpairs(l *lua.State) int {
	if l.GetMetaField(1, "__pairs") == lua.TypeNil {
		l.GetIter(1)
		l.PushIndex(1)
		l.Push(nil)
	} else {
		l.PushIndex(1)
		l.Call(1, 3)
	}
	return 3
}

func lpcall(l *lua.State) int {
	if msg := l.PCall(l.AbsIndex(-1)-1, -1, false); msg == nil {
		l.Push(true)
		n := l.AbsIndex(-1)
		if n > 1 {
			l.Insert(1)
		}
		return n
	} else {
		l.Push(false)
		if err, ok := msg.(error); ok {
			l.Push(err.Error())
		} else {
			l.Push(msg)
		}
		return 2
	}
}

func lprint(l *lua.State) int {
	n := l.AbsIndex(-1)
	for i := 1; i <= n; i++ {
		fmt.Print(l.ToString(i))
		if i < n {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

func lrawequal(l *lua.State) int {
	l.Push(l.CompareRaw(1, 2, lua.OpEqual))
	return 1
}

func lrawget(l *lua.State) int {
	l.PushIndex(2)
	l.GetTableRaw(1)
	return 1
}

func lrawlen(l *lua.State) int {
	l.Push(l.LengthRaw(1))
	return 1
}

func lrawset(l *lua.State) int {
	l.PushIndex(2)
	l.PushIndex(3)
	l.SetTableRaw(1)
	l.PushIndex(1)
	return 1
}

func lselect(l *lua.State) int {
	n := l.AbsIndex(-1)
	if l.TypeOf(1) == lua.TypeString && l.ToString(1) == "#" {
		l.Push(n - 1)
		return 1
	} else {
		if i := int(l.ToInteger(1)); i > 0 {
			return n - i
		} else if i < 0 && i > -n {
			return -i
		} else {
			panic("select: index out of range")
		}
	}
}

func lsetmetatable(l *lua.State) int {
	if l.GetMetaField(1, "__metatable") == lua.TypeNil {
		l.PushIndex(2)
		l.SetMetaTable(1)
		l.PushIndex(1)
		return 1
	} else {
		panic("setmetatable: cannot change a protected metatable")
	}
}

func ltonumber(l *lua.State) int {
	if v, err := l.TryInteger(1); err == nil {
		l.Push(v)
	} else if v, err := l.TryFloat(1); err == nil {
		l.Push(v)
	} else {
		l.Push(nil)
	}
	return 1
}

func ltostring(l *lua.State) int {
	l.Push(l.ToString(1))
	return 1
}

func ltype(l *lua.State) int {
	l.Push(l.TypeOf(1).String())
	return 1
}
