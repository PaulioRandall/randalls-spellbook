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

function isObject(v) {
	return (
		typeof v === 'object' && //
		!Array.isArray(v) &&
		v !== null
	)
}

function err(msg) {
	return new Error(`[ArrayUtil] ${msg}`)
}

// beforeLast returns the item second from last, or null
// if no such item exists.
function beforeLast(array) {
	const i = beforeLastIndex(array)
	return i < 0 ? null : array[i]
}

// beforeLastIndex returns the index of the second from
// last item. This will be negative if $array contains
// less than 2 items.
function beforeLastIndex(array) {
	return array.length - 2
}

// callAll invokes all items that satisfy
// `typeof item === 'function'` passing $args on each call.
function callAll(array, ...args) {
	array.forEach((item) => {
		if (typeof item === 'function') {
			item(...args)
		}
	})
}

// clear removes all items.
function clear(array) {
	array.splice(0)
}

// findByField returns the first matching item where the
// value of $field matches $value. Non-objects items are
// skipped during search. Null is returned if no matching
// item is found.
function findByField(array, field, value) {
	for (const item of array) {
		if (!isObject(item)) {
			continue
		}

		if (item[field] === value) {
			return item
		}
	}

	return null
}

// insert inserts $item at the $index. Throws an error if
// $index is out of bounds. Returns $item.
function insert(array, index, item) {
	if (!withinRange(array, index, true)) {
		throw err('Index is out of range')
	}

	array.splice(index, 0, item)
	return item
}

// insertAfter inserts $item in the slot after $refItem.
// Throws an error if $refItem is not found. Returns $item.
function insertAfter(array, refItem, item) {
	const i = array.indexOf(refItem)

	if (i < 0) {
		throw err("Reference item doesn't exist")
	}

	array.splice(i + 1, 0, item)
	return item
}

// insertBefore inserts $item in the slot before $refItem.
// Throws an error if $refItem is not found. Returns $item.
function insertBefore(array, refItem, item) {
	const i = array.indexOf(refItem)

	if (i < 0) {
		throw err("Reference item doesn't exist")
	}

	array.splice(i, 0, item)
}

// itemAfter returns the item after $refItem. Returns null
// if $refItem is not found or is the last item.
function itemAfter(array, refItem) {
	const i = array.indexOf(refItem)
	const lastIdx = lastIndex(array)
	return i < 0 || i >= lastIdx ? null : array[i + 1]
}

// itemBefore returns the item before $refItem. Returns
// null if $refItem is not found or is the first item.
function itemBefore(array, refItem) {
	const i = array.indexOf(refItem)
	return i <= 0 ? null : array[i - 1]
}

// last returns the last item in $array.
function last(array) {
	const i = lastIndex(array)
	return i < 0 ? null : array[i]
}

// lastIndex returns the index of the last item in $array.
function lastIndex(array) {
	return array.length - 1
}

// mapToField is a mapper returning the value of $field for
// all objects with a $field as an own property.
// Non-objects and objects without a $field property are
// skipped. This means the resultant array will be equal to
// or smaller than the $array, it may even have a length of
// zero.
function mapToField(array, field) {
	const result = []

	for (const item of array) {
		if (!isObject(item)) {
			continue
		}

		if (Object.hasOwn(item, field)) {
			result.push(item[field])
		}
	}

	return result
}

// remove deletes $item from $array, if it exists. $item is
// returned.
function remove(array, item) {
	const i = array.indexOf(item)

	if (i > -1) {
		array.splice(i, 1)
	}

	return item
}

// replace swaps the $currentItem with $newItem. An error
// is thrown if $currentItem is not found. Returns $newItem.
function replace(array, currentItem, newItem) {
	const i = array.indexOf(currentItem)

	if (i < 0) {
		throw err("Current item doesn't exist")
	}

	array.splice(i, 1, newItem)
	return newItem
}

// withinRange returns true if $index is within the bounds
// $array, i.e. references a valid array slot. If
// $includeLength is true, $array's length is considered a
// valid index.
function withinRange(array, index, includeLength = false) {
	return (
		(index >= 0 && index < array.length) || //
		(includeLength && index === array.length) //
	)
}

export default {
	beforeLast, //
	beforeLastIndex,
	callAll,
	clear,
	delete: remove,
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
}
