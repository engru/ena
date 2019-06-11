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
	"math/rand"
	"testing"

	"github.com/stretchr/testify/suite"
)

type priorityQueueTestSuite struct {
	suite.Suite

	pq *priorityQueue
}

func (s *priorityQueueTestSuite) SetupTest() {
	s.pq = NewPriorityQueue(0).(*priorityQueue)
}

func (s *priorityQueueTestSuite) TestAddOk() {
	type testCase struct {
		description string
		value       int
		priority    int64
		index       int
	}
	testCases := []testCase{
		{
			description: "add priority 1",
			value:       100,
			priority:    1,
			index:       0,
		},
		{
			description: "add priority 2",
			value:       200,
			priority:    2,
			index:       1,
		},
		{
			description: "add priority 0",
			value:       300,
			priority:    0,
			index:       0,
		},
		{
			description: "add priority 100",
			value:       400,
			priority:    100,
			index:       3,
		},
		{
			description: "add priority 1",
			value:       500,
			priority:    1,
			index:       1,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.description, func() {
			e := s.pq.Add(tc.value, tc.priority)
			s.Equal(tc.value, e.Value.(int))
			s.Equal(tc.index, e.index)
			s.Equal(tc.priority, e.Priority())
		})
	}

	testCases2 := []testCase{
		{
			description: "index 0",
			index:       0,
			priority:    0,
			value:       300,
		},
		{
			description: "index 1",
			index:       1,
			priority:    1,
			value:       500,
		},
		{
			description: "index 2",
			index:       2,
			priority:    1,
			value:       100,
		},
		{
			description: "index 3",
			index:       3,
			priority:    100,
			value:       400,
		},
		{
			description: "index 4",
			index:       4,
			priority:    2,
			value:       200,
		},
	}
	for i, e := range s.pq.e {
		tc := testCases2[i]
		s.Run(tc.description, func() {
			s.Equal(i, e.index)
			s.Equal(tc.index, e.index)
			s.Equal(tc.priority, e.priority)
			s.Equal(tc.value, e.Value.(int))
		})
		s.Equal(i, e.index)
	}
}

func (s *priorityQueueTestSuite) TestPeekNil() {
	s.Nil(s.pq.Peek())
}

func (s *priorityQueueTestSuite) TestPeekOk() {
	type testCase struct {
		description string
		value       int
		priority    int64
		index       int

		targetValue    int
		targetPriority int64
	}
	testCases := []testCase{
		{
			description:    "add priority 1",
			value:          100,
			priority:       1,
			index:          0,
			targetValue:    100,
			targetPriority: 1,
		},
		{
			description:    "add priority 2",
			value:          200,
			priority:       2,
			index:          1,
			targetValue:    100,
			targetPriority: 1,
		},
		{
			description:    "add priority 0",
			value:          300,
			priority:       0,
			index:          0,
			targetValue:    300,
			targetPriority: 0,
		},
		{
			description:    "add priority 100",
			value:          400,
			priority:       100,
			index:          3,
			targetValue:    300,
			targetPriority: 0,
		},
		{
			description:    "add priority -1",
			value:          500,
			priority:       -1,
			index:          0,
			targetValue:    500,
			targetPriority: -1,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.description, func() {
			e := s.pq.Add(tc.value, tc.priority)
			s.Equal(tc.value, e.Value.(int))
			s.Equal(tc.index, e.index)
			s.Equal(tc.priority, e.Priority())

			e = s.pq.Peek()
			s.NotNil(e)
			s.Equal(tc.targetValue, e.Value.(int))
			s.Equal(tc.targetPriority, e.Priority())
		})
	}
}

func (s *priorityQueueTestSuite) TestPopNil() {
	s.Nil(s.pq.Pop())
}

func (s *priorityQueueTestSuite) TestPopOk() {
	type testCase struct {
		description string
		value       int
		priority    int64
		index       int
	}
	testCases := []testCase{
		{
			description: "add priority 1",
			value:       100,
			priority:    1,
			index:       0,
		},
		{
			description: "add priority 2",
			value:       200,
			priority:    2,
			index:       1,
		},
		{
			description: "add priority 0",
			value:       300,
			priority:    0,
			index:       0,
		},
		{
			description: "add priority 100",
			value:       400,
			priority:    100,
			index:       3,
		},
		{
			description: "add priority -1",
			value:       500,
			priority:    -1,
			index:       0,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.description, func() {
			e := s.pq.Add(tc.value, tc.priority)
			s.Equal(tc.value, e.Value.(int))
			s.Equal(tc.index, e.index)
			s.Equal(tc.priority, e.Priority())
		})
	}

	type popTarget struct {
		priority int64
		value    int
	}
	targets := []popTarget{
		{
			priority: -1,
			value:    500,
		},
		{
			priority: 0,
			value:    300,
		},
		{
			priority: 1,
			value:    100,
		},
		{
			priority: 2,
			value:    200,
		},
		{
			priority: 100,
			value:    400,
		},
	}
	for _, tc := range targets {
		e := s.pq.Pop()
		s.NotNil(e)

		s.Equal(tc.priority, e.Priority())
		s.Equal(tc.value, e.Value.(int))
	}
}

func (s *priorityQueueTestSuite) TestRemoveQueueMatchFailed() {
	e := s.pq.Add(1, 1)
	s.NotNil(e)
	s.Equal(s.pq, e.pq)

	e2 := &Element{
		pq: nil,
	}
	err := s.pq.Remove(e2)
	s.EqualError(fmt.Errorf("PriorityQueue.Remove: QueueMatchFailed: Element[%v], Queue[%v]", e2.pq, s.pq), err.Error())
}

func (s *priorityQueueTestSuite) TestRemoveOutOfIndexFailed() {
	e := s.pq.Add(1, 1)
	s.NotNil(e)
	s.Equal(s.pq, e.pq)

	e.index = -1
	err := s.pq.Remove(e)
	s.EqualError(fmt.Errorf("PriorityQueue.Remove: OutOfIndex: Index[%v], Len[%v]", e.index, len(s.pq.e)), err.Error())

	e.index = 2
	err = s.pq.Remove(e)
	s.EqualError(fmt.Errorf("PriorityQueue.Remove: OutOfIndex: Index[%v], Len[%v]", e.index, len(s.pq.e)), err.Error())
}

func (s *priorityQueueTestSuite) TestRemovePriorityMatchFailed() {
	e := s.pq.Add(1, 1)
	s.NotNil(e)
	s.Equal(s.pq, e.pq)

	e2 := &Element{
		Value:    e.Value,
		priority: e.priority + 1,
		index:    e.index,
		pq:       e.pq,
	}

	err := s.pq.Remove(e2)
	s.EqualError(fmt.Errorf("PriorityQueue.Remove: PriorityMatchFailed: Element[%v], Queue[%v]", e2.priority, s.pq.e[e2.index].priority), err.Error())
}

func (s *priorityQueueTestSuite) TestRemoveOk() {
	s.initPriorityQueue()

	for s.pq.Size() > 0 {
		// random pick one element to remove
		i := rand.Intn(s.pq.Size())
		e := s.pq.e[i]
		err := s.pq.Remove(e)
		s.NoError(err)

		for _, v := range s.pq.e {
			if e == v {
				s.Fail("expect element[%v] been removed, found in queue", e)
			}
		}

		s.assertIsPriorityQueue()
	}
}

func (s *priorityQueueTestSuite) TestUpdateQueueMatchFailed() {
	e := s.pq.Add(1, 1)
	s.NotNil(e)
	s.Equal(s.pq, e.pq)

	e2 := &Element{
		pq: nil,
	}
	err := s.pq.Update(e2, e.priority+1)
	s.EqualError(fmt.Errorf("PriorityQueue.Update: QueueMatchFailed: Element[%v], Queue[%v]", e2.pq, s.pq), err.Error())
}

func (s *priorityQueueTestSuite) TestUpdateOutOfIndexFailed() {
	e := s.pq.Add(1, 1)
	s.NotNil(e)
	s.Equal(s.pq, e.pq)

	e.index = -1
	err := s.pq.Update(e, e.priority+1)
	s.EqualError(fmt.Errorf("PriorityQueue.Update: OutOfIndex: Index[%v], Len[%v]", e.index, len(s.pq.e)), err.Error())

	e.index = 2
	err = s.pq.Update(e, e.priority+1)
	s.EqualError(fmt.Errorf("PriorityQueue.Update: OutOfIndex: Index[%v], Len[%v]", e.index, len(s.pq.e)), err.Error())
}

func (s *priorityQueueTestSuite) TestUpdateOk() {
	s.initPriorityQueue()

	c := 100
	for c > 0 {
		c--
		// random pick one element to remove
		i := rand.Intn(s.pq.Size())
		e := s.pq.e[i]
		err := s.pq.Update(e, e.priority+1)
		s.NoError(err)

		s.assertIsPriorityQueue()
	}
}

func (s *priorityQueueTestSuite) TestUpdatePriorityNotChange() {
	e := s.pq.Add(1, 1)
	s.NotNil(e)
	s.Equal(s.pq, e.pq)

	err := s.pq.Update(e, e.priority)
	s.NoError(err)
}

func (s *priorityQueueTestSuite) TestSize() {
	s.Equal(0, s.pq.Size())

	type testCase struct {
		description string
		value       int
		priority    int64
		index       int
	}
	testCases := []testCase{
		{
			description: "add priority 1",
			value:       100,
			priority:    1,
			index:       0,
		},
		{
			description: "add priority 2",
			value:       200,
			priority:    2,
			index:       1,
		},
		{
			description: "add priority 0",
			value:       300,
			priority:    0,
			index:       0,
		},
		{
			description: "add priority 100",
			value:       400,
			priority:    100,
			index:       3,
		},
		{
			description: "add priority -1",
			value:       500,
			priority:    -1,
			index:       0,
		},
	}
	for i, tc := range testCases {
		s.Run(tc.description, func() {
			e := s.pq.Add(tc.value, tc.priority)
			s.Equal(tc.value, e.Value.(int))
			s.Equal(tc.index, e.index)
			s.Equal(tc.priority, e.Priority())

			s.Equal(i+1, s.pq.Size())
		})
	}
}

func (s *priorityQueueTestSuite) assertIsPriorityQueue() {
	for i, e := range s.pq.e {
		if i != e.index {
			s.Fail("expect element[%v] index[%v] equal index of slice, failed", e, i)
		}
	}

	es := make(elems, len(s.pq.e))
	copy(es, s.pq.e)

	heap.Init(s.pq.h)

	for i := range es {
		if es[i] != s.pq.e[i] {
			s.Fail("expect element[%v] at index[%v], got [%v]", es[i], i, s.pq.e[i])
		}
	}
}

func (s *priorityQueueTestSuite) initPriorityQueue() {
	type testCase struct {
		description string
		value       int
		priority    int64
		index       int
	}
	testCases := []testCase{
		{
			description: "add priority 1",
			value:       100,
			priority:    1,
			index:       0,
		},
		{
			description: "add priority 2",
			value:       200,
			priority:    2,
			index:       1,
		},
		{
			description: "add priority 0",
			value:       300,
			priority:    0,
			index:       0,
		},
		{
			description: "add priority 100",
			value:       400,
			priority:    100,
			index:       3,
		},
		{
			description: "add priority -1",
			value:       500,
			priority:    -1,
			index:       0,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.description, func() {
			e := s.pq.Add(tc.value, tc.priority)
			s.Equal(tc.value, e.Value.(int))
			s.Equal(tc.index, e.index)
			s.Equal(tc.priority, e.Priority())
		})
	}
}

func TestPriorityQueueTestSuite(t *testing.T) {
	s := &priorityQueueTestSuite{}
	suite.Run(t, s)
}
