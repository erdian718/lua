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

package util

import (
	"fmt"
	"os"
	"path/filepath"

	"ofunc/lua"
	"ofunc/lua/lmodbase"
	"ofunc/lua/lmodio"
	"ofunc/lua/lmodmath"
	"ofunc/lua/lmodos"
	"ofunc/lua/lmodstring"
	"ofunc/lua/lmodtable"
	"ofunc/lua/lmodutf8"
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

// Open opens the buildin modules.
func Open(l *lua.State) {
	l.Push(lmodbase.Open)
	l.Call(0, 0)

	l.Preload("string", lmodstring.Open)
	l.Preload("utf8", lmodutf8.Open)
	l.Preload("table", lmodtable.Open)
	l.Preload("math", lmodmath.Open)
	l.Preload("io", lmodio.Open)
	l.Preload("os", lmodos.Open)
}

// Run runs the specified script file.
func Run(l *lua.State, script string) error {
	f, err := os.Open(script)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := l.LoadText(f, f.Name(), 0); err != nil {
		return err
	}
	if msg := l.PCall(0, 0, true); msg != nil {
		return fmt.Errorf("%v", msg)
	}
	return nil
}

// Test runs the specified test script file.
func Test(l *lua.State, root string) error {
	pass, fail := 0, 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(info.Name()) != ".lua" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := l.LoadText(f, path, 0); err != nil {
			return err
		}
		if msg := l.PCall(0, 1, true); msg != nil {
			return fmt.Errorf("%v", msg)
		}
		if l.TypeOf(-1) != lua.TypeTable {
			return nil
		}

		fmt.Println("****", path)
		l.ForEachRaw(-1, func() bool {
			if l.TypeOf(-1) == lua.TypeFunction {
				fmt.Println("---- TEST", l.ToString(-2))
				if msg := l.PCall(0, 0, true); msg == nil {
					pass += 1
					fmt.Println("==== PASS")
				} else {
					fail += 1
					fmt.Println("==== FAIL:", msg)
				}
				l.Push(nil)
			}
			return true
		})
		l.Pop(1)
		return nil
	})
	fmt.Printf(">>>> PASS %v, FAIL %v\n", pass, fail)
	return err
}

// Benchmark runs the specified benchmark script file.
func Benchmark(root string) error {
	return nil // TODO
}

func lundefined(l *lua.State) int {
	panic("undefined: " + l.ToString(2))
	return 0
}
