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

package algo

import (
	"container/heap"
)

type Item struct {
	value    string
	priority int
	index    int
}

type PriorityQueue []*Item

func (q PriorityQueue) Len() int {
	return len(q)
}

func (q PriorityQueue) Less(i int, j int) bool {
	return q[i].priority > q[i].priority
}

func (q PriorityQueue) Swap(i int, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *PriorityQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*Item)
	item.index = n
	*q = append(*q, item)
}

func (q *PriorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.index = -1
	*q = old[0 : n-1]
	return item
}

func (q *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(q, item.index)
}
