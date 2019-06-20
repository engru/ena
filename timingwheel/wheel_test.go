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

package timingwheel

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lsytj0413/ena/delayqueue"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockDelayQueue struct {
	mock.Mock
}

func (m *mockDelayQueue) Chan() <-chan interface{} {
	return nil
}

func (m *mockDelayQueue) Size() int {
	return 0
}

func (m *mockDelayQueue) Offer(element interface{}, expiration int64) {
	args := m.Called(element, expiration)
	_ = args
}

func (m *mockDelayQueue) Poll(ctx context.Context) {

}

var _ delayqueue.DelayQueue = &mockDelayQueue{}

type wheelTestSuite struct {
	suite.Suite

	w  *wheel
	dq *mockDelayQueue
}

func (s *wheelTestSuite) SetupTest() {
	s.w = newWheel(3, 20, 4)
	s.dq = &mockDelayQueue{}
	defaultExecutor = blockExecutor

	s.Equal(int64(3), s.w.tick)
	s.Equal(int64(20), s.w.wheelSize)
	s.Equal(int64(60), s.w.interval)
	s.Equal(int64(3), s.w.currentTime)
	s.Equal(int(20), len(s.w.buckets))

	for _, b := range s.w.buckets {
		s.NotNil(b)
	}
}

func (s *wheelTestSuite) TearDownTest() {
	defaultExecutor = taskExecutor
}

func (s *wheelTestSuite) TestNewWheel() {
	tickMs, wheelSize, startMs := int64(1), int64(20), int64(0)
	s.w = newWheel(tickMs, wheelSize, startMs)
	s.Equal(tickMs, s.w.tick)
	s.Equal(wheelSize, s.w.wheelSize)
	s.Equal(tickMs*wheelSize, s.w.interval)
	s.Equal(int64(0), s.w.currentTime)
	s.Equal(int(wheelSize), len(s.w.buckets))

	for _, b := range s.w.buckets {
		s.NotNil(b)
	}
}

func (s *wheelTestSuite) TestAddFalse() {
	for i := int64(0); i < s.w.currentTime+s.w.tick; i++ {
		desp := fmt.Sprintf("exp %v", i)
		s.Run(desp, func() {
			s.Equal(false, s.w.add(&timerTask{
				expiration: int64(i),
			}, nil))
		})
	}
}

func (s *wheelTestSuite) TestAddOkCurrentWheel() {
	type testCase struct {
		desp             string
		t                *timerTask
		bucketIndex      int
		bucketExpiration int64
		isOfferCalled    bool
		currentTime      int64
	}
	testCases := []testCase{
		{
			desp: "task expiration 7, 2th bucket",
			t: &timerTask{
				expiration: 7,
			},
			bucketIndex:      2,
			bucketExpiration: 6,
			isOfferCalled:    true,
		},
		{
			desp: "task expiration 17, 5th bucket",
			t: &timerTask{
				expiration: 17,
			},
			bucketIndex:      5,
			bucketExpiration: 15,
			isOfferCalled:    true,
		},
		{
			desp: "task expiration 62, 0th bucket",
			t: &timerTask{
				expiration: 62,
			},
			bucketIndex:      0,
			bucketExpiration: 60,
			isOfferCalled:    true,
		},
		{
			desp: "task expiration 8, 2th bucket",
			t: &timerTask{
				expiration: 8,
			},
			bucketIndex:      2,
			bucketExpiration: 6,
			isOfferCalled:    false,
		},
		{
			desp:        "task expiration 16, 5th bucket",
			currentTime: 6,
			t: &timerTask{
				expiration: 16,
			},
			bucketIndex:      5,
			bucketExpiration: 15,
			isOfferCalled:    false,
		},
	}
	for _, tc := range testCases {
		s.dq = &mockDelayQueue{}
		s.Run(tc.desp, func() {
			if tc.currentTime > 0 {
				s.w.advanceClock(tc.currentTime)
			}

			s.dq.On("Offer", mock.MatchedBy(func(b *bucket) bool {
				return b == tc.t.b
			}), tc.bucketExpiration)
			s.Equal(true, s.w.add(tc.t, s.dq))

			s.Equal(tc.t.b, s.w.buckets[tc.bucketIndex])
			s.Equal(tc.bucketExpiration, s.w.buckets[tc.bucketIndex].Expiration())

			if tc.isOfferCalled {
				s.dq.AssertCalled(s.T(), "Offer", tc.t.b, tc.bucketExpiration)
			} else {
				s.dq.AssertNotCalled(s.T(), "Offer", tc.t.b, tc.bucketExpiration)
			}
		})
	}
}

func (s *wheelTestSuite) TestAddOkOverflowWheel() {
	s.Nil(s.w.overflowWheel)
	// the overflowwheel will cover the uplayerwheel range. at this EX:
	// 1. uplayer: tick(3), current(3), wheelsize(20), interval(60)
	// 2. uplayer: range(3-63)
	// 3. overflow layer:
	//   a. tick(60): uplayer interval
	//   b. current(0): truncate(uplayer current, tick)
	//   c. wheelsize(20): uplayer wheelsize
	//   d. interval(1200): tick*wheelsize
	// 4. overflow layer: range(0-1200)
	s.dq.On("Offer", mock.MatchedBy(func(b *bucket) bool {
		return true
	}), int64(60))
	s.Equal(true, s.w.add(&timerTask{
		expiration: 66,
	}, s.dq))
	s.NotNil(s.w.overflowWheel)
	s.NotNil(s.w.overflowWheel.buckets[1])
	s.Equal(int64(60), s.w.overflowWheel.buckets[1].Expiration())
}

func (s *wheelTestSuite) TestAddOrRun() {
	s.Run("run immediate after task", func() {
		v := 0
		t := &timerTask{
			expiration: 0,
			f: func(time.Time) {
				v = 1
			},
			t: taskAfter,
		}
		s.w.addOrRun(t, s.dq)
		s.Equal(1, v)
		s.Nil(t.b)
		s.dq.AssertNotCalled(s.T(), "Offer",
			mock.MatchedBy(func(b interface{}) bool {
				return true
			}),
			mock.MatchedBy(func(int64) bool {
				return true
			}),
		)
	})

	s.Run("run immediate nonstoped tick task", func() {
		v := 0
		t := &timerTask{
			expiration: 0,
			d:          time.Duration(7),
			stopped:    0,
			f: func(time.Time) {
				v = 1
			},
			t: taskTick,
		}
		s.dq = &mockDelayQueue{}
		s.dq.On("Offer",
			mock.MatchedBy(func(b interface{}) bool {
				return true
			}),
			mock.MatchedBy(func(int64) bool {
				return true
			}),
		)
		s.w.addOrRun(t, s.dq)
		s.Equal(1, v)
		s.NotNil(t.b)
		s.dq.AssertCalled(s.T(), "Offer",
			mock.MatchedBy(func(b interface{}) bool {
				return true
			}),
			mock.MatchedBy(func(int64) bool {
				return true
			}),
		)
	})

	s.Run("run immediate stopped tick task", func() {
		v := 0
		t := &timerTask{
			expiration: 0,
			d:          time.Duration(7),
			stopped:    1,
			f: func(time.Time) {
				v = 1
			},
			t: taskTick,
		}
		s.dq = &mockDelayQueue{}
		s.dq.On("Offer",
			mock.MatchedBy(func(b interface{}) bool {
				return true
			}),
			mock.MatchedBy(func(int64) bool {
				return true
			}),
		)
		s.w.addOrRun(t, s.dq)
		s.Equal(1, v)
		s.Nil(t.b)
		s.dq.AssertNotCalled(s.T(), "Offer",
			mock.MatchedBy(func(b interface{}) bool {
				return true
			}),
			mock.MatchedBy(func(int64) bool {
				return true
			}),
		)
	})
}

func (s *wheelTestSuite) TestAdvanceClock() {
	s.w.overflowWheel = newWheel(0, 0, 0)
	s.Equal(int64(0), s.w.overflowWheel.currentTime)

	exp := int64(100)
	s.w.advanceClock(exp)
	exp = int64(99)
	s.Equal(exp, s.w.currentTime)
	s.Equal(exp, s.w.overflowWheel.currentTime)
}

func TestWheelTestSuite(t *testing.T) {
	s := &wheelTestSuite{}
	suite.Run(t, s)
}
