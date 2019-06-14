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

type utilsTestSuite struct {
	suite.Suite
}

func (s *utilsTestSuite) TestTruncate() {
	type testCase struct {
		desp   string
		x      int64
		m      int64
		expect int64
	}
	testCases := []testCase{
		{
			desp:   "negative m",
			x:      1,
			m:      -1,
			expect: 1,
		},
		{
			desp:   "zero m",
			x:      100,
			m:      0,
			expect: 100,
		},
		{
			desp:   "divide exactly",
			x:      100,
			m:      2,
			expect: 100,
		},
		{
			desp:   "divide unexactly 01",
			x:      100,
			m:      3,
			expect: 99,
		},
		{
			desp:   "divide unexactly 02",
			x:      10,
			m:      3,
			expect: 9,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.desp, func() {
			actual := truncate(tc.x, tc.m)
			s.Equal(tc.expect, actual)
		})
	}
}

func TestUtilsTestSuite(t *testing.T) {
	s := &utilsTestSuite{}
	suite.Run(t, s)
}
