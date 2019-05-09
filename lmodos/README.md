# Operating System Facilities

This library is implemented through table `os`.

## Documentation

### os.abs(path)

Returns an absolute representation of `path`.

### os.clock()

Returns an approximation of the amount in seconds of CPU time used by the program.

### os.date([format[, time]])

Returns a string or a table containing date and time, formatted according to the given string format.

Please Refer to [Lua 5.3 Reference Manual](http://www.lua.org/manual/5.3/manual.html#6.9).

### os.difftime(t2, t1)

Returns the difference, in seconds, from time `t1` to time `t2` (where the times are values returned by `os.time`).

### os.execute(command[, arg1, ...])

Executes the named program with the given arguments.
Its first result is `true` if the command terminated successfully, or `nil` otherwise.
After this first result the function returns a string plus a number, as follows:
```
'exit': the command terminated normally; the following number is the exit status of the command.
'signal': the command was terminated by a signal; the following number is the signal that terminated the command.
```

### os.exists(filename)

Determines whether the file named by `filename` exists.

### os.exit([code])

Terminates the host program.
If `code` is `true`, the returned status is `EXIT_SUCCESS`.
If `code` is `false`, the returned status is `EXIT_FAILURE`.
If `code` is a number, the returned status is this number.
The default value for `code` is `true`.

### os.getenv(name)

Returns the value of the process environment variable `name`, or `nil` if the variable is not defined.

### os.join(elem1[, elem2, ...])

Joins any number of path elements into a single path, adding a separator if necessary.

### os.mkdir(path[, all])

If `all` is `false` or `nil`, creates a new directory with the specified `path`.
If `all` is `true`, creates a directory named `path`, along with any necessary parents.

### os.open(filename[, mode])

This function opens a file, in the mode specified in the string mode.
In case of success, it returns a new file handle.

The mode string can be any of the following:
```
'r': read mode (the default);
'w': write mode;
'a': append mode;
'r+': update mode, all previous data is preserved;
'w+': update mode, all previous data is erased;
'a+': append update mode, previous data is preserved, writing is only allowed at the end of file.
```

### os.popen(prog[, mode[, arg1, ...]])

Starts program `prog` in a separated process and returns a file handle that you can use to read data from this program (if mode is "r", the default) or to write data to this program (if mode is "w").

### os.remove(filename)

Deletes the file (or empty directory, on POSIX systems) with the given name.
If this function fails, it returns `nil`, plus a string describing the error.
Otherwise, it returns `true`.

### os.rename(oldname, newname)

Renames the file or directory named `oldname` to `newname`.
If this function fails, it returns `nil`, plus a string describing the error.
Otherwise, it returns `true`.

### os.root

The root path of the executable that started the current process.

### os.setlocale(locale[, category])

Sets the current locale of the program.
The function returns the name of the new locale, or nil if the request cannot be honored.

### os.stat(file)

Returns the information describing `file`.
`file` can be a string or a file handle.

### os.time([table])

Returns the current time when called without arguments, or a time representing the local date and time specified by the given table.

Please Refer to [Lua 5.3 Reference Manual](http://www.lua.org/manual/5.3/manual.html#6.9).

### os.tmpfile()

In case of success, returns a handle for a temporary file.
It is the caller's responsibility to remove the file when no longer needed.

### os.tmpname()

Returns a string with a file name that can be used for a temporary file.

### os.walk(root, walkfn)

Walks the file tree rooted at `root`, calling walkFn for each file or directory in the tree, including `root`.
All errors that arise visiting files and directories are filtered by `walkfn`.

The `walkfn` is similar to `Go filepath.WalkFunc`.
If the function returns `true` when invoked on a non-directory file, Walk skips the remaining files in the containing directory.

### info.name

Base name of the file.

### info.isdir

Reports whether m describes a directory.

### info.size

Length in bytes for regular files; system-dependent for others.

### info.modtime

Modification time.

### file:close()

Closes `file`.

### file:seek([whence[, offset]])

Sets and gets the file position, measured from the beginning of the file, to the position given by `offset` plus a base specified by the string `whence`, as follows:
```
'set': base is position 0 (beginning of the file);
'cur': base is current position;
'end': base is end of file;
```
In case of success, `seek` returns the final file position, measured in bytes from the beginning of the file. If `seek` fails, it returns nil, plus a string describing the error.

The default value for whence is 'cur', and for offset is 0.

### prog:close()

Closes `prog`.
