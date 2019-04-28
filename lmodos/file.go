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
