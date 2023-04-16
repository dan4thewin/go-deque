# go-deque
## what
The deque package (pronounced "deck") provides a generic golang implementation
of a circular double-ended queue backed by a slice.  Storage grows by doubling
and optionally shrinks.  Operations include Push, Pop, Peek, Shift, Unshift,
and PeekShift.

## why
This implementation aligns with my
[C deque library in CCAN](https://ccodearchive.net/info/deque.html).
I find it easier to reason about allocation and garbage generation with this
amortized approach.  At steady state, this deque needs no additional allocation,
and generates no garbage.

## who
Dan Good
