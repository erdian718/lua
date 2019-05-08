local test = {}

function test.print()
	assert(print('Hello, Go Lua v' .. _VERSION) == nil)
end

function test.error()
	local a, b, c = assert(true, 'ABC', 123)
	assert(a == true and b == 'ABC' and c == 123)
	assert(pcall(assert, true))
	assert(pcall(assert, 0))
	assert(pcall(assert, ''))
	assert(not pcall(error, 'ERROR'))
	assert(not pcall(assert, false))
	assert(not pcall(assert, nil))
	local x, y, z = pcall(function()
		return 'ABC', 123
	end)
	assert(x == true and y == 'ABC' and z == 123)
end

function test.metatable()
	local t1, t2 = {}, {}
	local m1, m2 = {}, {__metatable = 'ABC'}
	assert(setmetatable(t1, m1) == t1)
	assert(setmetatable(t2, m2) == t2)
	assert(getmetatable(t1) == m1)
	assert(getmetatable(t2) == 'ABC')
	assert(setmetatable(t1, m2) == t1)
	assert(getmetatable(t1) == 'ABC')
	assert(pcall(setmetatable, t2, m1) == false)
	assert(getmetatable(123) == nil)
	assert(getmetatable(nil) == nil)
end

function test.ipairs()
	local xs = {'A', 'B', 'C'}
	local iter, a, k = ipairs(xs)
	local v
	k, v = iter(a, k)
	assert(k == 1 and v == 'A')
	k, v = iter(a, k)
	assert(k == 2 and v == 'B')
	k, v = iter(a, k)
	assert(k == 3 and v == 'C')
	k, v = iter(a, k)
	assert(k == nil and v == nil)
end

function test.pairs()
	local xs = {k1 = 'A', k2 = 'B', k3 = 'C'}
	local iter, a, k = pairs(xs)
	local k1, v1 = iter(a, k)
	assert(xs[k1] == v1)
	local k2, v2 = iter(a, k1)
	assert(xs[k2] == v2)
	local k3, v3 = iter(a, k2)
	assert(xs[k3] == v3)
	local k4, v4 = iter(a, k3)
	assert(k4 == nil and v4 == nil)
	assert(k1 ~= k2 and k2 ~= k3 and k3 ~= k1)
end

function test.raw()
	local xs = setmetatable({}, {
		__index = error;
		__newindex = error;
		__len = error;
		__eq = error;
	})
	assert(rawequal(rawset(xs, 1, 'A'), xs))
	assert(rawget(xs, 1))
	assert(rawlen(xs) == 1)
	assert(pcall(rawget, 123, 1) == false)
	assert(pcall(rawset, 123, 1, 'A') == false)
end

function test.select()
	assert(select('#', 'A', true, 123) == 3)
	assert(select('#', 'A', true, 123, nil) == 4)
	local r1, r2, r3 = select(1, 'A', true, 123)
	assert(r1 == 'A' and r2 == true and r3 == 123)
	local r1, r2, r3 = select(2, 'A', true, 123)
	assert(r1 == true and r2 == 123 and r3 == nil)
	local r1, r2, r3 = select(3, 'A', true, 123)
	assert(r1 == 123 and r2 == nil and r3 == nil)
	local r1, r2, r3 = select(4, 'A', true, 123)
	assert(r1 == nil and r2 == nil and r3 == nil)
	local r1, r2, r3 = select(-1, 'A', true, 123)
	assert(r1 == 123 and r2 == nil and r3 == nil)
	local r1, r2, r3 = select(-2, 'A', true, 123)
	assert(r1 == true and r2 == 123 and r3 == nil)
	local r1, r2, r3 = select(-3, 'A', true, 123)
	assert(r1 == 'A' and r2 == true and r3 == 123)
	assert(pcall(select, -4, 'A', true, 123) == false)
	assert(select('#', select(2, 'A', true, '123')) == 2)
end

function test.tonumber()
	assert(tonumber(123) == 123)
	assert(tonumber(12.34) == 12.34)
	assert(tonumber('123') == 123)
	assert(tonumber('12.34') == 12.34)
	assert(tonumber(nil) == nil)
	assert(tonumber('ABC') == nil)
	assert(tonumber('0xABC') == 2748)
end

function test.tostring()
	assert(tostring(nil) == 'nil')
	assert(tostring(true) == 'true')
	assert(tostring(false) == 'false')
	assert(tostring(123) == '123')
	assert(tostring(12.34) == '12.34')
	assert(tostring(-12.34) == '-12.34')
	assert(tostring('ABC') == 'ABC')
	local x = setmetatable({k = 'ABC'}, {__tostring = function(v)
		return v.k .. '123'
	end})
	assert(tostring(x) == 'ABC123')
end

function test.type()
	assert(type(nil) == 'nil')
	assert(type(true) == 'boolean')
	assert(type(false) == 'boolean')
	assert(type(123) == 'number')
	assert(type('ABC') == 'string')
	assert(type({}) == 'table')
	assert(type(error) == 'function')
end

function test.package()
	local seq1 = require 'seq'
	local seq2 = require 'seq'
	assert(seq1 == seq2)
	assert(seq1.next() == 123456)
	assert(seq1.next() ~= seq1.next())

	local pkg1 = require 'test/pkg'
	local pkg2 = require 'test/pkg'
	assert(pkg1 == pkg2)
	assert(pkg1 == 123459)
end

return test
