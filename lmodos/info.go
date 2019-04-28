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
