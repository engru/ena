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
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/suite"

	"github.com/lsytj0413/ena/conc"
)

type timingWheelTestSuite struct {
	suite.Suite

	tw *timingWheel
}

func (s *timingWheelTestSuite) SetupTest() {
	tw, err := NewTimingWheel(WithTickDuration(time.Millisecond), WithSize(20))
	s.NoError(err)
	s.tw = tw.(*timingWheel)
	defaultExecutor = blockExecutor
}

func (s *timingWheelTestSuite) TearDownTest() {
	defaultExecutor = taskExecutor
}

func (s *timingWheelTestSuite) TestNewTimingWheelInvalidTick() {
	values := []time.Duration{
		time.Nanosecond,
		20 * time.Nanosecond,
		200 * time.Nanosecond,
		time.Microsecond,
		20 * time.Microsecond,
		200 * time.Microsecond,
		999 * time.Microsecond,
	}
	for _, v := range values {
		tw, err := NewTimingWheel(WithTickDuration(v), WithSize(20))
		s.Error(err, ErrInvalidTickValue.Error())
		s.Nil(tw)
	}
}

func (s *timingWheelTestSuite) TestNewTimingWheelInvalidWheelSize() {
	values := []int64{
		-100,
		-50,
		-20,
		-1,
		0,
	}
	for _, v := range values {
		tw, err := NewTimingWheel(WithTickDuration(time.Millisecond), WithSize(v))
		s.Error(err, ErrInvalidWheelSize.Error())
		s.Nil(tw)
	}
}

func (s *timingWheelTestSuite) TestNewTimingWheelOk() {
	values := []time.Duration{
		time.Millisecond,
		20 * time.Millisecond,
		200 * time.Millisecond,
		time.Second,
		20 * time.Second,
		200 * time.Second,
		time.Minute,
		time.Hour,
	}
	for _, v := range values {
		tw, err := NewTimingWheel(WithTickDuration(v), WithSize(20))
		s.NoError(err)
		s.NotNil(tw)
	}
}

func (s *timingWheelTestSuite) TestAfterFunc() {
	type testCase struct {
		description string
		d           time.Duration
	}
	testCases := []testCase{
		{
			description: "2 ms",
			d:           2 * time.Millisecond,
		},
		{
			description: "10 ms",
			d:           10 * time.Millisecond,
		},
		{
			description: "300 ms",
			d:           300 * time.Millisecond,
		},
		{
			description: "1 s",
			d:           time.Second,
		},
		{
			description: "1.5 s",
			d:           time.Second + 500*time.Millisecond,
		},
		{
			description: "3 s",
			d:           3 * time.Second,
		},
	}

	var wg conc.WaitGroupWrapper
	now := timeToMs(time.Now())
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan struct{})

	wg.Wrap(func() {
		defer cancel()

		c := 0
		for {
			<-ch
			c++
			if c >= len(testCases) {
				return
			}
		}
	})

	wg.Wrap(func() {
		for _, tc := range testCases {
			s.tw.AfterFunc(tc.d, func(tc testCase) func(time.Time) {
				return func(ct time.Time) {
					defer func() {
						ch <- struct{}{}
					}()
					n := timeToMs(ct)
					expect := now + int64(tc.d/time.Millisecond)
					s.T().Logf("receive: %s", tc.description)
					if n < expect || n > expect+10 {
						s.T().Fatalf("receive %s: expect[%v], got[%v]", tc.description, expect, n)
					}
				}
			}(tc))
		}
	})

	s.tw.Start()
	wg.Wrap(func() {
		defer s.tw.Stop()
		<-ctx.Done()
	})

	wg.Wait()
}

func (s *timingWheelTestSuite) TestTickFunc() {
	type testCase struct {
		description string
		d           time.Duration
		last        unsafe.Pointer
		skip        int32
		t           TimerTask
	}
	testCases := []*testCase{
		{
			description: "2 ms",
			d:           2 * time.Millisecond,
		},
		{
			description: "10 ms",
			d:           10 * time.Millisecond,
		},
		{
			description: "300 ms",
			d:           300 * time.Millisecond,
		},
		{
			description: "1 s",
			d:           time.Second,
		},
		{
			description: "1.5 s",
			d:           time.Second + 500*time.Millisecond,
		},
		{
			description: "3 s",
			d:           3 * time.Second,
		},
	}

	var wg conc.WaitGroupWrapper
	timeout := time.Second * 4
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.tw.Start()
	wg.Wrap(func() {
		for _, tc := range testCases {
			now := time.Now()
			atomic.StorePointer(&tc.last, unsafe.Pointer(&now))
			t, err := s.tw.TickFunc(tc.d, func(tc *testCase) func(time.Time) {
				return func(ct time.Time) {
					now, lastptr := ct, (atomic.LoadPointer(&tc.last))
					last := *((*time.Time)(lastptr))
					expect := last.Add(tc.d)

					if now.Before(last) {
						// the next handler has been called, skip this
						atomic.AddInt32(&tc.skip, 1)
						return
					}

					atomic.CompareAndSwapPointer(&tc.last, lastptr, unsafe.Pointer(&ct))
					if expect.After(now.Add(2*time.Millisecond)) || now.After(expect.Add(10*time.Millisecond)) {
						s.T().Fatalf("receive %s: expect[%v], got[%v], last[%v]", tc.description, expect, now, last)
					}
				}
			}(tc))
			s.NoError(err)
			tc.t = t
		}
	})

	<-ctx.Done()
	wg.Wait()

	for _, tc := range testCases {
		tc.t.Stop()
	}
	s.tw.Stop()

	for _, tc := range testCases {
		v := atomic.LoadInt32(&tc.skip)
		if v > 0 {
			s.T().Logf("%v: skip times[%v]", tc.description, v)
		}
	}
}

func (s *timingWheelTestSuite) TestTickFuncFailed() {
	type tickDuration struct {
		desp string
		d    time.Duration
	}
	type testCase struct {
		desp string
		tick time.Duration
		ds   []tickDuration
	}

	testCases := []testCase{
		{
			desp: "tick(10ms)",
			tick: 10 * time.Millisecond,
			ds: []tickDuration{
				{
					desp: "d(1ms)",
					d:    1 * time.Millisecond,
				},
				{
					desp: "d(2ms)",
					d:    2 * time.Millisecond,
				},
				{
					desp: "d(7ms)",
					d:    7 * time.Millisecond,
				},
				{
					desp: "d(9ms)",
					d:    9 * time.Millisecond,
				},
			},
		},
		{
			desp: "tick(10s)",
			tick: 10 * time.Second,
			ds: []tickDuration{
				{
					desp: "d(1ms)",
					d:    1 * time.Millisecond,
				},
				{
					desp: "d(1.7s)",
					d:    1700 * time.Millisecond,
				},
				{
					desp: "d(5ms)",
					d:    5 * time.Second,
				},
				{
					desp: "d(9.999s)",
					d:    9999 * time.Millisecond,
				},
			},
		},
	}

	for _, tc := range testCases {
		tw, err := NewTimingWheel(WithTickDuration(tc.tick), WithSize(20))
		s.NoError(err)
		for _, tcd := range tc.ds {
			s.Run(tc.desp+"-"+tcd.desp, func() {
				tt, err := tw.TickFunc(tcd.d, func(time.Time) {})
				s.Error(err, ErrInvalidTickFuncDurationValue.Error())
				s.Nil(tt)
			})
		}
	}
}

func (s *timingWheelTestSuite) TestTickFuncOk() {
	type tickDuration struct {
		desp string
		d    time.Duration
	}
	type testCase struct {
		desp string
		tick time.Duration
		ds   []tickDuration
	}

	testCases := []testCase{
		{
			desp: "tick(10ms)",
			tick: 10 * time.Millisecond,
			ds: []tickDuration{
				{
					desp: "d(10ms)",
					d:    10 * time.Millisecond,
				},
				{
					desp: "d(20ms)",
					d:    20 * time.Millisecond,
				},
				{
					desp: "d(21ms)",
					d:    21 * time.Millisecond,
				},
				{
					desp: "d(11ms)",
					d:    11 * time.Millisecond,
				},
			},
		},
		{
			desp: "tick(10s)",
			tick: 10 * time.Second,
			ds: []tickDuration{
				{
					desp: "d(10s)",
					d:    10000 * time.Millisecond,
				},
				{
					desp: "d(10.001s)",
					d:    10001 * time.Millisecond,
				},
				{
					desp: "d(12s)",
					d:    12 * time.Second,
				},
				{
					desp: "d(10000s)",
					d:    10000 * time.Millisecond,
				},
			},
		},
	}

	for _, tc := range testCases {
		tw, err := NewTimingWheel(WithTickDuration(tc.tick), WithSize(20))
		s.NoError(err)
		tw.Start()
		for _, tcd := range tc.ds {
			s.Run(tc.desp+"-"+tcd.desp, func() {
				tt, err := tw.TickFunc(tcd.d, func(time.Time) {})
				s.NoError(err)
				s.NotNil(tt)
			})
		}
		tw.Stop()
	}
}

func TestTimingWheelTestSuite(t *testing.T) {
	s := &timingWheelTestSuite{}
	suite.Run(t, s)
}
