// MIT License

// Copyright (c) 2019 soren yang

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package delayqueue

import (
	"container/heap"
	"fmt"
)

// Element for priority queue element
type Element struct {
	// Value for element
	Value interface{}

	// priority of the element, make it private to avoid change from element
	priority int64

	// index of the element in the slice
	index int

	// pq is the refer of priority queue
	pq *priorityQueue
}

// Priority return the priority value of element
func (e *Element) Priority() int64 {
	return e.priority
}

// Index return the element index in slice
func (e *Element) Index() int {
	return e.index
}

// String return the representation of element
func (e *Element) String() string {
	return fmt.Sprintf("*Element{Value:%v,priority:%v,index:%v,pq:%v}", e.Value, e.priority, e.index, e.pq)
}

// PriorityQueue for priority queue trait
type PriorityQueue interface {
	// Add element to the PriorityQueue, it will return the element witch been added
	Add(v interface{}, priority int64) *Element

	// Peek return the lowest priority element
	Peek() *Element

	// Pop return the lowest priority element and remove it
	Pop() *Element

	// Remove will remove the element from the priority queue
	Remove(v *Element) error

	// Update the element in the priority queue with the new priority
	Update(v *Element, priority int64) error

	// Size return the element size of queue
	Size() int
}

// elems is an slice of Elements, it implement the heap.Interface
type elems = []*Element

// heapi implement the heap.Interface
type heapi struct {
	pq *priorityQueue
}

// Len return the size of slice, implement heap.Len
func (h *heapi) Len() int {
	return len(h.pq.e)
}

// Less return the compare of slice element i and j, implement heap.Less
func (h *heapi) Less(i int, j int) bool {
	return h.pq.e[i].priority < h.pq.e[j].priority
}

// Swap change the slice element i and j, implement heap.Swap
func (h *heapi) Swap(i int, j int) {
	h.pq.e[i], h.pq.e[j] = h.pq.e[j], h.pq.e[i]
	h.pq.e[i].index, h.pq.e[j].index = h.pq.e[j].index, h.pq.e[i].index
}

// Push the value at the end of slice, implement heap.Push
func (h *heapi) Push(x interface{}) {
	h.pq.e = append(h.pq.e, x.(*Element))
}

// Pop the value at the last position of slice, implement heap.Pop
func (h *heapi) Pop() interface{} {
	old := h.pq.e
	n := len(old)
	if n == 0 {
		return nil
	}

	// set element to nil for GC
	x := old[n-1]
	old[n-1] = nil
	h.pq.e = old[0 : n-1]
	return x
}

// priorityQueue is a implement by min heap, the 0th element is the lowest value
type priorityQueue struct {
	e elems
	h *heapi
}

// NewPriorityQueue construct a PriorityQueue
func NewPriorityQueue() PriorityQueue {
	pq := &priorityQueue{
		e: make(elems, 0),
	}
	pq.h = &heapi{
		pq: pq,
	}
	return pq
}

// Add element to the PriorityQueue, it will return the element witch been added
func (pq *priorityQueue) Add(x interface{}, priority int64) *Element {
	e := &Element{
		Value:    x,
		priority: priority,
		index:    len(pq.e),
		pq:       pq,
	}
	heap.Push(pq.h, e)
	return e
}

// Peek return the lowest priority element
func (pq *priorityQueue) Peek() *Element {
	if len(pq.e) == 0 {
		return nil
	}

	return pq.e[0]
}

// Pop return the lowest priority element and remove it
func (pq *priorityQueue) Pop() *Element {
	if len(pq.e) == 0 {
		return nil
	}

	x := heap.Pop(pq.h)
	e := x.(*Element)
	e.index = -1
	e.pq = nil
	return e
}

// Remove will remove the element from the priority queue
func (pq *priorityQueue) Remove(e *Element) error {
	if e.pq != pq {
		return fmt.Errorf("PriorityQueue.Remove: QueueMatchFailed: Element[%v], Queue[%v]", e.pq, pq)
	}

	if e.index < 0 || e.index >= len(pq.e) {
		return fmt.Errorf("PriorityQueue.Remove: OutOfIndex: Index[%v], Len[%v]", e.index, len(pq.e))
	}
	if e.priority != pq.e[e.index].priority {
		return fmt.Errorf("PriorityQueue.Remove: PriorityMatchFailed: Element[%v], Queue[%v]", e.priority, pq.e[e.index].priority)
	}

	heap.Remove(pq.h, e.index)
	e.index = -1
	e.pq = nil
	return nil
}

// Update the element in the priority queue with the new priority
func (pq *priorityQueue) Update(e *Element, priority int64) error {
	if e.pq != pq {
		return fmt.Errorf("PriorityQueue.Update: QueueMatchFailed: Element[%v], Queue[%v]", e.pq, pq)
	}
	if e.index < 0 || e.index >= len(pq.e) {
		return fmt.Errorf("PriorityQueue.Update: OutOfIndex: Index[%v], Len[%v]", e.index, len(pq.e))
	}
	if e.priority == priority {
		// the priority doesn't change, just return nil as updated
		return nil
	}

	e.priority = priority
	heap.Fix(pq.h, e.index)
	return nil
}

func (pq *priorityQueue) Size() int {
	return len(pq.e)
}
