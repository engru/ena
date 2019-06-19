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
	"testing"

	"github.com/stretchr/testify/suite"
)

type bucketTestSuite struct {
	suite.Suite

	b *bucket
}

func (s *bucketTestSuite) SetupTest() {
	s.b = newBucket()
}

func (s *bucketTestSuite) TestSetExpiration() {
	s.Equal(int64(-1), s.b.expiration)

	s.Equal(false, s.b.SetExpiration(int64(-1)))
	s.Equal(int64(-1), s.b.Expiration())

	s.Equal(true, s.b.SetExpiration(int64(1)))
	s.Equal(int64(1), s.b.Expiration())
}

func (s *bucketTestSuite) TestAdd() {
	t := &timerTask{}
	s.b.Add(t)
	s.Equal(s.b, t.b)
	s.NotNil(t.e)
}

func (s *bucketTestSuite) TestRemoveFailed() {
	t := &timerTask{}
	s.b.Add(t)

	t.b = nil
	s.Equal(false, s.b.remove(t))
	s.Nil(t.b)
	s.NotNil(t.e)
}

func (s *bucketTestSuite) TestRemoveOk() {
	t := &timerTask{}
	s.b.Add(t)

	s.Equal(true, s.b.remove(t))
	s.Nil(t.b)
	s.Nil(t.e)
}

func (s *bucketTestSuite) TestFlush() {
	timers := []*timerTask{
		{},
		{},
	}
	for _, t := range timers {
		s.b.Add(t)
	}

	gotTimers := make([]*timerTask, 0, len(timers))
	reinsert := func(t *timerTask) {
		gotTimers = append(gotTimers, t)
	}

	s.b.Flush(reinsert)

	for _, t := range timers {
		s.Nil(t.b)
		s.Nil(t.e)
	}

	for i := range timers {
		s.Equal(timers[i], gotTimers[i])
	}
}

func TestBucketTestSuite(t *testing.T) {
	s := &bucketTestSuite{}
	suite.Run(t, s)
}
