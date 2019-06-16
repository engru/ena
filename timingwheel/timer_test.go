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
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockStopWheel struct {
	mock.Mock
	stopped bool
}

func (m *mockStopWheel) StopFunc(t *timerTask) (bool, error) {
	args := m.Called(t)
	if m.stopped {
		t.stopped = 1
	}
	return args.Get(0).(bool), args.Error(1)
}

type timerTaskTestSuite struct {
	suite.Suite

	t *timerTask
	w *mockStopWheel
}

func (s *timerTaskTestSuite) SetupTest() {
	s.t = &timerTask{}
	s.w = &mockStopWheel{}
	s.t.w = s.w
}

func (s *timerTaskTestSuite) TestStopStopped() {
	s.t.stopped = 1
	v, err := s.t.Stop()
	s.NoError(err)
	s.True(v)
}

func (s *timerTaskTestSuite) TestStopFailed() {
	err := fmt.Errorf("MockFailed")
	s.w.On("StopFunc", s.t).Return(false, err)
	v, verr := s.t.Stop()
	s.False(v)
	s.Error(verr, err.Error())
}

func (s *timerTaskTestSuite) TestStopOk() {
	s.w.On("StopFunc", s.t).Return(true, nil)
	v, verr := s.t.Stop()
	s.True(v)
	s.NoError(verr)
}

func (s *timerTaskTestSuite) TestStopThenStopped() {
	s.w.On("StopFunc", s.t).Return(false, nil)
	s.w.stopped = true
	v, verr := s.t.Stop()
	s.True(v)
	s.NoError(verr)
	s.Equal(s.t.stopped, uint32(1))
}

func TestTimerTaskTestSuite(t *testing.T) {
	s := &timerTaskTestSuite{}
	suite.Run(t, s)
}
