local test = {}
local io = require 'io'

function test.stdio()
	assert(io.stdin ~= nil)
	assert(io.stdout ~= nil)
	assert(io.stderr ~= nil)
	assert(io.type(io.stdin) == 'readwriter')
	assert(io.type(io.stdout) == 'readwriter')
	assert(io.type(io.stderr) == 'readwriter')
end

return test
