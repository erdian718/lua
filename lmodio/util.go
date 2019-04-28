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
