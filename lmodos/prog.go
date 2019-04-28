package lmodos

import (
	"io"
	"os/exec"

	"ofunc/lua"
)

type prog struct {
	cmd *exec.Cmd
	r   io.ReadCloser
	w   io.WriteCloser
}

// Read reads up to len(xs) bytes from the prog.
func (p prog) Read(xs []byte) (int, error) {
	return p.r.Read(xs)
}

// Write writes len(xs) bytes to the prog.
func (p prog) Write(xs []byte) (int, error) {
	return p.w.Write(xs)
}

//Close closes the prog, rendering it unusable for I/O.
func (p prog) Close() (err error) {
	err = p.cmd.Process.Kill()
	if e := p.cmd.Process.Release(); e != nil && err == nil {
		err = e
	}
	if p.r != nil {
		if e := p.r.Close(); e != nil && err == nil {
			err = e
		}
	} else if p.w != nil {
		if e := p.w.Close(); e != nil && err == nil {
			err = e
		}
	}
	return
}

func mateprog(l *lua.State) int {
	l.NewTable(0, 2)
	idx := l.AbsIndex(-1)

	l.Push("__index")
	l.PushIndex(idx)
	l.SetTableRaw(idx)

	l.Push("close")
	l.Push(lpclose)
	l.SetTableRaw(idx)

	return idx
}

func lpclose(l *lua.State) int {
	if c, ok := l.GetRaw(1).(io.Closer); ok {
		if err := c.Close(); err == nil {
			l.Push(true)
			return 1
		} else {
			l.Push(nil)
			l.Push(err.Error())
			return 2
		}
	} else {
		panic("os.prog: not a prog")
		return 0
	}
}
