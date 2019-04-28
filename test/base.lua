local test = {}

function test.print()
	assert(print("Hello, Go Lua v" .. _VERSION) == nil)
end

return test
