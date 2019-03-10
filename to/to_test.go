// MIT License

// Copyright (c) 2018 soren yang

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

package to

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type toTestSuite struct {
	suite.Suite
}

func (s *toTestSuite) TestString() {
	type testcase struct {
		desc   string
		expect string
	}
	testcases := []testcase{
		{
			desc:   "empty string",
			expect: "",
		},
		{
			desc:   "normal string",
			expect: "normal",
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			s.Equal(tc.expect, String(&v))
		})
	}
}

func (s *toTestSuite) TestStringNil() {
	expect := ""
	actual := String(nil)
	s.Equal(expect, actual)
}

func (s *toTestSuite) TestStringPtr() {
	type testcase struct {
		desc   string
		expect string
	}
	testcases := []testcase{
		{
			desc:   "empty string",
			expect: "",
		},
		{
			desc:   "normal string",
			expect: "normal",
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			s.Equal(tc.expect, *StringPtr(v))
		})
	}
}

func (s *toTestSuite) TestStringSlice() {
	type testcase struct {
		desc   string
		expect []string
	}
	testcases := []testcase{
		{
			desc:   "empty slice",
			expect: []string{},
		},
		{
			desc:   "normal string",
			expect: []string{"normal"},
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := StringSlice(&v)
			s.True(reflect.DeepEqual(tc.expect, actual))
		})
	}
}

func (s *toTestSuite) TestStringSliceHandlesNil() {
	var expect []string
	s.Equal(expect, StringSlice(nil))
}

func (s *toTestSuite) TestStringSlicePtr() {
	type testcase struct {
		desc   string
		expect []string
	}
	testcases := []testcase{
		{
			desc:   "empty slice",
			expect: []string{},
		},
		{
			desc:   "normal string",
			expect: []string{"normal"},
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := StringSlicePtr(v)
			s.True(reflect.DeepEqual(tc.expect, *actual))
		})
	}
}

func (s *toTestSuite) TestBool() {
	type testcase struct {
		desc   string
		expect bool
	}
	testcases := []testcase{
		{
			desc:   "true",
			expect: true,
		},
		{
			desc:   "false",
			expect: false,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Bool(&v)
			s.Equal(tc.expect, actual)
		})
	}
}

func (s *toTestSuite) TestBoolHandlesNil() {
	s.Equal(false, Bool(nil))
}

func (s *toTestSuite) TestBoolPtr() {
	type testcase struct {
		desc   string
		expect bool
	}
	testcases := []testcase{
		{
			desc:   "true",
			expect: true,
		},
		{
			desc:   "false",
			expect: false,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := BoolPtr(v)
			s.Equal(tc.expect, *actual)
		})
	}
}

func (s *toTestSuite) TestInt() {
	type testcase struct {
		desc   string
		expect int
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0,
		},
		{
			desc:   "positive",
			expect: 100,
		},
		{
			desc:   "negative",
			expect: -100,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Int(&v)
			s.Equal(tc.expect, actual)
		})
	}
}

func (s *toTestSuite) TestIntHandlesNil() {
	s.Equal(int(0), Int(nil))
}

func (s *toTestSuite) TestIntPtr() {
	type testcase struct {
		desc   string
		expect int
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0,
		},
		{
			desc:   "positive",
			expect: 100,
		},
		{
			desc:   "negative",
			expect: -100,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := IntPtr(v)
			s.Equal(tc.expect, *actual)
		})
	}
}

func (s *toTestSuite) TestInt32() {
	type testcase struct {
		desc   string
		expect int32
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0,
		},
		{
			desc:   "positive",
			expect: 100,
		},
		{
			desc:   "negative",
			expect: -100,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Int32(&v)
			s.Equal(tc.expect, actual)
		})
	}
}

func (s *toTestSuite) TestInt32HandlesNil() {
	s.Equal(int32(0), Int32(nil))
}

func (s *toTestSuite) TestInt32Ptr() {
	type testcase struct {
		desc   string
		expect int32
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0,
		},
		{
			desc:   "positive",
			expect: 100,
		},
		{
			desc:   "negative",
			expect: -100,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Int32Ptr(v)
			s.Equal(tc.expect, *actual)
		})
	}
}

func (s *toTestSuite) TestInt64() {
	type testcase struct {
		desc   string
		expect int64
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0,
		},
		{
			desc:   "positive",
			expect: 100,
		},
		{
			desc:   "negative",
			expect: -100,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Int64(&v)
			s.Equal(tc.expect, actual)
		})
	}
}

func (s *toTestSuite) TestInt64HandlesNil() {
	s.Equal(int64(0), Int64(nil))
}

func (s *toTestSuite) TestInt64Ptr() {
	type testcase struct {
		desc   string
		expect int64
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0,
		},
		{
			desc:   "positive",
			expect: 100,
		},
		{
			desc:   "negative",
			expect: -100,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Int64Ptr(v)
			s.Equal(tc.expect, *actual)
		})
	}
}

func (s *toTestSuite) TestFloat32() {
	type testcase struct {
		desc   string
		expect float32
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0.0,
		},
		{
			desc:   "positive",
			expect: 100.0,
		},
		{
			desc:   "negative",
			expect: -100.0,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Float32(&v)
			s.Equal(tc.expect, actual)
		})
	}
}

func (s *toTestSuite) TestFloat32HandlesNil() {
	s.Equal(float32(0), Float32(nil))
}

func (s *toTestSuite) TestFloat32Ptr() {
	type testcase struct {
		desc   string
		expect float32
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0.0,
		},
		{
			desc:   "positive",
			expect: 100.0,
		},
		{
			desc:   "negative",
			expect: -100.0,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Float32Ptr(v)
			s.Equal(tc.expect, *actual)
		})
	}
}

func (s *toTestSuite) TestFloat64() {
	type testcase struct {
		desc   string
		expect float64
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0.0,
		},
		{
			desc:   "positive",
			expect: 100.0,
		},
		{
			desc:   "negative",
			expect: -100.0,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Float64(&v)
			s.Equal(tc.expect, actual)
		})
	}
}

func (s *toTestSuite) TestFloat64HandlesNil() {
	s.Equal(float64(0), Float64(nil))
}

func (s *toTestSuite) TestFloat64Ptr() {
	type testcase struct {
		desc   string
		expect float64
	}
	testcases := []testcase{
		{
			desc:   "zero",
			expect: 0.0,
		},
		{
			desc:   "positive",
			expect: 100.0,
		},
		{
			desc:   "negative",
			expect: -100.0,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			v := tc.expect
			actual := Float64Ptr(v)
			s.Equal(tc.expect, *actual)
		})
	}
}

func TestToTetSuite(t *testing.T) {
	s := &toTestSuite{}
	suite.Run(t, s)
}
