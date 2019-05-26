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

package lmodbase

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"syscall/js"
)

// OpenSrc opens the Lua src file.
func OpenSrc(p string) (r io.ReadCloser, e error) {
	var wg sync.WaitGroup
	cbok := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resp := args[0]
		if resp.Get("ok").Bool() {
			return resp.Call("text")
		} else {
			e = errors.New(resp.Get("statusText").String())
			return nil
		}
	})
	defer cbok.Release()
	cbreader := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		r = ioutil.NopCloser(strings.NewReader(args[0].String()))
		wg.Done()
		return nil
	})
	defer cbreader.Release()
	cbfail := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e = errors.New(args[0].Get("message").String())
		wg.Done()
		return nil
	})
	defer cbfail.Release()
	wg.Add(1)
	go js.Global().Call("fetch", p).Call("then", cbok).Call("then", cbreader).Call("catch", cbfail)
	wg.Wait()
	return
}
