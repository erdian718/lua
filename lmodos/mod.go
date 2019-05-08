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
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ofunc/lua"
)

// The root path for the executable.
var Root string
var start = time.Now()

func init() {
	epath, _ := os.Executable()
	Root = filepath.Dir(epath)
}

// Open opens the module.
func Open(l *lua.State) int {
	minfo := metainfo(l)
	mfile := matefile(l)
	mprog := mateprog(l)
	l.NewTable(0, 32)

	l.Push("abs")
	l.Push(labs)
	l.SetTableRaw(-3)

	l.Push("clock")
	l.Push(lclock)
	l.SetTableRaw(-3)

	l.Push("date")
	l.Push(ldate)
	l.SetTableRaw(-3)

	l.Push("difftime")
	l.Push(ldifftime)
	l.SetTableRaw(-3)

	l.Push("execute")
	l.Push(lexecute)
	l.SetTableRaw(-3)

	l.Push("exists")
	l.Push(lexists)
	l.SetTableRaw(-3)

	l.Push("exit")
	l.Push(lexit)
	l.SetTableRaw(-3)

	l.Push("getenv")
	l.Push(lgetenv)
	l.SetTableRaw(-3)

	l.Push("join")
	l.Push(ljoin)
	l.SetTableRaw(-3)

	l.Push("mkdir")
	l.Push(lmkdir)
	l.SetTableRaw(-3)

	l.Push("open")
	l.PushClosure(lopen, mfile)
	l.SetTableRaw(-3)

	l.Push("popen")
	l.PushClosure(lpopen, mprog)
	l.SetTableRaw(-3)

	l.Push("remove")
	l.Push(lremove)
	l.SetTableRaw(-3)

	l.Push("rename")
	l.Push(lrename)
	l.SetTableRaw(-3)

	l.Push("root")
	l.Push(Root)
	l.SetTableRaw(-3)

	l.Push("setlocale")
	l.Push(lsetlocale)
	l.SetTableRaw(-3)

	l.Push("stat")
	l.PushClosure(lstat, minfo)
	l.SetTableRaw(-3)

	l.Push("time")
	l.Push(ltime)
	l.SetTableRaw(-3)

	l.Push("tmpfile")
	l.PushClosure(ltmpfile, mfile)
	l.SetTableRaw(-3)

	l.Push("tmpname")
	l.Push(ltmpname)
	l.SetTableRaw(-3)

	l.Push("walk")
	l.PushClosure(lwalk, minfo)
	l.SetTableRaw(-3)

	return 1
}

func labs(l *lua.State) int {
	if p, err := filepath.Abs(l.ToString(1)); err == nil {
		l.Push(p)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func lclock(l *lua.State) int {
	l.Push(float64(time.Now().Sub(start)) / float64(time.Second))
	return 1
}

func ldate(l *lua.State) int {
	f := l.OptString(1, "%c")
	t := time.Unix(l.OptInteger(2, time.Now().Unix()), 0)
	if strings.HasPrefix(f, "!") {
		t = t.UTC()
		f = f[1:]
	}
	if strings.HasPrefix(f, "*t") {
		l.NewTable(0, 8)
		idx := l.AbsIndex(-1)

		l.Push("year")
		l.Push(t.Year())
		l.SetTableRaw(idx)

		l.Push("month")
		l.Push(int64(t.Month()))
		l.SetTableRaw(idx)

		l.Push("day")
		l.Push(t.Day())
		l.SetTableRaw(idx)

		l.Push("hour")
		l.Push(t.Hour())
		l.SetTableRaw(idx)

		l.Push("min")
		l.Push(t.Minute())
		l.SetTableRaw(idx)

		l.Push("sec")
		l.Push(t.Second())
		l.SetTableRaw(idx)

		l.Push("wday")
		l.Push(int64(t.Weekday()) + 1)
		l.SetTableRaw(idx)

		l.Push("yday")
		l.Push(t.YearDay())
		l.SetTableRaw(idx)
	} else {
		l.Push(strftime(t, f))
	}
	return 1
}

func ldifftime(l *lua.State) int {
	x := l.ToInteger(1)
	y := l.ToInteger(2)
	l.Push(x - y)
	return 1
}

func lexecute(l *lua.State) int {
	n := l.AbsIndex(-1)
	if n < 1 {
		l.Push(true)
		return 1
	}

	args := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		args = append(args, l.ToString(i))
	}
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err == nil {
		l.Push(true)
		l.Push("exit")
		l.Push(0)
	} else if e, ok := err.(*exec.ExitError); ok {
		l.Push(nil)
		if e.Exited() {
			l.Push("exit")
		} else {
			l.Push("signal")
		}

		if status, ok := e.Sys().(syscall.WaitStatus); ok {
			switch {
			case status.Exited():
				l.Push(int64(status.ExitCode))
			case status.Signaled():
				l.Push(status.Signal().String())
			case status.Stopped():
				l.Push(status.StopSignal().String())
			default:
				l.Push(nil)
			}
		}
	} else {
		l.Push(nil)
		l.Push(nil)
		l.Push(err.Error())
	}
	return 3
}

func lexists(l *lua.State) int {
	if _, err := os.Stat(l.ToString(1)); err == nil {
		l.Push(true)
		return 1
	} else if os.IsNotExist(err) {
		l.Push(false)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func lexit(l *lua.State) int {
	if l.TypeOf(1) == lua.TypeNumber {
		os.Exit(int(l.ToInteger(1)))
	} else if l.IsNil(1) || l.ToBoolean(1) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
	return 0
}

func lgetenv(l *lua.State) int {
	if v := os.Getenv(l.ToString(1)); v == "" {
		l.Push(nil)
	} else {
		l.Push(v)
	}
	return 1
}

func ljoin(l *lua.State) int {
	n := l.AbsIndex(-1)
	xs := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		xs = append(xs, l.ToString(i))
	}
	l.Push(filepath.Join(xs...))
	return 1
}

func lmkdir(l *lua.State) int {
	mk := os.Mkdir
	if l.ToBoolean(2) {
		mk = os.MkdirAll
	}
	if err := mk(l.ToString(1), os.ModeDir); err == nil {
		l.Push(true)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func lopen(l *lua.State) int {
	flag := 0
	mode := l.OptString(2, "r")
	if strings.HasSuffix(mode, "b") {
		mode = strings.TrimSuffix(mode, "b")
	}
	switch mode {
	case "r":
		flag = os.O_RDONLY
	case "w":
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case "a":
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case "r+":
		flag = os.O_RDWR
	case "w+":
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	case "a+":
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	default:
		l.Push(nil)
		l.Push("os.open: unknown mode: " + mode)
		return 2
	}
	if f, err := os.OpenFile(l.ToString(1), flag, 0666); err == nil {
		l.Push(f)
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func lpopen(l *lua.State) int {
	n := l.AbsIndex(-1)
	name := l.ToString(1)
	mode := l.OptString(2, "r")
	args := make([]string, 0, n-1)
	for i := 3; i <= n; i++ {
		args = append(args, l.ToString(i))
	}

	p := prog{}
	p.cmd = exec.Command(name, args...)
	ok := true
	switch mode {
	case "r":
		if r, err := p.cmd.StdoutPipe(); err == nil {
			p.r = r
			l.Push(p)
			l.PushIndex(lua.FirstUpVal - 1)
			l.SetMetaTable(-2)
			l.Push(nil)
		} else {
			ok = false
			l.Push(nil)
			l.Push(err.Error())
		}
	case "w":
		if w, err := p.cmd.StdinPipe(); err == nil {
			p.w = w
			l.Push(p)
			l.PushIndex(lua.FirstUpVal - 1)
			l.SetMetaTable(-2)
			l.Push(nil)
		} else {
			ok = false
			l.Push(nil)
			l.Push(err.Error())
		}
	default:
		ok = false
		l.Push(nil)
		l.Push("unknown popen mode: " + mode)
	}
	if ok {
		if err := p.cmd.Start(); err != nil {
			l.Push(nil)
			l.Push(err.Error())
		}
	}
	return 2
}

func lremove(l *lua.State) int {
	rm := os.Remove
	if l.ToBoolean(2) {
		rm = os.RemoveAll
	}
	if err := rm(l.ToString(1)); err == nil {
		l.Push(true)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func lrename(l *lua.State) int {
	old := l.ToString(1)
	new := l.ToString(2)
	if err := os.Rename(old, new); err == nil {
		l.Push(true)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func lsetlocale(l *lua.State) int {
	return 0
}

func lstat(l *lua.State) int {
	var info os.FileInfo
	var err error
	switch l.TypeOf(1) {
	case lua.TypeString:
		info, err = os.Stat(l.ToString(1))
	case lua.TypeUserData:
		if f, ok := l.GetRaw(1).(*os.File); ok {
			info, err = f.Stat()
		} else {
			panic("not string or fileinfo: " + l.ToString(1))
		}
	default:
		panic("not string or fileinfo: " + l.ToString(1))
	}

	if err == nil {
		l.Push(info)
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func ltime(l *lua.State) int {
	if l.IsNil(1) {
		l.Push(time.Now().Unix())
	} else {
		l.Push("year")
		l.GetTable(1)
		year := int(l.ToInteger(-1))

		l.Push("month")
		l.GetTable(1)
		month := time.Month(l.ToInteger(-1))

		l.Push("day")
		l.GetTable(1)
		day := int(l.ToInteger(-1))

		l.Push("hour")
		l.GetTable(1)
		hour := int(l.OptInteger(-1, 12))

		l.Push("min")
		l.GetTable(1)
		min := int(l.OptInteger(-1, 0))

		l.Push("sec")
		l.GetTable(1)
		sec := int(l.OptInteger(-1, 0))

		l.Push(time.Date(year, month, day, hour, min, sec, 0, time.Local).Unix())
	}
	return 1
}

func ltmpfile(l *lua.State) int {
	if f, err := ioutil.TempFile("", ""); err == nil {
		l.Push(f)
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		return 1
	} else {
		l.Push(nil)
		l.Push(err.Error())
		return 2
	}
}

func ltmpname(l *lua.State) int {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return 0
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	l.Push(name)
	return 1
}

func lwalk(l *lua.State) int {
	err := filepath.Walk(l.ToString(1), func(path string, info os.FileInfo, err error) error {
		l.PushIndex(2)
		l.Push(path)
		l.Push(info)
		l.PushIndex(lua.FirstUpVal - 1)
		l.SetMetaTable(-2)
		if err == nil {
			l.Push(nil)
		} else {
			l.Push(err.Error())
		}
		l.Call(3, 1)
		switch l.TypeOf(-1) {
		case lua.TypeNil:
			err = nil
		case lua.TypeBoolean:
			if l.ToBoolean(-1) {
				err = filepath.SkipDir
			} else {
				err = errors.New("stop")
			}
		default:
			err = errors.New(l.ToString(-1))
		}
		l.Pop(1)
		return err
	})
	if err == nil {
		l.Push(nil)
	} else {
		l.Push(err.Error())
	}
	return 1
}
