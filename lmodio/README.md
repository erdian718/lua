# Input and Output Facilities

This library is implemented through table `io`.

## Documentation

### io.buffer([buf])

Creates and initializes a new buffer using `buf` as its initial contents.

### io.copy(r, w[, n])

Copies `n` bytes (or until an error) from `r` to `w`.
It returns the number of bytes copied and the earliest error encountered while copying.

### io.limit(r, n)

Returns a Reader that reads from `r` but stops with `eof` after `n` bytes.

### io.read(r[, n])

Reads up to `n` bytes.
It returns the data it read and any error encountered.

### io.readfile(filename)

Reads the file named by `filename` and returns the contents.

### io.scanner(r[, split])

Returns a new Scanner to read from `r`.
The buildin `split` functions are `byte`, `line`, `rune` and `word`, defaults to `line`.

### io.stderr

Standard error file.

### io.stdin

Standard input file.

### io.stdout

Standard output file.

### io.type(x)

Returns the I/O type of `x`: `readwriter`, `reader`, `writer` or `nil`.

### io.write(w, data)

Writes `data` to `w`.
It returns the number of bytes written and any error encountered that caused the write to stop early.

### io.writefile(filename, data)

Writes `data` to a file named by `filename`.
If the file does not exist, `writefile` creates it.
Otherwise `writefile` truncates it before writing.
