local test = {}
local string = require 'string'
local os = require 'os'

function test.bytechar()
	local a = string.byte('ABC')
	assert(string.char(a) == 'A')

	local b = string.byte('ABC', 2)
	assert(string.char(b) == 'B')

	local a, b, c = string.byte('ABC', 1, -1)
	assert(string.char(a, b, c) == 'ABC')

	local b, c = string.byte('ABC', 2, -1)
	assert(string.char(b, c) == 'BC')
end

function test.find()
	local i, j = string.find('Hello World', 'World')
	assert(i == 7 and j == 11)

	local i, j = string.find('Hello World', 'World', 3)
	assert(i == 7 and j == 11)

	local i, j = string.find('Hello World', 'World', 12)
	assert(i == nil and j == nil)

	local i, j = string.find('Hello World', '%a+')
	assert(i == 1 and j == 5)

	local i, j = string.find('Hello World', '%a+', 3)
	assert(i == 3 and j == 5)

	local i, j = string.find('Hello World', '%a+', 1, true)
	assert(i == nil and j == nil)

	local i, j = string.find('Hello %a+ World', '%a+', 1, true)
	assert(i == 7 and j == 9)
end

function test.format()
	assert(string.format('%b', 2) == '10')
	assert(string.format('0x%x', 15) == '0xf')
	assert(string.format('%d', 0xf) == '15')
	assert(string.format('%.2f', 3.1415926) == '3.14')
	assert(string.format('%v\t%v\t%v', 1, 2, 3) == '1\t2\t3')
end

function test.gmatch()
	local t = {}
	local s = 'from=world, to=Lua'
	for k, v in string.gmatch(s, '(%w+)=(%w+)') do
		t[k] = v
	end
	assert(t.from == 'world')
	assert(t.to == 'Lua')
end

function test.gsub()
	assert(string.gsub('hello world', '(%w+)', '%1 %1') == 'hello hello world world')
	assert(string.gsub('hello world', '%w+', '%0 %0', 1) == 'hello hello world')
	assert(string.gsub('hello world from Lua', '(%w+)%s*(%w+)', '%2 %1') == 'world hello Lua from')

	local t = {name='lua', version='5.3'}
	assert(string.gsub('$name-$version.tar.gz', '%$(%w+)', t) == 'lua-5.3.tar.gz')
end

function test.len()
	assert(string.len('') == 0)
	assert(string.len('ABC') == 3)
end

function test.lower()
	assert(string.lower('AbC') == 'abc')
end

function test.match()
	assert(string.match('a.b.c.txt', '.+%.(%w+)$') == 'txt')
end

function test.rep()
	assert(string.rep('', 0) == '')
	assert(string.rep('', -1) == '')
	assert(string.rep('*', 1) == '*')
	assert(string.rep('*', 3) == '***')
	assert(string.rep('*', 3, '-') == '*-*-*')
	assert(string.rep('*', 0, '-') == '')
end

function test.reverse()
	assert(string.reverse('') == '')
	assert(string.reverse('ABC') == 'CBA')
end

function test.sub()
	assert(string.sub('', 1) == '')
	assert(string.sub('ABC', 1) == 'ABC')
	assert(string.sub('ABC', 2) == 'BC')
	assert(string.sub('ABC', 3) == 'C')
	assert(string.sub('ABC', 4) == '')
	assert(string.sub('ABC', 5) == '')
end

function test.upper()
	assert(string.upper('aBc') == 'ABC')
end

return test
