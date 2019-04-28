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
)

const (
	// Registry table index
	RegistryIndex = -1000000 - iota
	// Globals table index
	GlobalsIndex
	// To get a specific upvalue use "FirstUpVal-<upvalue index>"
	FirstUpVal
)

// State is the central arbitrator of all Lua operations.
type State struct {
	// Add a native stack trace to errors that have attached stack traces.
	NativeTrace bool

	stack    *stack
	registry *table
	global   *table
	meta     [nType]*table
}

// NewState creates a new State, ready to use.
func NewState() *State {
	l := &State{
		stack: newStack(),
	}

	l.global = newTable(l, 0, 64)
	l.global.SetRaw("_G", l.global)

	l.registry = newTable(l, 0, 32)
	l.registry.SetRaw("LUA_RIDX_GLOBALS", l.global)

	return l
}

// Helper
func (l *State) get(i int) value {
	switch {
	case i == RegistryIndex:
		return l.registry
	case i == GlobalsIndex:
		return l.global
	case i <= FirstUpVal:
		return l.stack.cFrame().getUp(FirstUpVal - i)
	case i > 0:
		return l.stack.Get(i - 1)
	case i < 0:
		return l.stack.Get(i)
	default:
		return nil
	}
}

// Used to create the return values for the compiler API functions (nothing else!).
func (l *State) asFunc(proto *funcProto, env *table) *function {
	up := make([]*upValue, len(proto.upVals))
	for i := range up {
		def := proto.upVals[i].makeUp()
		// Don't set name or index! name may come in from debug info, index is meaningless when closed.
		def.closed = true
		def.absIdx = -1
		up[i] = def
	}

	// Top level functions must have their first upvalue as _ENV
	if len(up) > 0 {
		if name := up[0].name; name != "_ENV" && name != "" {
			panic(errors.New("top level function without _ENV or _ENV in improper position"))
		}
		up[0].val = env
	}

	return &function{
		proto: *proto,
		up:    up,
	}
}
