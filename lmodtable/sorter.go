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
