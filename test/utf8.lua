local test = {}
local utf8 = require 'utf8'
local string = require 'string'

function test.char()
	assert(utf8.char() == '')
	assert(utf8.char(72, 101, 108, 108, 111, 65292, 19990, 30028) == 'Hello，世界')
end

function test.charpattern()
	local xs = {string.byte(utf8.charpattern, 1, -1)}
	assert(#xs == 14)
	assert(xs[1] == 91)
	assert(xs[2] == 0)
	assert(xs[3] == 45)
	assert(xs[4] == 127)
	assert(xs[5] == 194)
	assert(xs[6] == 45)
	assert(xs[7] == 244)
	assert(xs[8] == 93)
	assert(xs[9] == 91)
	assert(xs[10] == 128)
	assert(xs[11] == 45)
	assert(xs[12] == 191)
	assert(xs[13] == 93)
	assert(xs[14] == 42)
end

function test.codes()
	local s = 'Hello，世界'
	local next, x, p = utf8.codes(s)
	local p, c = next(x, p)
	assert(p == 1 and c == 72)
	local p, c = next(x, p)
	assert(p == 2 and c == 101)
	local p, c = next(x, p)
	assert(p == 3 and c == 108)
	local p, c = next(x, p)
	assert(p == 4 and c == 108)
	local p, c = next(x, p)
	assert(p == 5 and c == 111)
	local p, c = next(x, p)
	assert(p == 6 and c == 65292)
	local p, c = next(x, p)
	assert(p == 9 and c == 19990)
	local p, c = next(x, p)
	assert(p == 12 and c == 30028)
end

function test.codepoint()
	local s = 'Hello，世界'
	local xs = {utf8.codepoint(s)}
	assert(#xs == 1)
	assert(xs[1] == 72)

	local xs = {utf8.codepoint(s, 3)}
	assert(#xs == 1)
	assert(xs[1] == 108)

	local xs = {utf8.codepoint(s, 4, -2)}
	assert(#xs == 5)
	assert(xs[1] == 108)
	assert(xs[2] == 111)
	assert(xs[3] == 65292)
	assert(xs[4] == 19990)
	assert(xs[5] == 30028)
end

function test.len()
	local s = 'Hello，世界'
	assert(utf8.len(s) == 8)
	assert(utf8.len(s, 3) == 6)
	assert(utf8.len(s, 4, -2) == 5)
	local n, i = utf8.len(s, 8)
	assert(n == nil and i == 8)
end

function test.offset()
	local s = 'Hello，世界'
	assert(utf8.offset(s, 2) == 2)
	assert(utf8.offset(s, -2) == 9)
	assert(utf8.offset(s, 2, 2) == 3)
	assert(utf8.offset(s, -2, 2) == nil)
	assert(utf8.offset(s, -2, -3) == 6)
	assert(utf8.offset(s, 2, -6) == 12)
	assert(utf8.offset(s, 0, 2) == 2)
	assert(utf8.offset(s, 0, -2) == 12)
end

return test
