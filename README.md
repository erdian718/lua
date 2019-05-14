# Go Lua Compiler and VM

This is a Lua 5.3 VM and compiler written in [Go](http://golang.org/). This is intended to allow easy embedding into Go programs, with minimal fuss and bother.

This repository is forked from [milochristiansen/lua](https://github.com/milochristiansen/lua), and I made some incompatible changes. So it's unlikely to be merged into the original repository.

The strftime function and string pattern matching is currently copied from [yuin/gopher-lua](https://github.com/yuin/gopher-lua). It may be rewritten in the future.

## Usage

```go
package main

import (
	"ofunc/lua/util"
)

func main() {
	l := util.NewState()
	util.Run(l, "main.lua")
}
```

Please refer to the README for standard libraries.

## Dependencies

[Go 1.12 or later](http://golang.org/)

## Modules

* [lmodoffice](https://github.com/ofunc/lmodoffice) - A simple Lua module for converting various office documents into OOXML format files.
* [lmodxlsx](https://github.com/ofunc/lmodxlsx) - plandem/xlsx bindings for Lua.
* [stream](https://github.com/ofunc/stream) - A simple lazy list module for Lua.

## Missing Stuff

The following standard functions are not available:

* `collectgarbage` (not possible, VM uses the Go collector)
* `xpcall` (VM has no concept of a message handler)
* `next` (I don't need it at this time)
* `load` (violates my security policy)
* `dofile` (violates my security policy)
* `loadfile` (violates my security policy)
* `string.dump` (violates my security policy)
* `string.pack` (I don't need it at this time)
* `string.packsize` (I don't need it at this time)
* `string.unpack` (I don't need it at this time)

* * *

The following standard modules are not available:

* `package` (violates my security policy)
* `debug` (violates my security policy)
* `coroutine` (no coroutine support yet, just use goroutine)

* * *

The following modules are not implemented exactly as the Lua 5.3 specification requires:

* `math.random` (initialize the seed by default using program startup time)
* `os` (Go style functions)
* `io` (Go style functions)

* * *

There are a few things that are implemented exactly as the Lua 5.3 specification requires, where the reference
Lua implementation does not follow the specification exactly:

* The `#` (length) operator always returns the number of positive integer keys. When the table is a sequence, it's exactly equal to the sequence's length.
* Modulo operator (`%`) is implemented the same way most languages implement it, not the way Lua does. This does not matter unless you are using negative operands.

* * *

The following *core language* features are not supported:

* Hexadecimal floating point literals are not supported at this time.
* Weak references of any kind are not supported. This is because the VM uses Go's garbage collector, and it does not support weak references.
* Finalizers are not supported.

## TODO

* More tests.
* Better error hints.
* Code optimization.
