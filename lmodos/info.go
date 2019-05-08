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

package lmodos

import (
	"os"

	"ofunc/lua"
)

func metainfo(l *lua.State) int {
	l.NewTable(0, 2)
	idx := l.AbsIndex(-1)

	l.Push("__index")
	l.Push(func(l *lua.State) int {
		info, ok := l.GetRaw(1).(os.FileInfo)
		if !ok {
			panic("not a fileinfo: " + l.ToString(1))
		}
		switch l.ToString(2) {
		case "name":
			l.Push(info.Name())
		case "isdir":
			l.Push(info.IsDir())
		case "size":
			l.Push(info.Size())
		case "modtime":
			l.Push(info.ModTime().Unix())
		default:
			l.Push(nil)
		}
		return 1
	})
	l.SetTableRaw(idx)
	return idx
}
