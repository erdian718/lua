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

package lmodtable

import (
	"ofunc/lua"
)

type sorter lua.State

// Len is the number of elements in the collection.
func (self *sorter) Len() int {
	return (*lua.State)(self).Length(1)
}

// Less reports whether the element with index i should sort before the element with index j.
func (self *sorter) Less(i, j int) bool {
	i, j = i+1, j+1
	l := (*lua.State)(self)

	l.PushIndex(2)
	l.Push(i)
	l.GetTable(1)
	l.Push(j)
	l.GetTable(1)
	l.Call(2, 1)
	less := l.ToBoolean(-1)
	l.Pop(1)
	return less
}

// Swap swaps the elements with indexes i and j.
func (self *sorter) Swap(i, j int) {
	i, j = i+1, j+1
	l := (*lua.State)(self)

	l.Push(j)
	l.Push(i)
	l.GetTable(1)
	l.Push(i)
	l.Push(j)
	l.GetTable(1)
	l.SetTable(1)
	l.SetTable(1)
}
