/*
	Source: p102-array-util-js

	MIT License

	Copyright (c) 2026 Paul Williams

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

import ArrayUtil from './ArrayUtil.js'

const {
	beforeLast, //
	beforeLastIndex,
	callAll,
	clear,
	// delete,
	findByField,
	insert,
	insertAfter,
	insertBefore,
	itemAfter,
	itemBefore,
	last,
	lastIndex,
	mapToField,
	remove,
	replace,
	withinRange,
} = ArrayUtil

const A = 'A'
const B = 'B'
const C = 'C'
const D = 'D'

const ObjA = { name: 'A', position: 1 }
const ObjB = { name: 'B' }
const ObjC = { name: 'C' }
const ObjD = { name: 'D', position: 4 }

describe('beforeLast()', () => {
	test('returns null for list with 1 item', () => {
		const exp = beforeLast([A])
		expect(exp).toEqual(null)
	})

	test('returns correct item', () => {
		const exp = beforeLast([A, B, C])
		expect(exp).toEqual(B)
	})
})

describe('beforeLastIndex()', () => {
	test('returns -1 for list with 1 item', () => {
		const exp = beforeLastIndex([A])
		expect(exp).toEqual(-1)
	})

	test('returns 2 for list with 4 items', () => {
		const exp = beforeLastIndex([A, B, C, D])
		expect(exp).toEqual(2)
	})
})

test('callAll() calls all functions', () => {
	const called = []
	const calledWith = []

	const fA = (...args) => {
		called.push(A)
		calledWith.push(args)
	}

	const fB = (...args) => {
		called.push(B)
		calledWith.push(args)
	}

	const list = [fA, C, fB, D]
	callAll(list, 'rum', 'whiskey')

	expect(called).toEqual([A, B])
	expect(calledWith).toEqual([
		['rum', 'whiskey'],
		['rum', 'whiskey'],
	])
})

test('clear() removes all items', () => {
	const list = [A, B, C]
	clear(list)
	expect(list).toEqual([])
})

test('findByField() finds the first valid item', () => {
	const A = { id: 1 }
	const B = { id: 2 }
	const C = { id: 3 }

	const list = [A, B, C]
	const result = findByField(list, 'id', 2)
	expect(result).toEqual(B)
})

test('itemBefore()', () => {
	const list = [A, B, C]
	expect(itemBefore(list, A)).toEqual(null)
	expect(itemBefore(list, B)).toEqual(A)
	expect(itemBefore(list, C)).toEqual(B)
	expect(itemBefore(list, D)).toEqual(null)
})

test('itemAfter()', () => {
	const list = [A, B, C]
	expect(itemAfter(list, A)).toEqual(B)
	expect(itemAfter(list, B)).toEqual(C)
	expect(itemAfter(list, C)).toEqual(null)
	expect(itemAfter(list, D)).toEqual(null)
})

describe('insert()', () => {
	test('puts item in correct place', () => {
		const list = [A, C]
		insert(list, 1, B)
		expect(list).toEqual([A, B, C])
	})

	test('puts item at end of list', () => {
		const list = [A, B]
		insert(list, 2, C)
		expect(list).toEqual([A, B, C])
	})

	test('throws if index is out of bounds', () => {
		const f = () => insert([A, C], 5, B)
		expect(f).toThrow(Error)
	})
})

describe('insertAfter()', () => {
	test('puts item in correct place', () => {
		const list = [A, C]
		insertAfter(list, A, B)
		expect(list).toEqual([A, B, C])
	})

	test('throws if ref item not in list', () => {
		const f = () => insertAfter([A, C], D, B)
		expect(f).toThrow(Error)
	})
})

describe('insertBefore()', () => {
	test('puts item in correct place', () => {
		const list = [A, C]
		insertBefore(list, C, B)
		expect(list).toEqual([A, B, C])
	})

	test('throws if ref item not in list', () => {
		const f = () => insertBefore([A, C], D, B)
		expect(f).toThrow(Error)
	})
})

describe('last()', () => {
	test('returns null for empty list', () => {
		const exp = last([])
		expect(exp).toEqual(null)
	})

	test('returns correct item', () => {
		const exp = last([A, B, C])
		expect(exp).toEqual(C)
	})
})

describe('mapToField()', () => {
	test('returns all mapped values', () => {
		const given = [ObjA, ObjB, ObjC, ObjD]
		const exp = mapToField(given, 'name')
		expect(exp).toEqual(['A', 'B', 'C', 'D'])
	})

	test('returns only own property values', () => {
		const given = [ObjA, ObjB, ObjC, ObjD]
		const exp = mapToField(given, 'position')
		expect(exp).toEqual([1, 4])
	})

	test('ignores non-objects', () => {
		const given = [ObjA, ObjB, 3, ['Four', 'Five']]
		const exp = mapToField(given, 'name')
		expect(exp).toEqual(['A', 'B'])
	})
})

describe('remove()', () => {
	test('remove correct item', () => {
		const list = [A, B, C]
		remove(list, B)
		expect(list).toEqual([A, C])
	})

	test('remove nothing when item not in list', () => {
		const list = [A, B, C]
		remove(list, D)
		expect(list).toEqual([A, B, C])
	})
})

describe('replace()', () => {
	test('swaps correct items', () => {
		const list = [A, B, C]
		replace(list, C, D)
		expect(list).toEqual([A, B, D])
	})

	test('throws if current item is not in list', () => {
		const f = () => replace([A, B], D, C)
		expect(f).toThrow(Error)
	})
})

describe('withinRange()', () => {
	test('with length excluded', () => {
		const f = (i) => withinRange([A, B, C], i)

		expect(f(-1)).toEqual(false)
		expect(f(0)).toEqual(true)
		expect(f(1)).toEqual(true)
		expect(f(2)).toEqual(true)
		expect(f(3)).toEqual(false)
	})

	test('with length included', () => {
		const f = (i) => withinRange([A, B, C], i, true)

		expect(f(-1)).toEqual(false)
		expect(f(0)).toEqual(true)
		expect(f(1)).toEqual(true)
		expect(f(2)).toEqual(true)
		expect(f(3)).toEqual(true)
		expect(f(4)).toEqual(false)
	})
})
