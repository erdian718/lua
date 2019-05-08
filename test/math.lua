local test = {}
local math = require 'math'

local eps = 1e-8

function test.type()
	assert(math.type(0) == 'integer')
	assert(math.type(123) == 'integer')
	assert(math.type(-123) == 'integer')
	assert(math.type(0.0) == 'float')
	assert(math.type(12.34) == 'float')
	assert(math.type(-12.34) == 'float')
	assert(math.type('123') == nil)
end

function test.abs()
	assert(math.abs(0) == 0)
	assert(math.abs(-0) == 0)
	assert(math.abs(-123) == 123)
	assert(math.abs(123) == 123)
	assert(math.abs(-12.34) == 12.34)
	assert(math.abs(12.34) == 12.34)
	assert(math.type(math.abs(0)) == 'integer')
	assert(math.type(math.abs(123)) == 'integer')
	assert(math.type(math.abs(-123)) == 'integer')
	assert(math.type(math.abs(12.34)) == 'float')
	assert(math.type(math.abs(-12.34)) == 'float')
end

function test.number()
	assert(math.huge + 1 <= math.huge)
	assert(-math.huge - 1 >= -math.huge)
	assert(math.maxinteger + 1 == math.mininteger)
	assert(math.mininteger - 1 == math.maxinteger)
	assert(math.abs(math.pi-3.14159265358979323846264338327950288419716939937510582097494459) <= eps)
	assert(math.abs(math.e-2.71828182845904523536028747135266249775724709369995957496696763) <= eps)
end

function test.random()
	local r = math.random()
	assert(0 <= r and r < 1)

	local r = math.random(100)
	assert(1 <= r and r <= 100)
	assert(math.type(r) == 'integer')

	local r = math.random(-100, 100)
	assert(-100 <= r and r <= 100)
	assert(math.type(r) == 'integer')

	local s = math.random(1000000)
	math.randomseed(s)
	local r1 = math.random()
	math.randomseed(s)
	local r2 = math.random()
	assert(r1 == r2)
end

function test.ceil()
	assert(math.ceil(0) == 0)
	assert(math.ceil(123) == 123)
	assert(math.ceil(12.34) == 13)
	assert(math.ceil(-12.34) == -12)
end

function test.floor()
	assert(math.floor(0) == 0)
	assert(math.floor(123) == 123)
	assert(math.floor(12.34) == 12)
	assert(math.floor(-12.34) == -13)
end

function test.deg()
	assert(math.abs(math.deg(0)-0) <= eps)
	assert(math.abs(math.deg(math.pi)-180) <= eps)
	assert(math.abs(math.deg(123)-7047.3808801091) <= eps)
end

function test.rad()
	assert(math.abs(math.rad(0)-0) <= eps)
	assert(math.abs(math.rad(180)-math.pi) <= eps)
	assert(math.abs(math.rad(123)-2.146754979953) <= eps)
end

function test.fmod()
	assert(math.fmod(5, 3) == 2)
	assert(math.fmod(6, 3) == 0)
	assert(math.fmod(-5, 3) == -2)
	assert(math.fmod(-6, 3) == 0)
	assert(math.fmod(5, -3) == 2)
	assert(math.fmod(6, -3) == 0)
	assert(math.fmod(-5, -3) == -2)
	assert(math.fmod(-6, -3) == 0)
	assert(math.abs(math.fmod(math.pi, math.e)-0.42331082513075) <= eps)
end

function test.log()
	assert(math.abs(math.log(math.e)-1) <= eps)
	assert(math.abs(math.log(math.pi)-1.1447298858494) <= eps)
	assert(math.abs(math.log(10, 10)-1) <= eps)
	assert(math.abs(math.log(math.pi, 10)-0.49714987269413) <= eps)
end

function test.atan()
	assert(math.abs(math.atan(123)-1.562666424615) <= eps)
	assert(math.abs(math.atan(123, -456)-2.8781261125549) <= eps)
	assert(math.abs(math.atan(123, 456)-0.26346654103492) <= eps)
	assert(math.abs(math.atan(-123, 456)+0.26346654103492) <= eps)
end

function test.max()
	assert(math.max(5, 2, 8, 2, 1) == 8)
	assert(math.max(8, 2, 5, 2, 1) == 8)
	assert(math.max(5, 2, 1, 2, 8) == 8)
end

function test.min()
	assert(math.min(5, 2, 1, 2, 8) == 1)
	assert(math.min(8, 2, 5, 2, 1) == 1)
	assert(math.min(1, 2, 5, 2, 8) == 1)
end

function test.modf()
	local a, b = math.modf(0)
	assert(a == 0 and b == 0)
	local a, b = math.modf(12.34)
	assert(a == 12 and math.abs(0.34-b) <= eps)
	local a, b = math.modf(-12.34)
	assert(a == -12 and math.abs(-0.34-b) <= eps)
end

function test.tointeger()
	assert(math.tointeger(0) == 0)
	assert(math.tointeger(123) == 123)
	assert(math.tointeger(-123) == -123)
	assert(math.tointeger('123') == 123)
	assert(math.tointeger('-123') == -123)
	assert(math.tointeger('0xF1') == 241)
	assert(math.tointeger(12.34) == nil)
	assert(math.tointeger('12.34') == nil)
	assert(math.tointeger('ABC') == nil)
end

function test.ult()
	assert(not math.ult(0, 0))
	assert(not math.ult(-1, 0))
	assert(math.ult(0, -1))
	assert(math.ult(123, 456))
end

function test.normfuncs()
	assert(math.abs(math.acos(0.123)-1.447484051603) <= eps)
	assert(math.abs(math.asin(0.123)-0.12331227519187) <= eps)
	assert(math.abs(math.cos(123)+0.88796890669186) <= eps)
	assert(math.abs(math.exp(1.23)-3.4212295362897) <= eps)
	assert(math.abs(math.sin(123)+0.45990349068959) <= eps)
	assert(math.abs(math.sqrt(123)-11.090536506409) <= eps)
	assert(math.abs(math.tan(123)-0.51792747158566) <= eps)
end

return test
