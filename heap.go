package heap

import "iter"

func Concurrent[E any](heap *Heap[E]) *Heap[E] {
	return heap
}

// type Heap[E any] interface {
// 	Len() int
// 	Push(E)
// 	Pop() E
// 	Remove(int) E
// 	Fix(int)
// 	Queue() iter.Seq[E]
// }

// New creates a new heap from the given elements and less function.
// The complexity is O(n) where n = len(e).
func New[E any](e []E, less func(E, E) bool) *Heap[E] {
	h := &Heap[E]{e: e, l: less}
	Init(h)
	return h
}

// Heap is a min-heap of elements of type E.
// It is not safe for concurrent use.
type Heap[E any] struct {
	e []E
	l func(E, E) bool
}

// el returns the element at index i in the heap.
func (h *Heap[E]) el(i int) E {
	return h.e[i]
}

// Len returns the number of elements in the heap.
func (h *Heap[E]) Len() int {
	return len(h.e)
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Push(x E) {
	h.e = append(h.e, x)
	h.up(h.Len() - 1)
}

// less mimics sort.Interface, making code easier to compare / port.
func (h *Heap[E]) less(i, j int) bool {
	return h.l(h.el(i), h.el(j))
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to [Remove](h, 0).
func (h *Heap[E]) Pop() E {
	n := h.Len() - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

// pop removes and returns the last element
// just a helper to make things readable, don't expose
func (h *Heap[E]) pop() E {
	n := h.Len() - 1
	ret := h.el(n)
	h.e = h.e[:n]
	return ret
}

// swap swaps elements
// just a helper to make things readable, don't expose
func (h *Heap[E]) swap(i, j int) {
	h.e[i], h.e[j] = h.el(j), h.el(i)
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func Init[E any](h *Heap[E]) {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Remove(i int) E {
	n := h.Len() - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.pop()
}

// Queue works through the heap in sorted order.
// Be careful about concurrent use.
func (h *Heap[E]) Queue() iter.Seq[E] {
	return func(yield func(E) bool) {
		for h.Len() > 0 {
			if !yield(h.Pop()) {
				break
			}
		}
	}
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling [Remove](h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

func (h *Heap[E]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Heap[E]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}
