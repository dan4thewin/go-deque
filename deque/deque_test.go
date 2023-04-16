package deque

import (
	"reflect"
	"testing"
)

func check(t *testing.T, op func() (int, bool), expected []int, empty bool) {
	for _, en := range expected {
		n, ok := op()
		if !ok {
			t.Error("got false, expected true")
		}
		if n != en {
			t.Errorf("got %v, expected %v", n, en)
		}
	}
	if empty {
		n, ok := op()
		if ok {
			t.Error("got true, expected false")
		}
		if n != 0 {
			t.Errorf("got %v, expected %v", n, 0)
		}
	}
}

func TestPushPop(t *testing.T) {
	d := Deque[int]{}
	d.Push(1, 2, 3)
	check(t, d.Pop, []int{3, 2, 1}, true)
}

func TestUnshiftShift(t *testing.T) {
	d := Deque[int]{}
	d.Unshift(1, 2, 3)
	check(t, d.Shift, []int{3, 2, 1}, true)
}

func TestPushShift(t *testing.T) {
	d := Deque[int]{}
	d.Push(1, 2, 3)
	check(t, d.Shift, []int{1, 2, 3}, true)
}

func TestUnshiftPop(t *testing.T) {
	d := Deque[int]{}
	d.Unshift(1, 2, 3)
	check(t, d.Pop, []int{1, 2, 3}, true)
}

func TestPeek(t *testing.T) {
	d := Deque[int]{}
	d.Push(1, 2, 3)
	// repeat to ensure deque is unchanged by operations
	for range []int{1, 2} {
		check(t, d.Peek, []int{3}, false)
		check(t, d.PeekShift, []int{1}, false)
		if n := d.Len(); n != 3 {
			t.Errorf("length %d, expected %d", n, 3)
		}
	}
}

func TestGrowShrink(t *testing.T) {
	checkcap := func(dd Deque[int], ec int) {
		if c := dd.Cap(); c != ec {
			t.Errorf("capacity %d, expected %d", c, ec)
		}
	}

	// no shrinking
	d := Deque[int]{Minsize: 2}
	d.Push(1, 2)
	checkcap(d, 2)
	d.Push(3, 4)
	checkcap(d, 4)
	d.Push(5)
	checkcap(d, 8)
	check(t, d.Shift, []int{1, 2, 3, 4, 5}, true)
	checkcap(d, 8)

	// shrink and maintain a portion of items
	d = Deque[int]{Minsize: 2, Shrink: ShrinkAt20Pct}
	d.Push(1, 2)
	checkcap(d, 2)
	d.Push(3, 4)
	checkcap(d, 4)
	d.Push(5)
	checkcap(d, 8)
	check(t, d.Shift, []int{1, 2, 3}, false)
	checkcap(d, 8)
	check(t, d.Shift, []int{4}, false)
	checkcap(d, 2)
	check(t, d.Shift, []int{5}, true)

	// shrink without keeping any items
	d = Deque[int]{Minsize: 2, Shrink: ShrinkIfEmpty}
	d.Push(1, 2, 3, 4, 5)
	checkcap(d, 8)
	check(t, d.Shift, []int{1, 2, 3, 4}, false)
	checkcap(d, 8)
	check(t, d.Shift, []int{5}, true)
	checkcap(d, 2)

	// exercise grow with wraparound case, where head is after tail in slice
	d = Deque[int]{Minsize: 4}
	d.Push(1, 2, 3, 4)
	checkcap(d, 4)
	check(t, d.Shift, []int{1, 2}, false)
	d.Push(5)
	checkcap(d, 4)
	d.Push(6, 7, 8)
	checkcap(d, 8)
	check(t, d.Shift, []int{3, 4, 5, 6, 7, 8}, true)
}

func TestSlices(t *testing.T) {
	checkcap := func(dd Deque[int], ec int) {
		if c := dd.Cap(); c != ec {
			t.Errorf("capacity %d, expected %d", c, ec)
		}
	}

	checklen := func(dd Deque[int], el int) {
		if l := dd.Len(); l != el {
			t.Errorf("length %d, expected %d", l, el)
		}
	}

	d := Deque[int]{}
	d.WrapSlice([]int{1, 2, 3, 4})
	checkcap(d, 4)
	checklen(d, 4)
	check(t, d.Pop, []int{4, 3, 2, 1}, true)

	tmp := []int{1, 2, 3, 4}
	d.WrapSlice(tmp[:0])
	checkcap(d, 4)
	checklen(d, 0)

	// exercise ToSlice() with wraparound case
	d.WrapSlice([]int{5, 6, 7, 8})
	check(t, d.Shift, []int{5, 6}, false)
	d.Push(9)
	s := d.ToSlice()
	es := []int{7, 8, 9}
	if !reflect.DeepEqual(s, es) {
		t.Errorf("got %v, expected %v", s, es)
	}

	check(t, d.Pop, []int{9, 8, 7}, true)
	s2 := d.ToSlice()
	if !reflect.DeepEqual(s2, []int{}) {
		t.Errorf("got %v, expected empty slice", s)
	}
}

func TestPushPopString(t *testing.T) {
	d := Deque[string]{}
	d.Push("foo", "bar", "baz")
	for _, en := range []string{"baz", "bar", "foo"} {
		n, ok := d.Pop()
		if !ok {
			t.Errorf("got false, expected true")
		}
		if n != en {
			t.Errorf("got %v, expected %v", n, en)
		}
	}
	n, ok := d.Pop()
	if ok {
		t.Error("got true, expected false")
	}
	if n != "" {
		t.Errorf("got %v, expected empty string", n)
	}
}
