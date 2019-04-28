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

// Go Lua Compiler and VM
//
// This is a Lua 5.3 VM and compiler written in Go.
// This is intended to allow easy embedding into Go programs, with minimal fuss and bother.
package lua

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
)

// AbsIndex converts the given index into an absolute index.
// Use -1 as the index to get the number of items currently on the stack.
func (l *State) AbsIndex(i int) int {
	if i >= 0 || i <= RegistryIndex {
		return i
	}
	// Need to add 2 so we get a 1 based index.
	return l.stack.TopIndex() + i + 2
}

// Pop removes the top n items from the stack.
func (l *State) Pop(n int) {
	l.stack.Pop(n)
}

// Push pushes the given value onto the stack.
// If the value is not one of nil, float32, float64, int, int32, int64, string or bool,
// it is converted to a userdata value before being pushed.
func (l *State) Push(v interface{}) {
	switch x := v.(type) {
	case nil, bool, string, float64, int64, *table, *function, *userdata:
	case float32:
		v = float64(x)
	case int:
		v = int64(x)
	case int32:
		v = int64(x)
	case func(*State) int:
		v = &function{
			native: x,
			up: []*upValue{
				{
					name:   "_ENV",
					index:  -1,
					closed: true,
					val:    l.global,
					absIdx: -1,
				},
			},
		}
	default:
		v = &userdata{
			data: x,
		}
	}
	l.stack.Push(v)
}

// PushClosure pushes a native function as a closure.
func (l *State) PushClosure(f func(*State) int, v ...int) {
	c := len(v)
	if c == 0 {
		l.Push(f)
		return
	}
	c++

	fn := &function{
		native: f,
		up:     make([]*upValue, c),
	}
	// ALL native functions ALWAYS have their first upvalue set to the global table.
	// This differs from standard Lua, but doesn't hurt anything.
	fn.up[0] = &upValue{
		name:   "_ENV",
		index:  -1,
		closed: true,
		val:    l.global,
		absIdx: -1,
	}
	for i := 1; i < c; i++ {
		fn.up[i] = &upValue{
			name:   "(native upvalue)",
			index:  -1,
			closed: true,
			val:    l.get(v[i-1]),
			absIdx: -1,
		}
	}
	l.stack.Push(fn)
}

// PushIndex pushes a copy of the value at the given index onto the stack.
func (l *State) PushIndex(i int) {
	l.stack.Push(l.get(i))
}

// Insert takes the item from the TOS and inserts it at the given stack index.
func (l *State) Insert(i int) {
	if i >= 1 {
		i = i - 1
	}
	v := l.get(-1)
	l.Pop(1)
	l.stack.Insert(i, v)
}

// GetRaw gets the raw data for a Lua value.
// Lua types use the following mapping:
//     nil -> nil
//     number -> int64 or float64
//     string -> string
//     bool -> bool
//     table -> string: "table: <pointer as hexadecimal>"
//     function -> string: "function: <pointer as hexadecimal>"
//     userdata -> The raw user data value
func (l *State) GetRaw(i int) interface{} {
	switch v := l.get(i).(type) {
	case *userdata:
		return v.data
	case nil:
		return v
	case float64:
		return v
	case int64:
		return v
	case string:
		return v
	case bool:
		return v
	default:
		return toString(v)
	}
}

// Set sets the value at index d to the value at index s (d = s).
// Trying to set the registry or an invalid index will do nothing.
// Setting an absolute index will never fail, the stack will be extended as needed.
// Be careful not to waste stack space or you could run out of memory!
// This function is mostly for setting up-values and things like that.
func (l *State) Set(d, s int) {
	v := l.get(s)
	switch {
	case d == RegistryIndex:
		// Do nothing.
	case d == GlobalsIndex:
		// Do nothing.
	case d <= FirstUpVal:
		l.stack.cFrame().setUp(d-FirstUpVal, v)
	case d >= 1:
		l.stack.Set(d-1, v)
	case d < 0:
		l.stack.Set(d, v)
	default:
		// Do nothing.
	}
}

// SetGlobal pops a value from the stack and sets it as the new value of global name.
func (l *State) SetGlobal(name string) {
	l.global.SetRaw(name, l.stack.Get(-1))
	l.stack.Pop(1)
}

// SetUpVal sets upvalue "i" in the function at "f" to the value at "v".
// If the upvalue index is out of range, "f" is not a function, or the upvalue is not closed, false is returned and nothing is done, else returns true and sets the upvalue.
// Any other functions that share this upvalue will also be affected!
func (l *State) SetUpValue(f, i, v int) bool {
	fn, ok := l.get(f).(*function)
	if !ok || i >= len(fn.up) {
		return false
	}
	def := fn.up[i]
	if !def.closed {
		return false
	}
	def.val = l.get(v)
	return true
}

// Preload adds the given loader function for "require".
func (l *State) Preload(name string, loader func(*State) int) {
	loaded, ok := l.registry.GetRaw("_PRELOAD").(*table)
	if !ok {
		loaded = newTable(l, 0, 16)
		l.registry.SetRaw("_PRELOAD", loaded)
	}
	l.Push(loader)
	loaded.SetRaw(name, l.get(-1))
	l.Pop(1)
}

// IsNil check if the value at the given index is nil. Nonexistent values are always nil.
func (l *State) IsNil(i int) bool {
	return l.get(i) == nil
}

// TypeOf returns the type of the value at the given index.
func (l *State) TypeOf(i int) TypeID {
	return typeOf(l.get(i))
}

// STypeOf returns the sub-type of the value at the given index.
func (l *State) STypeOf(i int) STypeID {
	return stypeOf(l.get(i))
}

// ToBoolean reads a value from the stack at the given index and interprets it as a boolean.
func (l *State) ToBoolean(i int) bool {
	return toBoolean(l.get(i))
}

// ToString reads a value from the stack at the given index and formats it as a string.
// This will call a __tostring metamethod if provided.
// This is safe if no metamethods are called, but may panic if the metamethod errors out.
func (l *State) ToString(i int) string {
	v := l.get(i)
	if m := l.getMetaField(v, "__tostring"); m == nil {
		return toString(v)
	} else {
		l.Push(m)
		l.Push(v)
		l.Call(1, 1)
		r := l.stack.Get(-1)
		l.Pop(1)
		return toString(r)
	}
}

// ToFloat reads a floating point value from the stack at the given index.
// If the value is not an float and cannot be converted to one this may panic.
func (l *State) ToFloat(i int) float64 {
	return toFloat(l.get(i))
}

// ToInteger reads an integer value from the stack at the given index.
// If the value is not an integer and cannot be converted to one this may panic.
func (l *State) ToInteger(i int) int64 {
	return toInteger(l.get(i))
}

// TryFloat attempts to read the value at the given index as a floating point number.
func (l *State) TryFloat(i int) (float64, error) {
	return tryFloat(l.get(i))
}

// TryInteger attempts to read the value at the given index as a integer number.
func (l *State) TryInteger(i int) (int64, error) {
	return tryInteger(l.get(i))
}

// OptString is the same as ToString, except the given default is returned if the value is nil or non-existent.
func (l *State) OptString(i int, d string) string {
	if l.IsNil(i) {
		return d
	} else {
		return l.ToString(i)
	}
}

// OptFloat is the same as ToFloat, except the given default is returned if the value is nil or non-existent.
func (l *State) OptFloat(i int, d float64) float64 {
	if v := l.get(i); v == nil {
		return d
	} else {
		return toFloat(v)
	}
}

// OptInteger is the same as ToInt, except the given default is returned if the value is nil or non-existent.
func (l *State) OptInteger(i int, d int64) int64 {
	if v := l.get(i); v == nil {
		return d
	} else {
		return toInteger(v)
	}
}

// Arith performs the specified the arithmetic operator with the top two items on the stack (or just the top item for OpUMinus and OpBinNot).
// The result is pushed onto the stack.
// This may raise an error if they values are not appropriate for the given operator.
func (l *State) Arith(op opCode) {
	a := l.stack.Get(-2)
	b := a
	if op != OpUMinus && op != OpBinNot {
		b = l.stack.Get(-1)
	}
	l.stack.Pop(2)
	l.stack.Push(l.arith(op, a, b))
}

// Compare performs the specified the comparison operator with the items at the given stack indexes.
// This may raise an error if they values are not appropriate for the given operator.
func (l *State) Compare(i1, i2 int, op opCode) bool {
	a := l.get(i1)
	b := l.get(i2)
	return l.compare(op, a, b, false)
}

// CompareRaw is exactly like Compare, but without meta-methods.
func (l *State) CompareRaw(i1, i2 int, op opCode) bool {
	a := l.get(i1)
	b := l.get(i2)
	return l.compare(op, a, b, true)
}

// NewTable creates a new table with "as" preallocated array elements and "hs" preallocated hash elements.
func (l *State) NewTable(as, hs int) {
	l.stack.Push(newTable(l, as, hs))
}

// GetTable reads from the table at the given index, popping the key from the stack and pushing the result.
// The type of the pushed object is returned.
// This may raise an error if the value is not a table or is lacking the __index meta method.
func (l *State) GetTable(i int) TypeID {
	v := l.getTable(l.get(i), l.stack.Get(-1))
	l.Pop(1)
	l.Push(v)
	return typeOf(v)
}

// GetTableRaw is like GetTable except it ignores meta methods.
// This may raise an error if the value is not a table.
func (l *State) GetTableRaw(i int) TypeID {
	x := l.get(i)
	k := l.stack.Get(-1)
	l.Pop(1)
	if t, ok := x.(*table); ok {
		v := t.GetRaw(k)
		l.Push(v)
		return typeOf(v)
	} else {
		panic(errors.New("not a table: " + toString(x)))
	}
}

// SetTable writes to the table at the given index, popping the key and value from the stack.
// This may raise an error if the value is not a table or is lacking the __newindex meta method.
// The value must be on TOS, the key TOS-1.
func (l *State) SetTable(i int) {
	l.setTable(l.get(i), l.stack.Get(-2), l.stack.Get(-1))
	l.Pop(2)
}

// SetTableRaw is like SetTable except it ignores meta methods.
// This may raise an error if the value is not a table.
func (l *State) SetTableRaw(i int) {
	x := l.get(i)
	k := l.stack.Get(-2)
	v := l.stack.Get(-1)
	l.Pop(2)
	if t, ok := x.(*table); ok {
		t.SetRaw(k, v)
	} else {
		panic(errors.New("not a table: " + toString(x)))
	}
}

// GetIter pushes a table iterator onto the stack.
// This value is type "userdata" and has a "__call" meta method. Calling the iterator will
// push the next key/value pair onto the stack. The key is not required for the next
// iteration, so unlike Next you must pop both values.
// The end of iteration is signaled by returning a single nil value.
// If the given value is not a table this will raise an error.
func (l *State) GetIter(i int) {
	x := l.get(i)
	if t, ok := x.(*table); ok {
		l.Push(newTableIter(t))
		l.NewTable(0, 1)
		l.Push("__call")
		l.Push(func(l *State) int {
			i := l.GetRaw(1).(*tableIter)
			k, v := i.Next()
			if k == nil {
				l.Push(k)
				return 1
			}
			l.Push(k)
			l.Push(v)
			return 2
		})
		l.SetTableRaw(-3)
		l.SetMetaTable(-2)
	} else {
		panic(errors.New("not a table: " + toString(x)))
	}
}

// ForEach is a fancy version of ForEachRaw that respects metamethods (to be specific, __pairs).
func (l *State) ForEach(t int, f func() bool) {
	t = l.AbsIndex(t)
	if typ := l.GetMetaField(t, "__pairs"); typ == TypeNil {
		l.GetIter(t)    // iter
		l.PushIndex(-1) // iter iter
		l.PushIndex(1)  // iter iter key
		l.Push(nil)     // iter iter key value
	} else {
		l.PushIndex(t)  // meta tbl
		l.Call(1, 3)    // iter key value
		l.PushIndex(-3) // iter key value iter
		l.Insert(-3)    // iter iter key value
	}
	l.Call(2, 2) // iter key value
	for !l.IsNil(-2) {
		if !f() {
			break
		}
		l.PushIndex(-3) // iter key value iter
		l.Insert(-3)    // iter iter key value
		l.Call(2, 2)    // iter key value
	}
	l.Pop(3)
}

// ForEachRaw is a simple wrapper around GetIter and is provided as a convenience.
// The given function is called once for every item in the table at t.
// For each call of the function the value is at -1 and the key at -2.
// You MUST keep the stack balanced inside the function!
// Do not pop the key and value off the stack before returning!
// The value returned by the iteration function determines if ForEach should return early.
// Return false to break, return true to continue to the next iteration.
func (l *State) ForEachRaw(t int, f func() bool) {
	l.GetIter(t)    // iter
	l.PushIndex(-1) // iter iter
	l.Call(0, 2)    // iter key value
	for !l.IsNil(-2) {
		if !f() {
			break
		}
		l.Pop(2)        // iter
		l.PushIndex(-1) // iter iter
		l.Call(0, 2)    // iter key value
	}
	l.Pop(3)
}

// Returns the "length" of the item at the given index, exactly like the "#" operator would.
// If this calls a meta method it may raise an error if the length is not an integer.
func (l *State) Length(i int) int {
	v := l.get(i)
	if s, ok := v.(string); ok {
		return len(s)
	} else if m := l.getMetaField(v, "__len"); m != nil {
		if f, ok := m.(*function); ok {
			l.Push(f)
			l.Push(v)
			l.Call(1, 1)
			r := l.stack.Get(-1)
			l.Pop(1)
			return int(toInteger(r))
		} else {
			panic(errors.New("meta method __len is not a function: " + toString(m)))
		}
	} else if t, ok := v.(*table); ok {
		return t.Length()
	} else {
		panic(errors.New("not a string or table and has no __len meta method: " + toString(v)))
	}
}

// Returns the length of the table or string at the given index. This does not call meta methods.
// If the value is not a table or string this will raise an error.
func (l *State) LengthRaw(i int) int {
	v := l.get(i)
	if s, ok := v.(string); ok {
		return len(s)
	} else if t, ok := v.(*table); ok {
		return t.Length()
	} else {
		panic(errors.New("not a string or table: " + toString(v)))
	}
}

// GetMetaField pushes the meta method with the given name for the item at the given index onto the stack, then returns the type of the pushed item.
// If the item does not have a meta table or does not have the specified method this does nothing and returns TypNil.
func (l *State) GetMetaField(i int, name string) TypeID {
	m := l.getMetaField(l.get(i), name)
	if m != nil {
		l.Push(m)
	}
	return typeOf(m)
}

// GetMetaTable gets the meta table for the value at the given index and pushes it onto the stack.
// If the value does not have a meta table then this returns false and pushes nothing.
func (l *State) GetMetaTable(i int) bool {
	if m := l.getMetaTable(l.get(i)); m == nil {
		return false
	} else {
		l.Push(m)
		return true
	}
}

// SetMetaTable pops a table from the stack and sets it as the meta table of the value at the given index.
// If the value is not a userdata or table then the meta table is set for ALL values of that type!
// If you try to set a metatable that is not a table or try to pass an invalid type this will raise an error.
func (l *State) SetMetaTable(i int) {
	x := l.get(i)
	t := l.stack.Get(-1)
	m, ok := t.(*table)
	l.stack.Pop(1)
	if !ok && t != nil {
		panic(errors.New("not a table or nil: " + toString(t)))
	}
	switch v := x.(type) {
	case *table:
		v.meta = m
	case *userdata:
		v.meta = m
	case nil:
		l.meta[TypeNil] = m
	case float64:
		l.meta[TypeNumber] = m
	case int64:
		l.meta[TypeNumber] = m
	case string:
		l.meta[TypeString] = m
	case bool:
		l.meta[TypeBoolean] = m
	case *function:
		l.meta[TypeFunction] = m
	default:
		panic(errors.New("unknown type: " + toString(x)))
	}
}

// Dump converts the Lua function at the given index to a binary chunk.
// The returned value may be used with LoadBinary to get a function equivalent to the dumped function (but without the original function's up values).
// Currently the "strip" argument does nothing.
// This (obviously) only works with Lua functions, trying to dump a native function or a non-function value will raise an error.
func (l *State) Dump(i int, strip bool) []byte {
	v := l.get(i)
	f, ok := l.get(i).(*function)
	if !ok {
		panic(errors.New("not a function: " + toString(v)))
	}
	if f.native != nil {
		panic(errors.New("can't dump native function"))
	}
	return dumpBin(&f.proto)
}

// LoadBinary loads a binary chunk into memory and pushes the result onto the stack.
// If there is an error it is returned and nothing is pushed.
// Set env to 0 to use the default environment.
func (l *State) LoadBinary(in io.Reader, name string, env int) error {
	proto, err := loadBin(in, name)
	if err != nil {
		return err
	}
	envv := l.global
	if env != 0 {
		x, ok := l.get(env), false
		envv, ok = x.(*table)
		if !ok {
			return errors.New("not a table: " + toString(x))
		}
	}
	l.stack.Push(l.asFunc(proto, envv))
	return nil
}

// LoadText loads a text chunk into memory and pushes the result onto the stack.
// If there is an error it is returned and nothing is pushed.
// Set env to 0 to use the default environment.
func (l *State) LoadText(in io.Reader, name string, env int) error {
	source, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	proto, err := compSource(string(source), name, 1)
	if err != nil {
		return err
	}
	envv := l.global
	if env != 0 {
		x, ok := l.get(env), false
		envv, ok = x.(*table)
		if !ok {
			return errors.New("not a table: " + toString(x))
		}
	}
	l.stack.Push(l.asFunc(proto, envv))
	return nil
}

// Call runs a function with the given number of arguments and results.
// The function must be on the stack just before the first argument.
// If this raises an error the stack is NOT unwound!
// Call this only from code that is below a call to PCall unless you want your State to be permanently trashed!
func (l *State) Call(args, rtns int) {
	if args < 0 {
		panic(errors.New("invalid arg count: " + toString(int64(args))))
	}
	fi := -(args + 1) // Generate a relative index for the function
	l.call(fi, args, rtns, false)
}

// PCall is exactly like Call, except instead of panicking when it encounters an error the error is cleanly recovered and returned.
// On error the stack is reset to the way it was before the call minus the function and it's arguments, the State may then be reused.
func (l *State) PCall(args, rtns int, trace bool) (msg interface{}) {
	frames := len(l.stack.frames)
	top := len(l.stack.data) - args - 1
	defer func() {
		msg = recover()
		if msg != nil {
			// Print trace
			if trace {
				sources := []string{}
				lines := []int{}
				for i := len(l.stack.frames) - 1; i >= frames; i-- {
					frame := l.stack.frames[i]
					if frame.fn.native == nil {
						sources = append(sources, frame.fn.proto.source)
						if int(frame.pc) < len(frame.fn.proto.lineInfo) {
							lines = append(lines, frame.fn.proto.lineInfo[frame.pc])
						} else if len(frame.fn.proto.lineInfo) > 0 {
							lines = append(lines, frame.fn.proto.lineInfo[len(frame.fn.proto.lineInfo)-1])
						} else {
							lines = append(lines, -1)
						}
					} else {
						sources = append(sources, "(native code)")
						lines = append(lines, -1)
					}
				}
				fmt.Println("error:", msg)
				for i := range sources {
					if lines[i] == -1 {
						fmt.Printf("    \"%v\"\n", sources[i])
						continue
					}
					fmt.Printf("    \"%v\": <line: %v>\n", sources[i], lines[i])
				}
				if l.NativeTrace {
					buf := make([]byte, 4096)
					buf = buf[:runtime.Stack(buf, true)]
					fmt.Printf("\nNative Trace:\n%s\n", buf)
				}
			}

			// Before we strip the stack we need to close all upvalues in the section we will be stripping, just in
			// case a closure was assigned to another upvalue.
			l.stack.frames[len(l.stack.frames)-1].closeUpAbs(top)

			// Make sure the stack is back to the way we found it, minus the function and it's arguments.
			l.stack.frames = l.stack.frames[:frames]
			for i := len(l.stack.data) - 1; i >= top; i-- {
				l.stack.data[i] = nil
			}
			l.stack.data = l.stack.data[:top]
		}
	}()
	l.Call(args, rtns)
	return
}

// Error pops a value off the top of the stack and raises it as an error.
func (l *State) Error() {
	msg := l.get(-1)
	l.stack.Pop(1)
	panic(msg)
}

// PrintStack prints some stack information for sanity checking during test runs.
func (l *State) PrintStack() {
	n := l.AbsIndex(-1)
	fmt.Println("++++++++")
	fmt.Println("D:", len(l.stack.data))
	fmt.Println("F:", len(l.stack.frames))
	for i := 1; i <= n; i++ {
		fmt.Printf("%v: %v\n", i, l.ToString(i))
	}
	fmt.Println("--------")
}
