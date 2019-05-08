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

package lmodio

import (
	"io"

	"ofunc/lua"
)

func toReader(l *lua.State, i int) io.Reader {
	if r, ok := l.GetRaw(i).(io.Reader); ok {
		return r
	} else {
		panic("io: not a reader: " + l.ToString(i))
	}
}

func toWriter(l *lua.State, i int) io.Writer {
	if w, ok := l.GetRaw(i).(io.Writer); ok {
		return w
	} else {
		panic("io: not a writer: " + l.ToString(i))
	}
}

func errmsg(err error) interface{} {
	if err == nil {
		return nil
	} else if err == io.EOF {
		return "eof"
	} else {
		return err.Error()
	}
}
