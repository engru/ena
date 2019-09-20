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
	"context"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type delayQueueTestSuite struct {
	suite.Suite

	dq *delayQueue
}

func (s *delayQueueTestSuite) SetupTest() {
	poll = pollImpl
	s.dq = New(1).(*delayQueue)
}

func (s *delayQueueTestSuite) TestOffer() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer func() {
			s.T().Logf("Poll done")
			wg.Done()
		}()

		s.dq.Poll(ctx)
	}()
	go func() {
		defer func() {
			s.T().Logf("Receive done")
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case v := <-s.dq.Chan():
				s.T().Logf("Receive %v", v.(int))
			}
		}
	}()

	time.Sleep(10 * time.Millisecond)
	s.Equal(int32(1), atomic.LoadInt32(&s.dq.sleeping))

	n := defaultTimer.Now()
	n += 1000
	n += 10
	s.T().Logf("Offer first element: priority %v", n)
	s.dq.Offer(n, n)
	s.Equal(1, s.dq.Size())
	// TODO(yangsonglin): difficult to test, we expect to see 0
	v := atomic.LoadInt32(&s.dq.sleeping)
	if v != 0 && v != 1 {
		s.FailNow("sleeping[%v] doesn't be 0 or 1")
	}

	// wait be sleeping
	time.Sleep(10 * time.Millisecond)
	s.Equal(int32(1), atomic.LoadInt32(&s.dq.sleeping))

	// won't been wakeup
	n += 20
	s.T().Logf("Offer second element: priority %v", n)
	s.dq.Offer(n, n)
	s.Equal(2, s.dq.Size())
	s.Equal(int32(1), atomic.LoadInt32(&s.dq.sleeping))

	// been wakeup
	n -= 40
	s.T().Logf("Offer third element: priority %v", n)
	s.dq.Offer(n, n)
	s.Equal(3, s.dq.Size())
	// TODO(yangsonglin): difficult to test, we expect to see 0
	v = atomic.LoadInt32(&s.dq.sleeping)
	if v != 0 && v != 1 {
		s.FailNow("sleeping[%v] doesn't be 0 or 1")
	}

	// wait be sleeping
	time.Sleep(10 * time.Millisecond)
	s.Equal(int32(1), atomic.LoadInt32(&s.dq.sleeping))

	cancel()
	wg.Wait()
}

func (s *delayQueueTestSuite) TestPoll() {
	poll = func(ctx context.Context, q *delayQueue) bool {
		return false
	}
	defer func() {
		poll = pollImpl
	}()

	s.dq.Poll(context.Background())
	s.Equal(int32(0), atomic.LoadInt32(&s.dq.sleeping))
}

func (s *delayQueueTestSuite) TestPollImplNullItem() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	// cancel context
	n := defaultTimer.Now()
	r := pollImpl(ctx, s.dq)
	s.Equal(false, r)
	s.True(defaultTimer.Now()-n <= 2)
	s.Equal(int32(1), atomic.LoadInt32(&s.dq.sleeping))

	// been wakeup
	go func() {
		s.dq.wakeupC <- struct{}{}
	}()
	r = pollImpl(context.Background(), s.dq)
	s.Equal(true, r)
	s.Equal(int32(1), atomic.LoadInt32(&s.dq.sleeping))
}

func (s *delayQueueTestSuite) TestPollImplAfterDelayItem() {
	// test the item should been fired
	s.dq.Offer(1, 0)
	s.Equal(1, s.dq.Size())

	// fired
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		v := <-s.dq.Chan()
		s.Equal(1, v.(int))
	}()

	r := pollImpl(context.Background(), s.dq)
	s.Equal(true, r)
	s.Equal(int32(0), atomic.LoadInt32(&s.dq.sleeping))
	s.Equal(0, s.dq.Size())
	wg.Wait()

	// been wakeup
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	go func(wakeupC chan struct{}) {
		wakeupC <- struct{}{}
	}(s.dq.wakeupC)

	s.dq.Offer(1, 0)
	s.Equal(1, s.dq.Size())
	r = pollImpl(ctx, s.dq)
	s.Equal(false, r)
	s.Equal(int32(0), atomic.LoadInt32(&s.dq.sleeping))
	s.Equal(1, s.dq.Size())
}

func (s *delayQueueTestSuite) TestPollImplBeforeDelayItemFired() {
	// test the item shouldn't been fired
	n := defaultTimer.Now() + 1000
	s.dq.Offer(1, n)
	s.Equal(1, s.dq.Size())

	// wait been fired
	r := pollImpl(context.Background(), s.dq)
	s.Equal(true, r)
	s.Equal(int32(0), atomic.LoadInt32(&s.dq.sleeping))
	s.Equal(1, s.dq.Size())
}

func (s *delayQueueTestSuite) TestPollImplBeforeDelayItemWakeup() {
	// test the item shouldn't been fired
	n := defaultTimer.Now() + 1000

	// been wakeup
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	go func() {
		s.dq.wakeupC <- struct{}{}
	}()
	s.dq.Offer(1, n)
	s.Equal(1, s.dq.Size())
	r := pollImpl(ctx, s.dq)
	s.Equal(true, r)
	s.Equal(int32(1), atomic.LoadInt32(&s.dq.sleeping))
	s.Equal(1, s.dq.Size())
}

func (s *delayQueueTestSuite) TestPollImplBeforeDelayItemCancel() {
	// test the item shouldn't been fired
	n := defaultTimer.Now() + 1000

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	s.dq.Offer(1, n)
	s.Equal(1, s.dq.Size())
	r := pollImpl(ctx, s.dq)
	s.Equal(false, r)
	s.Equal(int32(1), atomic.LoadInt32(&s.dq.sleeping))
	s.Equal(1, s.dq.Size())
}

func (s *delayQueueTestSuite) TestSize() {
	ctx, cancel := context.WithCancel(context.Background())

	type testCase struct {
		value       int
		expireation int64
	}
	delay := int64(2000) // 2000ms, 2s
	testCases := []testCase{
		{
			value:       1,
			expireation: defaultTimer.Now() + int64(rand.Intn(int(delay))),
		},
		{
			value:       2,
			expireation: defaultTimer.Now() + int64(rand.Intn(int(delay))),
		},
		{
			value:       3,
			expireation: defaultTimer.Now() + int64(rand.Intn(int(delay))),
		},
		{
			value:       4,
			expireation: defaultTimer.Now() + int64(rand.Intn(int(delay))),
		},
	}
	sortTestCases := make([]testCase, len(testCases))
	copy(sortTestCases, testCases)
	sort.Slice(sortTestCases, func(i int, j int) bool {
		return sortTestCases[i].expireation < sortTestCases[j].expireation
	})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer func() {
			s.T().Logf("Poll done")
			wg.Done()
		}()

		s.dq.Poll(ctx)
	}()

	go func() {
		defer func() {
			s.T().Logf("Receive done")
			wg.Done()
		}()

		i := 0
		for v := range s.dq.Chan() {
			n := defaultTimer.Now()
			diff := n - sortTestCases[i].expireation
			if diff < 0 || diff > 10 {
				s.Failf("ExpireationPrecesionFailed", "Receive at[%v], estimate[%v]", n, sortTestCases[i].expireation)
			}

			s.T().Logf("Receive %v", v.(int))
			s.Equal(sortTestCases[i].value, v)
			i++
			if i >= len(sortTestCases) {
				cancel()
				return
			}
		}
	}()

	for i, tc := range testCases {
		s.dq.Offer(tc.value, tc.expireation)
		s.Equal(i+1, s.dq.Size())
	}

	wg.Wait()

	s.Equal(0, s.dq.Size())
}

func TestDelayQueueTestSuite(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	s := &delayQueueTestSuite{}
	suite.Run(t, s)
}
