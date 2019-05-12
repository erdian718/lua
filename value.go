/*
Copyright 2016-2017 by Milo Christiansen

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

package lua

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type TypeID int
type STypeID int

const (
	TypeUnknown TypeID = iota
	TypeNil
	TypeNumber
	TypeString
	TypeBoolean
	TypeTable
	TypeFunction
	TypeUserData

	nType int = iota
)

const (
	STypeUnknown STypeID = iota
	STypeInteger
	STypeFloat
)

var typeNames = []string{"unknown", "nil", "number", "string", "boolean", "table", "function", "userdata"}
var stypeNames = []string{"unknown", "integer", "float"}

func (typ TypeID) String() string {
	return typeNames[typ]
}

func (typ STypeID) String() string {
	return stypeNames[typ]
}

type value interface{}

type userdata struct {
	data interface{}
	meta *table
}

func typeOf(v value) TypeID {
	switch v.(type) {
	case nil:
		return TypeNil
	case float64, int64:
		return TypeNumber
	case string:
		return TypeString
	case bool:
		return TypeBoolean
	case *table:
		return TypeTable
	case *function:
		return TypeFunction
	case *userdata:
		return TypeUserData
	default:
		return TypeUnknown
	}
}

func stypeOf(v value) STypeID {
	switch v.(type) {
	case float64:
		return STypeFloat
	case int64:
		return STypeInteger
	default:
		return STypeUnknown
	}
}

func toBoolean(v value) bool {
	switch x := v.(type) {
	case bool:
		return x
	case nil:
		return false
	default:
		return true
	}
}

func toString(v value) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	case int64:
		return strconv.FormatInt(x, 10)
	case bool:
		return strconv.FormatBool(x)
	case nil:
		return "nil"
	default:
		return typeOf(x).String() + ": " + strconv.FormatUint(uint64(reflect.ValueOf(x).Pointer()), 16)
	}
}

func tryFloat(v value) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case int64:
		return float64(x), nil
	case string:
		return strconv.ParseFloat(strings.TrimSpace(x), 64)
	default:
		return 0, errors.New("can't convert to float: " + toString(x))
	}
}

func toFloat(v value) float64 {
	if x, err := tryFloat(v); err == nil {
		return x
	} else {
		panic(err)
	}
}

func tryInteger(v value) (int64, error) {
	switch x := v.(type) {
	case int64:
		return x, nil
	case float64:
		y := int64(x)
		if float64(y) == x {
			return y, nil
		} else {
			return 0, errors.New("can't convert to integer: " + toString(x))
		}
	case string:
		x = strings.TrimSpace(x)
		if y, err := strconv.ParseInt(x, 0, 64); err == nil {
			return y, nil
		} else if y, err := strconv.ParseFloat(x, 64); err == nil {
			z := int64(y)
			if float64(z) == y {
				return z, nil
			} else {
				return 0, errors.New("can't convert to integer: " + x)
			}
		} else {
			return 0, err
		}
	default:
		return 0, errors.New("can't convert to integer: " + toString(x))
	}
}

func toInteger(v value) int64 {
	if x, err := tryInteger(v); err == nil {
		return x
	} else {
		panic(err)
	}
}

func (l *State) getTable(t, k value) value {
	tbl, ok := t.(*table)
	if ok {
		if v := tbl.Get(k); v != nil {
			return v
		}
	}

	meth := l.getMetaField(t, "__index")
	if meth != nil {
		if tbl, ok := meth.(*table); ok {
			return l.getTable(tbl, k)
		}

		f, ok := meth.(*function)
		if !ok {
			panic("Meta method __index is not a table or function.")
		}

		l.Push(f)
		l.Push(t)
		l.Push(k)
		l.Call(2, 1)
		rtn := l.stack.Get(-1)
		l.Pop(1)
		return rtn
	}

	if ok {
		return tbl.Get(k)
	}
	panic("Value is not a table and has no __index meta method.")
	panic("UNREACHABLE")
}

func (l *State) setTable(t, k, v value) {
	tbl, ok := t.(*table)
	if ok && tbl.Get(k) != nil {
		tbl.Set(k, v)
		return
	}

	meth := l.getMetaField(t, "__newindex")
	if meth != nil {
		if t, ok := meth.(*table); ok {
			l.setTable(t, k, v)
			return
		}

		f, ok := meth.(*function)
		if !ok {
			panic("Meta method __newindex is not a table or function.")
		}

		l.Push(f)
		l.Push(t)
		l.Push(k)
		l.Push(v)
		l.Call(3, 0)
		return
	}
	if ok {
		tbl.Set(k, v)
		return
	}
	panic("Value is not a table and has no __newindex meta method.")
	panic("UNREACHABLE")
}

var mathMeta = [...]string{
	"__add",
	"__sub",
	"__mul",
	"__mod",
	"__pow",
	"__div",
	"__idiv",
	"__band",
	"__bor",
	"__bxor",
	"__shl",
	"__shr",
	"__unm",
	"__bnot",
}

func (l *State) tryMathMeta(op opCode, a, b value) value {
	if op < OpAdd || op > OpBinNot {
		panic("Operator passed to tryMathMeta out of range.")
	}
	name := mathMeta[op-OpAdd]

	meta := l.getMetaField(a, name)
	if meta == nil {
		meta = l.getMetaField(b, name)
		if meta == nil {
			panic("Neither operand has a " + name + " meta method.")
		}
	}

	f, ok := meta.(*function)
	if !ok {
		panic("Meta method " + name + " is not a function.")
	}

	l.Push(f)
	l.Push(a)
	l.Push(b)
	l.Call(2, 1)
	rtn := l.stack.Get(-1)
	l.Pop(1)
	return rtn
}

func (l *State) arith(op opCode, a, b value) value {
	switch op {
	case OpAdd:
		ia, oka := a.(int64)
		ib, okb := b.(int64)
		if oka && okb {
			return ia + ib
		}

		fa, erra := tryFloat(a)
		fb, errb := tryFloat(b)
		if erra == nil && errb == nil {
			return fa + fb
		}

		return l.tryMathMeta(op, a, b)
	case OpSub:
		ia, oka := a.(int64)
		ib, okb := b.(int64)
		if oka && okb {
			return ia - ib
		}

		fa, erra := tryFloat(a)
		fb, errb := tryFloat(b)
		if erra == nil && errb == nil {
			return fa - fb
		}

		return l.tryMathMeta(op, a, b)
	case OpMul:
		ia, oka := a.(int64)
		ib, okb := b.(int64)
		if oka && okb {
			return ia * ib
		}

		fa, erra := tryFloat(a)
		fb, errb := tryFloat(b)
		if erra == nil && errb == nil {
			return fa * fb
		}

		return l.tryMathMeta(op, a, b)
	case OpMod:
		ia, oka := a.(int64)
		ib, okb := b.(int64)
		if oka && okb {
			return ia % ib
		}

		fa, erra := tryFloat(a)
		fb, errb := tryFloat(b)
		if erra == nil && errb == nil {
			return math.Mod(fa, fb)
		}

		return l.tryMathMeta(op, a, b)
	case OpPow:
		fa, erra := tryFloat(a)
		fb, errb := tryFloat(b)
		if erra == nil && errb == nil {
			return math.Pow(fa, fb)
		}

		return l.tryMathMeta(op, a, b)
	case OpDiv:
		fa, erra := tryFloat(a)
		fb, errb := tryFloat(b)
		if erra == nil && errb == nil {
			return fa / fb
		}

		return l.tryMathMeta(op, a, b)
	case OpIDiv:
		ia, erra := tryInteger(a)
		ib, errb := tryInteger(b)
		if erra == nil && errb == nil {
			return ia / ib
		}

		return l.tryMathMeta(op, a, b)
	case OpBinAND:
		ia, erra := tryInteger(a)
		ib, errb := tryInteger(b)
		if erra == nil && errb == nil {
			return ia & ib
		}

		return l.tryMathMeta(op, a, b)
	case OpBinOR:
		ia, erra := tryInteger(a)
		ib, errb := tryInteger(b)
		if erra == nil && errb == nil {
			return ia | ib
		}

		return l.tryMathMeta(op, a, b)
	case OpBinXOR:
		ia, erra := tryInteger(a)
		ib, errb := tryInteger(b)
		if erra == nil && errb == nil {
			return ia ^ ib
		}

		return l.tryMathMeta(op, a, b)
	case OpBinShiftL:
		ia, erra := tryInteger(a)
		ib, errb := tryInteger(b)
		if erra == nil && errb == nil {
			if ib < 0 {
				return int64(uint64(ia) >> uint64(-ib))
			} else {
				return int64(uint64(ia) << uint64(ib))
			}
		}

		return l.tryMathMeta(op, a, b)
	case OpBinShiftR:
		ia, erra := tryInteger(a)
		ib, errb := tryInteger(b)
		if erra == nil && errb == nil {
			if ib < 0 {
				return int64(uint64(ia) << uint64(-ib))
			} else {
				return int64(uint64(ia) >> uint64(ib))
			}
		}

		return l.tryMathMeta(op, a, b)
	case OpUMinus:
		ia, oka := a.(int64)
		if oka {
			return -ia
		}

		fa, erra := tryFloat(a)
		if erra == nil {
			return -fa
		}

		return l.tryMathMeta(op, a, b)
	case OpBinNot:
		ia, erra := tryInteger(a)
		if erra == nil {
			return ^ia
		}

		return l.tryMathMeta(op, a, b)
	default:
		panic("Invalid opCode passed to arith")
		panic("UNREACHABLE")
	}
}

var cmpMeta = [...]string{
	"__eq",
	"__lt",
	"__le", // if this does not exist then try !lt(b, a)
}

func (l *State) tryCmpMeta(op opCode, a, b value) bool {
	if op < OpEqual || op > OpLessOrEqual {
		panic("Operator passed to tryCmpMeta out of range.")
	}
	name := cmpMeta[op-OpEqual]

	var meta value
	tryLEHack := false
try:
	meta = l.getMetaField(a, name)
	if meta == nil {
		meta = l.getMetaField(b, name)
		if meta == nil {
			if name == "__le" {
				tryLEHack = true
				name = "__lt"
				goto try
			}
			if name == "__eq" {
				return a == b // Fall back to raw equality.
			}

			return false
		}
	}

	l.Push(meta)
	if tryLEHack {
		l.Push(b)
		l.Push(a)
	} else {
		l.Push(a)
		l.Push(b)
	}
	l.Call(2, 1)
	rtn := toBoolean(l.stack.Get(-1))
	l.Pop(1)
	if tryLEHack {
		return !rtn
	}
	return rtn
}

func (l *State) compare(op opCode, a, b value, raw bool) bool {
	tm := true
	t := typeOf(a)
	if t != typeOf(b) {
		tm = false
	}

	switch op {
	case OpEqual:
		if tm {
			switch t {
			case TypeNil:
				return true // Obviously.
			case TypeNumber:
				ia, oka := a.(int64)
				ib, okb := b.(int64)
				if oka && okb {
					return ia == ib
				}
				return toFloat(a) == toFloat(b)
			case TypeString:
				return a.(string) == b.(string)
			case TypeBoolean:
				return a.(bool) == b.(bool)
			}
		}

		if raw {
			return a == b
		}
		return l.tryCmpMeta(op, a, b)

	case OpLessThan:
		if tm {
			switch t {
			case TypeNumber:
				ia, oka := a.(int64)
				ib, okb := b.(int64)
				if oka && okb {
					return ia < ib
				}
				return toFloat(a) < toFloat(b)
			case TypeString:
				return a.(string) < b.(string) // Fix me, should be locale sensitive, not lexical
			}
		}

		if raw {
			return false
		}
		return l.tryCmpMeta(op, a, b)

	case OpLessOrEqual:
		if tm {
			switch t {
			case TypeNumber:
				ia, oka := a.(int64)
				ib, okb := b.(int64)
				if oka && okb {
					return ia <= ib
				}
				return toFloat(a) <= toFloat(b)
			case TypeString:
				return a.(string) <= b.(string) // Fix me, should be locale sensitive, not lexical
			}
		}

		if raw {
			return false
		}
		return l.tryCmpMeta(op, a, b)
	default:
		panic("Invalid comparison operator.")
		panic("UNREACHABLE")
	}
}

func toStringConcat(v value) string {
	switch v2 := v.(type) {
	case nil:
		panic("Attempt to concatenate a nil value.")
		panic("UNREACHABLE")
	case float64:
		return fmt.Sprintf("%g", v2)
	case int64:
		return fmt.Sprintf("%d", v2)
	case string:
		return v2
	case bool:
		panic("Attempt to concatenate a bool value.")
		panic("UNREACHABLE")
	case *table:
		panic("Attempt to concatenate a table value.")
		panic("UNREACHABLE")
	case *function:
		panic("Attempt to concatenate a function value.")
		panic("UNREACHABLE")
	case *userdata:
		panic("Attempt to concatenate a userdata value.")
		panic("UNREACHABLE")
	default:
		panic("Invalid type passed to toStringConcat.")
		panic("UNREACHABLE")
	}
}

func (l *State) getMetaTable(v value) *table {
	switch x := v.(type) {
	case *table:
		return x.meta
	case *userdata:
		return x.meta
	default:
		return l.meta[typeOf(v)]
	}
}

func (l *State) getMetaField(v value, name string) value {
	m := l.getMetaTable(v)
	if m == nil {
		return nil
	}
	return m.Get(name)
}
