local test = {}

function test.table()
	local t = {}
	assert(#t == 0)
	assert(t.test == nil)
	assert(t[1] == nil)
	assert(t['demo'] == nil)

	local t = {'A', 'B', 'C', 'D', 'E'}
	assert(#t == 5)
	assert(t[0] == nil)
	assert(t[1] == 'A')
	assert(t[2] == 'B')
	assert(t[3] == 'C')
	assert(t[4] == 'D')
	assert(t[5] == 'E')
	assert(t[6] == nil)
	t[4] = nil
	t[5] = nil
	assert(#t == 3)

	local t = {
		k1 = 123;
		k2 = 'demo';
		['k3'] = true;
	}
	assert(#t == 0)
	assert(t['k1'] == 123)
	assert(t.k2 == 'demo')
	assert(t.k3 == true)

	local t = {
		k1 = 'vk1';
		[-1] = 'v-1';
		[0] = 'v0';
		k2 = 'vk2';
		[1] = 'v1';
		[2] = 'v2';
		k3 = 'vk3';
		[3] = 'v3';
		[4] = 'v4';
		[5] = 'v5';
		k4 = 'vk4';
	}
	assert(#t == 5)
	assert(t.k1 == 'vk1')
	assert(t.k2 == 'vk2')
	assert(t.k3 == 'vk3')
	assert(t.k4 == 'vk4')
	assert(t[-1] == 'v-1')
	assert(t[0] == 'v0')
	assert(t[1] == 'v1')
	assert(t[2] == 'v2')
	assert(t[3] == 'v3')
	assert(t[4] == 'v4')
	assert(t[5] == 'v5')
	local n = 0
	for k, v in pairs(t) do
		n = n + 1
	end
	assert(n == 11)
end

return test
