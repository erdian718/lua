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

// Utility for Lua
package util

import (
	"fmt"

	"github.com/ofunc/lua"
	"github.com/ofunc/lua/lmodbase"
	"github.com/ofunc/lua/lmodos"
)

const (
	// Version number
	Version = lmodbase.Version
)

// The root path for the executable.
var Root string

func init() {
	Root = lmodos.Root
}

// NewState creates a new State, opens the buildin modules, and disables undefined variables.
func NewState() *lua.State {
	l := lua.NewState()
	Strict(l)
	Open(l)
	return l
}

// AddPath adds the searchpath for "require".
func AddPath(path string) {
	lmodbase.Paths = append(lmodbase.Paths, path)
}

// Strict disables undeclared variables.
func Strict(l *lua.State) {
	l.NewTable(0, 2)
	l.Push("__index")
	l.Push(lundefined)
	l.SetTableRaw(-3)
	l.Push("__newindex")
	l.Push(lundefined)
	l.SetTableRaw(-3)
	l.SetMetaTable(lua.GlobalsIndex)
}

// Run runs the specified Lua src file.
func Run(l *lua.State, src string) error {
	r, err := lmodbase.OpenSrc(src)
	if err != nil {
		return err
	}
	defer r.Close()
	if err := l.LoadText(r, src, 0); err != nil {
		return err
	}
	if msg := l.PCall(0, 0, true); msg != nil {
		return fmt.Errorf("%v", msg)
	}
	return nil
}

func lundefined(l *lua.State) int {
	panic("undefined: " + l.ToString(2))
	return 0
}
