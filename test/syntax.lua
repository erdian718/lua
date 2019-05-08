local test = {}

function test.syntax()
	local x = {}
	(x or {}).test = 123
	print(x.test)
end

return test
