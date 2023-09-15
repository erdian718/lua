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

package lmodmath

import (
	"math"
	"math/rand"
	"time"

	"github.com/ofunc/lua"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Open opens the module.
func Open(l *lua.State) int {
	l.NewTable(0, 32)

	l.Push("abs")
	l.Push(labs)
	l.SetTableRaw(-3)

	l.Push("acos")
	l.Push(lacos)
	l.SetTableRaw(-3)

	l.Push("asin")
	l.Push(lasin)
	l.SetTableRaw(-3)

	l.Push("atan")
	l.Push(latan)
	l.SetTableRaw(-3)

	l.Push("ceil")
	l.Push(lceil)
	l.SetTableRaw(-3)

	l.Push("cos")
	l.Push(lcos)
	l.SetTableRaw(-3)

	l.Push("deg")
	l.Push(ldeg)
	l.SetTableRaw(-3)

	l.Push("e")
	l.Push(math.E)
	l.SetTableRaw(-3)

	l.Push("exp")
	l.Push(lexp)
	l.SetTableRaw(-3)

	l.Push("floor")
	l.Push(lfloor)
	l.SetTableRaw(-3)

	l.Push("fmod")
	l.Push(lfmod)
	l.SetTableRaw(-3)

	l.Push("huge")
	l.Push(math.Inf(1))
	l.SetTableRaw(-3)

	l.Push("log")
	l.Push(llog)
	l.SetTableRaw(-3)

	l.Push("max")
	l.Push(lmax)
	l.SetTableRaw(-3)

	l.Push("maxinteger")
	l.Push(int64(math.MaxInt64))
	l.SetTableRaw(-3)

	l.Push("min")
	l.Push(lmin)
	l.SetTableRaw(-3)

	l.Push("mininteger")
	l.Push(int64(math.MinInt64))
	l.SetTableRaw(-3)

	l.Push("modf")
	l.Push(lmodf)
	l.SetTableRaw(-3)

	l.Push("pi")
	l.Push(math.Pi)
	l.SetTableRaw(-3)

	l.Push("rad")
	l.Push(lrad)
	l.SetTableRaw(-3)

	l.Push("random")
	l.Push(lrandom)
	l.SetTableRaw(-3)

	l.Push("randomseed")
	l.Push(lrandomseed)
	l.SetTableRaw(-3)

	l.Push("sin")
	l.Push(lsin)
	l.SetTableRaw(-3)

	l.Push("sqrt")
	l.Push(lsqrt)
	l.SetTableRaw(-3)

	l.Push("tan")
	l.Push(ltan)
	l.SetTableRaw(-3)

	l.Push("tointeger")
	l.Push(ltointeger)
	l.SetTableRaw(-3)

	l.Push("type")
	l.Push(ltype)
	l.SetTableRaw(-3)

	l.Push("ult")
	l.Push(lult)
	l.SetTableRaw(-3)

	return 1
}

func labs(l *lua.State) int {
	if v, err := l.TryInteger(1); err == nil {
		if v >= 0 {
			l.Push(v)
		} else {
			l.Push(-v)
		}
	} else {
		l.Push(math.Abs(l.ToFloat(1)))
	}
	return 1
}

func lacos(l *lua.State) int {
	l.Push(math.Acos(l.ToFloat(1)))
	return 1
}

func lasin(l *lua.State) int {
	l.Push(math.Asin(l.ToFloat(1)))
	return 1
}

func latan(l *lua.State) int {
	x := l.ToFloat(1)
	y := l.OptFloat(2, 1)
	l.Push(math.Atan2(x, y))
	return 1
}

func lceil(l *lua.State) int {
	v := math.Ceil(l.ToFloat(1))
	x := int64(v)
	if float64(x) == v {
		l.Push(x)
	} else {
		l.Push(v)
	}
	return 1
}

func lcos(l *lua.State) int {
	l.Push(math.Cos(l.ToFloat(1)))
	return 1
}

func ldeg(l *lua.State) int {
	l.Push(180 * l.ToFloat(1) / math.Pi)
	return 1
}

func lexp(l *lua.State) int {
	l.Push(math.Exp(l.ToFloat(1)))
	return 1
}

func lfloor(l *lua.State) int {
	v := math.Floor(l.ToFloat(1))
	x := int64(v)
	if float64(x) == v {
		l.Push(x)
	} else {
		l.Push(v)
	}
	return 1
}

func lfmod(l *lua.State) int {
	x := l.ToFloat(1)
	y := l.ToFloat(2)
	l.Push(math.Mod(x, y))
	return 1
}

func llog(l *lua.State) int {
	v := math.Log(l.ToFloat(1))
	if !l.IsNil(2) {
		v = v / math.Log(l.ToFloat(2))
	}
	l.Push(v)
	return 1
}

func lmax(l *lua.State) int {
	n, m := l.AbsIndex(-1), 1
	for i := 1; i <= n; i++ {
		if l.Compare(m, i, lua.OpLessThan) {
			m = i
		}
	}
	l.PushIndex(m)
	return 1
}

func lmin(l *lua.State) int {
	n, m := l.AbsIndex(-1), 1
	for i := 1; i <= n; i++ {
		if l.Compare(i, m, lua.OpLessThan) {
			m = i
		}
	}
	l.PushIndex(m)
	return 1
}

func lmodf(l *lua.State) int {
	a, b := math.Modf(l.ToFloat(1))
	x := int64(a)
	if float64(x) == a {
		l.Push(x)
	} else {
		l.Push(a)
	}
	l.Push(b)
	return 2
}

func lrad(l *lua.State) int {
	l.Push(math.Pi * l.ToFloat(1) / 180)
	return 1
}

func lrandom(l *lua.State) int {
	switch l.AbsIndex(-1) {
	case 0:
		l.Push(rand.Float64())
	case 1:
		l.Push(rand.Int63n(l.ToInteger(1)) + 1)
	default:
		m := l.ToInteger(1)
		n := l.ToInteger(2)
		l.Push(rand.Int63n(n-m+1) + m)
	}
	return 1
}

func lrandomseed(l *lua.State) int {
	rand.Seed(l.ToInteger(1))
	return 0
}

func lsin(l *lua.State) int {
	l.Push(math.Sin(l.ToFloat(1)))
	return 1
}

func lsqrt(l *lua.State) int {
	l.Push(math.Sqrt(l.ToFloat(1)))
	return 1
}

func ltan(l *lua.State) int {
	l.Push(math.Tan(l.ToFloat(1)))
	return 1
}

func ltointeger(l *lua.State) int {
	if v, err := l.TryInteger(1); err == nil {
		l.Push(v)
	} else {
		l.Push(nil)
	}
	return 1
}

func ltype(l *lua.State) int {
	if t := l.STypeOf(1); t == lua.STypeUnknown {
		l.Push(nil)
	} else {
		l.Push(t.String())
	}
	return 1
}

func lult(l *lua.State) int {
	x := l.ToInteger(1)
	y := l.ToInteger(2)
	l.Push(uint64(x) < uint64(y))
	return 1
}
