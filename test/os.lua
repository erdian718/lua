local test = {}
local os = require 'os'
local string = require 'string'
local math = require 'math'

function test.abs()
	local a = os.abs('os.lua')
	assert(a:match('.*[\\/](.+)$') == 'os.lua')
end

function test.clock()
	local c1 = os.clock()
	os.sleep(100)
	local c2 = os.clock()
	assert(type(c1) == 'number' and type(c2) == 'number')
	assert(c1 >= 0 and c2-c1 > 0.09999)
end

function test.date()
	assert(type(os.date()) == 'string')

	local t1 = os.date('*t')
	local t2 = os.date('*t', os.time())
	assert(t1.year == t2.year and t1.month == t2.month and t1.day == t2.day)
	assert(math.type(t2.year) == 'integer' and t2.year > 1000 and t2.year < 9999)
	assert(math.type(t2.month) == 'integer' and t2.month >= 1 and t2.month <= 12)
	assert(math.type(t2.day) == 'integer' and t2.day >= 1 and t2.day <= 31)

	assert(string.match(os.date('%Y/%m/%d'), '%d%d%d%d/%d%d/%d%d') ~= nil)
end

function test.difftime()
	local t1 = os.time()
	local x = os.date('*t', t1)
	x.day = x.day + 1
	local t2 = os.time(x)
	local d = os.difftime(t2, t1)
	assert(d >= 86399 and d <= 86401)
end

function test.exists()
	assert(os.exists('test/os.lua'))
	assert(not os.exists('demo/os.lua'))
end

function test.join()
	assert(string.match(os.join('a', 'b', 'c'), 'a[\\/]b[\\/]c') ~= nil)
end

function test.root()
	assert(type(os.root) == 'string')
end

function test.stat()
	local info = os.stat('test/os.lua')
	assert(info.name == 'os.lua')
	assert(not info.isdir)
	assert(math.type(info.size) == 'integer')
	assert(os.difftime(os.time(), info.modtime) > 0)
end

function test.tmpname()
	local t1 = os.tmpname()
	local t2 = os.tmpname()
	assert(t1 ~= t2)
end

return test
