// +build !js !wasm

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

	"github.com/ofunc/lua"
)

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

		fmt.Println("::", path)
		l.ForEachRaw(-1, func() bool {
			if l.TypeOf(-1) == lua.TypeFunction {
				fmt.Println("  -> TESTING", l.ToString(-2))
				if msg := l.PCall(0, 0, true); msg == nil {
					pass += 1
					fmt.Println("     PASS")
				} else {
					fail += 1
					fmt.Println("     FAIL:", msg)
				}
				l.Push(nil)
			}
			return true
		})
		l.Pop(1)
		return nil
	})
	fmt.Printf("=> PASS %v, FAIL %v\n", pass, fail)
	return err
}

// Benchmark runs the specified benchmark script file.
func Benchmark(root string) error {
	return nil // TODO
}
