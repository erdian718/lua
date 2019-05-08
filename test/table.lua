local test = {}
local table = require 'table'

function test.concat()
	assert(table.concat({}, '-') == '')
	assert(table.concat({'A', 'B', 'C'}, '-') == 'A-B-C')
	assert(table.concat({111, 222, 333}, '-') == '111-222-333')
	assert(table.concat({'A', 'B', 'C'}, '-', 2, 3) == 'B-C')
	assert(table.concat({'A', 'B', 'C'}, '-', 3, 3) == 'C')
	assert(table.concat({'A', 'B', 'C'}, '-', 3, 2) == '')
end

function test.insert()
	local xs = {}
	table.insert(xs, 'A')
	assert(#xs == 1)
	assert(xs[1] == 'A')

	table.insert(xs, 'B')
	assert(#xs == 2)
	assert(xs[1] == 'A' and xs[2] == 'B')

	table.insert(xs, 1, 'C')
	assert(#xs == 3)
	assert(xs[1] == 'C' and xs[2] == 'A' and xs[3] == 'B')

	table.insert(xs, 2, 'D')
	assert(#xs == 4)
	assert(xs[1] == 'C' and xs[2] == 'D' and xs[3] == 'A' and xs[4] == 'B')
end

function test.move()
	local xs = {'A', 'B', 'C', 'D'}
	local ys = {111, 222, 333}
	table.move(xs, 1, 4, 1, ys)
	assert(#ys == 4)
	assert(ys[1] == 'A' and ys[2] == 'B' and ys[3] == 'C' and ys[4] == 'D')

	local xs = {'A', 'B', 'C', 'D'}
	local ys = {111, 222, 333}
	table.move(xs, 2, 4, 2, ys)
	assert(#ys == 4)
	assert(ys[1] == 111 and ys[2] == 'B' and ys[3] == 'C' and ys[4] == 'D')

	local xs = {'A', 'B', 'C', 'D'}
	table.move(xs, 1, 4, 2)
	assert(#xs == 5)
	assert(xs[1] == 'A' and xs[2] == 'A' and xs[3] == 'B' and xs[4] == 'C' and xs[5] == 'D')
end

function test.pack()
	local xs = table.pack()
	assert(xs.n == 0)

	local xs = table.pack('A')
	assert(xs.n == 1)
	assert(xs[1] == 'A')

	local xs = table.pack('A', 'B', 'C')
	assert(xs.n == 3)
	assert(xs[1] == 'A' and xs[2] == 'B' and xs[3] == 'C')

	local xs = table.pack('A', nil, 'C')
	assert(xs.n == 3)
	assert(xs[1] == 'A' and xs[2] == nil and xs[3] == 'C')
end

function test.remove()
	local xs = {'A', 'B', 'C', 'D'}
	table.remove(xs)
	assert(#xs == 3)
	assert(xs[1] == 'A' and xs[2] == 'B' and xs[3] == 'C')

	table.remove(xs, 2)
	assert(#xs == 2)
	assert(xs[1] == 'A' and xs[2] == 'C')

	table.remove(xs, 1)
	assert(#xs == 1)
	assert(xs[1] == 'C')
end

function test.sort()
	local xs = {1, 2, 5, 4, 3}
	table.sort(xs)
	assert(#xs == 5)
	assert(xs[1] == 1 and xs[2] == 2 and xs[3] == 3 and xs[4] == 4 and xs[5] == 5)

	local xs = {1, 2, 5, 4, 3}
	table.sort(xs, function(a, b)
		return a > b
	end)
	assert(#xs == 5)
	assert(xs[1] == 5 and xs[2] == 4 and xs[3] == 3 and xs[4] == 2 and xs[5] == 1)
end

function test.unpack()
	assert(select('#', table.unpack({})) == 0)

	assert(select('#', table.unpack({'A', 'B', 'C'})) == 3)
	local a, b, c = table.unpack({'A', 'B', 'C'})
	assert(a == 'A' and b == 'B' and c == 'C')

	assert(select('#', table.unpack({'A', 'B', 'C', 'D'}, 2, 3)) == 2)
	local b, c = table.unpack({'A', 'B', 'C', 'D'}, 2, 3)
	assert(b == 'B' and c == 'C')
end

return test
