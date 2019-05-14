/*
Copyright 2019 by ofunc

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package lua

import (
	"math"
	"reflect"
)

// table is the VM's table type.
type table struct {
	array []value
	hash  map[value]value
	base  int
	nums  []int
	meta  *table
}

func newTable(l *State, as, hs int) *table {
	t := new(table)
	if as > 0 {
		t.array = make([]value, as)
	}
	if hs > 0 {
		t.hash = make(map[value]value, hs)
	}
	t.base = 1

	for {
		u := t.base << 1
		if u <= 0 || u > as {
			break
		}
		t.base = u
	}
	return t
}

// Length returns the raw table length as would be returned by the length operator.
func (tbl *table) Length() int {
	return tbl.count(0, len(tbl.nums))
}

// Count returns the raw table pairs number.
func (tbl *table) Count() int {
	return tbl.count(0, 1) + len(tbl.hash)
}

// Get reads the value at index k from the table without using any meta methods.
func (tbl *table) Get(k value) value {
	switch idx := k.(type) {
	case int64:
		if idx >= 1 && idx <= int64(len(tbl.array)) {
			return tbl.array[idx-1]
		}
	case float64:
		if math.IsNaN(idx) {
			return nil
		}
		if i := int64(idx); float64(i) == idx {
			if i >= 1 && i <= int64(len(tbl.array)) {
				return tbl.array[i-1]
			}
		}
	case nil:
		return nil
	}
	return tbl.hash[k]
}

// Set sets a key k in the table to the value v without using any meta methods.
func (tbl *table) Set(k, v value) {
	switch idx := k.(type) {
	case int64:
		if idx >= 1 {
			tbl.seti(idx, v)
			return
		}
	case float64:
		if math.IsNaN(idx) {
			panic("cannot set nan table index")
		}
		if i := int64(idx); float64(i) == idx {
			k = i
			if i >= 1 {
				tbl.seti(i, v)
				return
			}
		}
	case nil:
		panic("cannot set nil table index")
	}

	if v == nil {
		delete(tbl.hash, k)
	} else {
		if tbl.hash == nil {
			tbl.hash = make(map[value]value)
		}
		tbl.hash[k] = v
	}
}

// GetIter returns the table iterator.
func (tbl *table) GetIter() func() (value, value) {
	array := tbl.array
	hash := reflect.ValueOf(tbl.hash).MapRange()
	i, n := 0, len(array)
	return func() (value, value) {
		for ; i < n; i++ {
			if v := array[i]; v != nil {
				i++
				return int64(i), v
			}
		}
		if hash.Next() {
			k := hash.Key().Interface()
			v := hash.Value().Interface()
			return k, v
		}
		return nil, nil
	}
}

func (tbl *table) seti(i int64, v value) {
	n := int64(len(tbl.array))
	if i <= n {
		x := tbl.array[i-1]
		if x != v {
			if v == nil {
				tbl.nums[tbl.index(i)] -= 1
			} else if x == nil {
				tbl.nums[tbl.index(i)] += 1
			}
			tbl.array[i-1] = v
		}
	} else {
		var x value
		if tbl.hash != nil {
			x = tbl.hash[i]
		}
		if x != v {
			if v == nil {
				tbl.nums[tbl.index(i)] -= 1
				delete(tbl.hash, i)
			} else if x == nil {
				tbl.nums[tbl.index(i)] += 1
				if tbl.extend(i) {
					tbl.array[i-1] = v
				} else {
					if tbl.hash == nil {
						tbl.hash = make(map[value]value)
					}
					tbl.hash[i] = v
				}
			} else {
				tbl.hash[i] = v
			}
		}
	}
}

func (tbl *table) extend(idx int64) bool {
	m := len(tbl.nums) - 1
	u := tbl.base << uint(m)
	if u <= 0 {
		return false
	}
	for s := tbl.Length(); m >= 1 && int64(u) >= idx; m-- {
		u >>= 1
		if s >= u {
			break
		}
		s -= tbl.nums[m]
	}
	u <<= 1
	if m < 1 || int64(u) < idx {
		return false
	}
	tbl.base = u
	array := make([]value, u)
	copy(array, tbl.array)
	tbl.array = array
	for k, v := range tbl.hash {
		if i, ok := k.(int64); ok && i >= 1 && i <= int64(u) {
			array[i-1] = v
			delete(tbl.hash, k)
		}
	}
	nums := make([]int, len(tbl.nums)-m)
	nums[0] = tbl.count(0, m+1)
	copy(nums[1:], tbl.nums[1+m:])
	tbl.nums = nums
	return true
}

func (tbl *table) count(i, j int) int {
	n := len(tbl.nums)
	if j > n {
		j = n
	}
	s := 0
	for k := i; k < j; k++ {
		s += tbl.nums[k]
	}
	return s
}

func (tbl *table) index(i int64) int {
	m := 0
	if int64(len(tbl.array)) < i {
		u := tbl.base
		for {
			m++
			u <<= 1
			if u <= 0 || int64(u) >= i {
				break
			}
		}
	}
	if m >= len(tbl.nums) {
		nums := make([]int, m+1)
		copy(nums, tbl.nums)
		tbl.nums = nums
	}
	return m
}
