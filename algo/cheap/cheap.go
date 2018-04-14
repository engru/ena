// Copyright (c) 2018 soren yang
//
// Licensed under the MIT License
// you may not use this file except in complicance with the License.
// You may obtain a copy of the License at
//
//     https://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cheap extend container/heap for common usage, such as top, pop, push, update and remove. default sort is Min-Heap
package cheap

import "container/heap"

// Interface is defined for cheap usage, Any type that implements this could be use
type Interface interface {
	// Priority for sort
	Priority() uint64
	// Id for identify the each item, use the pointer is ok
	Id() uint64
}

// Heap interface for usage
type Heap interface {
	Top() interface{}
	PopX() interface{}
	PushX(x Interface)
	Update(x Interface)
	Remove(x Interface)
}

// default implement for Heap interface
type defHeap struct {
	array []Interface
	// keep the array.index for each node
	keyMap map[uint64]int
}

// NewHeap cons a Heap object
func NewHeap() Heap {
	h := &defHeap{keyMap: make(map[uint64]int)}
	heap.Init(h)
	return h
}

// Implement sort.Interface.Len
func (h *defHeap) Len() int {
	return len(h.array)
}

// Implement sort.Interface.Less
func (h *defHeap) Less(i int, j int) bool {
	return h.array[i].Priority() < h.array[j].Priority()
}

// Implement sort.Interface.Swap
func (h *defHeap) Swap(i int, j int) {
	// swap inode in array
	h.array[i], h.array[j] = h.array[j], h.array[i]

	// update keyMap
	h.keyMap[h.array[i].Id()] = i
	h.keyMap[h.array[j].Id()] = j
}

// Implement heap.Interface.Push
func (h *defHeap) Push(x interface{}) {
	n, _ := x.(Interface)
	h.keyMap[n.Id()] = len(h.array)
	h.array = append(h.array, n)
}

// Implement heap.Interface.Pop
func (h *defHeap) Pop() interface{} {
	old := h.array
	n := len(old)
	x := old[n-1]

	// set element to nil for GC
	old[n-1] = nil
	h.array = old[0 : n-1]
	delete(h.keyMap, x.Id())
	return x
}

func (h *defHeap) Top() interface{} {
	if h.Len() != 0 {
		return h.array[0]
	}

	return nil
}

func (h *defHeap) PopX() interface{} {
	x := heap.Pop(h)
	return x
}

func (h *defHeap) PushX(x Interface) {
	heap.Push(h, x)
}

func (h *defHeap) Update(n Interface) {
	index, ok := h.keyMap[n.Id()]
	if ok {
		// heap.Remove(h, index)
		// heap.Push(h, n)

		// use heap.Fix
		heap.Fix(h, index)
	}
}

func (h *defHeap) Remove(n Interface) {
	index, ok := h.keyMap[n.Id()]
	if ok {
		heap.Remove(h, index)
	}
}
