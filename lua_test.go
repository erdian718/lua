package lua_test

import (
	"fmt"
	"testing"

	"ofunc/lua"
	"ofunc/lua/util"
)

func TestMain(m *testing.M) {
	seq := 123456
	l := util.NewState()
	l.Preload("seq", func(l *lua.State) int {
		l.NewTable(0, 2)
		l.Push("next")
		l.Push(func(l *lua.State) int {
			l.Push(seq)
			seq += 1
			return 1
		})
		l.SetTableRaw(-3)
		return 1
	})
	util.AddPath("")
	if err := util.Test(l, "test"); err != nil {
		fmt.Println("error:", err)
	}
}
