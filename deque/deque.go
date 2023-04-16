// Copyright 2023 Dan Good. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package deque implements operations on a circular
// double-ended queue backed by a slice.  A deque
// grows by doubling to amortize allocations.
package deque

// Slice size to use when none is specified.
const DefaultSize = 32

// By default, a deque never shrinks.  A deque can be configured
// to shrink back to the defined minimum size either when it
// is empty, or when the length is less than or equal to 20 percent
// of the capacity, provided that length fits within the minimum size.
const (
	ShrinkNever = iota
	ShrinkIfEmpty
	ShrinkAt20Pct
)

// Deque tracks where to enqueue or dequeue for both
// sides of the deque.  A zero-valued deque is usable
// and will allocate on first enqueue.
type Deque[T any] struct {
	Minsize, Shrink int
	head, tail, len int
	dat             []T
}

// A deque changes size by copying into a new slice.
// In the new slice, head is always 0.
func (d *Deque[T]) resize(size int) {
	tmp := make([]T, size)
	if d.len > 0 {
		end := cap(d.dat)
		if end-d.head > d.len {
			end = d.head + d.len
		}
		count := copy(tmp, d.dat[d.head:end])
		if d.len > count {
			copy(tmp[count:], d.dat[:d.len-count])
		}
	}
	d.dat = tmp
	d.head = 0
	d.tail = d.len
	if d.tail == 0 {
		d.tail = size
	}
	d.tail--
}

func (d *Deque[T]) grow(add int) {
	if d.len+add > cap(d.dat) {
		if d.Minsize <= 0 {
			d.Minsize = DefaultSize
		}
		size := cap(d.dat)
		if size == 0 {
			size = d.Minsize
		}
		for size < d.len+add {
			size *= 2
		}
		d.resize(size)
	}
}

func (d *Deque[T]) shrink() {
	if d.Shrink == ShrinkNever {
		return
	}
	if d.len > d.Minsize || cap(d.dat) == d.Minsize {
		return
	}
	if d.Shrink == ShrinkAt20Pct && d.len*5 > cap(d.dat) {
		return
	}
	if d.Shrink == ShrinkIfEmpty && d.len > 0 {
		return
	}
	d.resize(d.Minsize)
}

// Push a single value - only called after grow().
func (d *Deque[T]) push(v T) {
	d.len++
	d.tail++
	if d.tail == cap(d.dat) {
		d.tail = 0
	}
	d.dat[d.tail] = v
}

// Enqueue values onto the end of the deque.
func (d *Deque[T]) Push(v ...T) {
	d.grow(len(v))
	for _, x := range v {
		d.push(x)
	}
}

// Unshift a single value - only called after grow().
func (d *Deque[T]) unshift(v T) {
	d.len++
	if d.head == 0 {
		d.head = cap(d.dat)
	}
	d.head--
	d.dat[d.head] = v
}

// Enqueue values onto the head of the deque.
func (d *Deque[T]) Unshift(v ...T) {
	d.grow(len(v))
	for _, x := range v {
		d.unshift(x)
	}
}

// Return and remove a single value from the end of the deque,
// and optionally shrink.  When empty, return a zero value and false.
func (d *Deque[T]) Pop() (v T, ok bool) {
	if d.len > 0 {
		d.len--
		v, ok = d.dat[d.tail], true
		if d.tail == 0 {
			d.tail = cap(d.dat)
		}
		d.tail--
		d.shrink()
	}
	return
}

// Return and remove a single value from the head of the deque,
// and optionally shrink.  When empty, return a zero value and false.
func (d *Deque[T]) Shift() (v T, ok bool) {
	if d.len > 0 {
		d.len--
		v, ok = d.dat[d.head], true
		d.head++
		if d.head == cap(d.dat) {
			d.head = 0
		}
		d.shrink()
	}
	return
}

// Return a single value from the end of the deque.
// When empty, return a zero value and false.
func (d *Deque[T]) Peek() (v T, ok bool) {
	if d.len > 0 {
		v, ok = d.dat[d.tail], true
	}
	return
}

// Return a single value from the head of the deque.
// When empty, return a zero value and false.
func (d *Deque[T]) PeekShift() (v T, ok bool) {
	if d.len > 0 {
		v, ok = d.dat[d.head], true
	}
	return
}

// Length of the deque
func (d *Deque[T]) Len() int {
	return d.len
}

// Capacity of the deque
func (d *Deque[T]) Cap() int {
	return cap(d.dat)
}

// Return a slice of the deque arranged with head equal to 0.
func (d *Deque[T]) ToSlice() []T {
	if d.len == 0 {
		return []T{}
	}
	if d.head != 0 {
		d.resize(cap(d.dat))
	}
	return d.dat[:d.len]
}

// Use a provided slice as the initial backing store for the deque.
// The next resize() will replace the slice.
func (d *Deque[T]) WrapSlice(dat []T) {
	d.dat = dat[:cap(dat)]
	d.len = len(dat)
	d.head = 0
	d.tail = d.len
	if d.tail == 0 {
		d.tail = cap(d.dat)
	}
	d.tail--
}
