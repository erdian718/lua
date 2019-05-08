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

func matefile(l *lua.State) int {
	l.NewTable(0, 2)
	idx := l.AbsIndex(-1)

	l.Push("__index")
	l.PushIndex(idx)
	l.SetTableRaw(idx)

	l.Push("close")
	l.Push(lclose)
	l.SetTableRaw(idx)

	l.Push("seek")
	l.Push(lseek)
	l.SetTableRaw(idx)

	return idx
}

func lclose(l *lua.State) int {
	if err := toFile(l, 1).Close(); err == nil {
		l.Push(true)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func lseek(l *lua.State) int {
	whence := os.SEEK_CUR
	switch name := l.OptString(2, "cur"); name {
	case "set":
		whence = os.SEEK_SET
	case "cur":
		whence = os.SEEK_CUR
	case "end":
		whence = os.SEEK_END
	default:
		l.Push(nil)
		l.Push("invalid whence: " + name)
		return 2
	}
	offset := l.OptInteger(3, 0)
	if ret, err := toFile(l, 1).Seek(offset, whence); err == nil {
		l.Push(ret)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func toFile(l *lua.State, i int) *os.File {
	if f, ok := l.GetRaw(1).(*os.File); ok {
		return f
	} else {
		panic("not a file: " + l.ToString(i))
		return nil
	}
}
