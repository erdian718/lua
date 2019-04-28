package lmodio

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"ofunc/lua"
)

// Open opens the module.
func Open(l *lua.State) int {
	l.NewTable(0, 16)

	l.Push("buffer")
	l.Push(lbuffer)
	l.SetTableRaw(-3)

	l.Push("copy")
	l.Push(lcopy)
	l.SetTableRaw(-3)

	l.Push("limit")
	l.Push(llimit)
	l.SetTableRaw(-3)

	l.Push("read")
	l.Push(lread)
	l.SetTableRaw(-3)

	l.Push("readfile")
	l.Push(lreadfile)
	l.SetTableRaw(-3)

	l.Push("scanner")
	l.Push(lscanner)
	l.SetTableRaw(-3)

	l.Push("type")
	l.Push(ltype)
	l.SetTableRaw(-3)

	l.Push("write")
	l.Push(lwrite)
	l.SetTableRaw(-3)

	l.Push("writefile")
	l.Push(lwritefile)
	l.SetTableRaw(-3)

	l.NewTable(0, 2)

	l.Push("__index")
	l.Push(lindex)
	l.SetTableRaw(-3)

	l.Push("__newindex")
	l.Push(lnewindex)
	l.SetTableRaw(-3)

	l.SetMetaTable(-2)
	return 1
}

func lbuffer(l *lua.State) int {
	l.Push(bytes.NewBufferString(l.OptString(1, "")))
	return 1
}

func lcopy(l *lua.State) int {
	r := toReader(l, 1)
	w := toWriter(l, 2)
	n := l.OptInteger(3, -1)

	var k int64
	var err error
	if n >= 0 {
		k, err = io.CopyN(w, r, n)
	} else {
		k, err = io.Copy(w, r)
	}
	l.Push(k)
	l.Push(errmsg(err))
	return 2
}

func llimit(l *lua.State) int {
	n := l.ToInteger(2)
	l.Push(io.LimitReader(toReader(l, 1), n))
	return 1
}

func lread(l *lua.State) int {
	r := toReader(l, 1)
	n := l.OptInteger(2, -1)

	var err error
	var xs []byte
	if n >= 0 {
		var k int
		xs = make([]byte, n)
		k, err = r.Read(xs)
		xs = xs[:k]
	} else {
		xs, err = ioutil.ReadAll(r)
	}
	l.Push(string(xs))
	l.Push(errmsg(err))
	return 2
}

func lreadfile(l *lua.State) int {
	xs, err := ioutil.ReadFile(l.ToString(1))
	l.Push(string(xs))
	l.Push(errmsg(err))
	return 2
}

func lscanner(l *lua.State) int {
	scanner := bufio.NewScanner(toReader(l, 1))
	var split bufio.SplitFunc
	switch typ := l.TypeOf(2); typ {
	case lua.TypeNil:
		split = bufio.ScanLines
	case lua.TypeString:
		switch name := l.ToString(2); name {
		case "byte":
			split = bufio.ScanBytes
		case "line":
			split = bufio.ScanLines
		case "rune":
			split = bufio.ScanRunes
		case "word":
			split = bufio.ScanWords
		default:
			panic("io.scanner: invalid argument #2: unknown split: " + name)
		}
	default:
		split = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			l.PushIndex(lua.FirstUpVal - 1)
			l.Push(string(data))
			l.Push(atEOF)
			l.Call(2, 3)
			advance = int(l.ToInteger(-3))
			if xs := l.OptString(-2, ""); xs != "" {
				token = []byte(xs)
			}
			if !l.IsNil(-1) {
				err = errors.New(l.ToString(-1))
			}
			l.Pop(3)
			return
		}
	}
	scanner.Split(split)
	l.PushClosure(func(l *lua.State) int {
		if scanner.Scan() {
			l.Push(scanner.Text())
			return 1
		} else {
			l.Push(nil)
			l.Push(errmsg(scanner.Err()))
			return 2
		}
	}, 2)
	return 1
}

func ltype(l *lua.State) int {
	switch l.GetRaw(1).(type) {
	case io.Reader:
		l.Push("reader")
	case io.Writer:
		l.Push("writer")
	case io.ReadWriter:
		l.Push("readwriter")
	default:
		l.Push(nil)
	}
	return 1
}

func lwrite(l *lua.State) int {
	n, err := toWriter(l, 1).Write([]byte(l.OptString(2, "")))
	l.Push(n)
	l.Push(errmsg(err))
	return 2
}

func lwritefile(l *lua.State) int {
	err := ioutil.WriteFile(l.ToString(1), []byte(l.ToString(2)), 0666)
	l.Push(errmsg(err))
	return 1
}

func lindex(l *lua.State) int {
	switch l.ToString(2) {
	case "stderr":
		l.Push(os.Stderr)
	case "stdin":
		l.Push(os.Stdin)
	case "stdout":
		l.Push(os.Stdout)
	default:
		l.Push(nil)
	}
	return 1
}

func lnewindex(l *lua.State) int {
	switch key := l.ToString(2); key {
	case "stderr":
		if f, ok := l.GetRaw(3).(*os.File); ok {
			os.Stderr = f
		} else {
			panic("io.stderr: not a file")
		}
	case "stdin":
		if f, ok := l.GetRaw(3).(*os.File); ok {
			os.Stdin = f
		} else {
			panic("io.stdin: not a file")
		}
	case "stdout":
		if f, ok := l.GetRaw(3).(*os.File); ok {
			os.Stdout = f
		} else {
			panic("io.stdout: not a file")
		}
	default:
		panic("io: invalid field: " + key)
	}
	return 0
}
