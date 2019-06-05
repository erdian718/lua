# Access to the WebAssembly host environment

This library is implemented through table `js`.

## Documentation

### js.global

The JavaScript global object, usually `window` or `global`.

### js.type(x)

Returns the JavaScript type of `x`.

### js.value(x)

Returns `x` as a JavaScript value.

### js.new(x[, ...])

Uses JavaScript's `new` operator with value x as constructor and the given arguments.

### js.free(f)

Frees up resources allocated for function `f`.
