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

package common

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type sortTestSuite struct {
	suite.Suite
}

func (p *sortTestSuite) TestSortIntOk() {
	testcase := []struct {
		arg1   []int
		target []int
	}{
		{
			arg1:   []int{3, 2, 1},
			target: []int{1, 2, 3},
		},
		{
			arg1:   []int{1, 7, 3, 6, 20, -1},
			target: []int{-1, 1, 3, 6, 7, 20},
		},
	}
	for _, tc := range testcase {
		Sort(tc.arg1, func(data interface{}, i int, j int) bool {
			d := data.([]int)
			return d[i] < d[j]
		}, func(data interface{}) int {
			d := data.([]int)
			return len(d)
		}, func(data interface{}, i int, j int) {
			d := data.([]int)
			d[i], d[j] = d[j], d[i]
		})

		p.Equal(tc.target, tc.arg1)
	}
}

func (p *sortTestSuite) TestSortStringOk() {
	testcase := []struct {
		arg1   []string
		target []string
	}{
		{
			arg1:   []string{"3", "2", "1"},
			target: []string{"1", "2", "3"},
		},
		{
			arg1:   []string{"b", "a", "ac"},
			target: []string{"a", "ac", "b"},
		},
	}
	for _, tc := range testcase {
		Sort(tc.arg1, func(data interface{}, i int, j int) bool {
			d := data.([]string)
			return d[i] < d[j]
		}, func(data interface{}) int {
			d := data.([]string)
			return len(d)
		}, func(data interface{}, i int, j int) {
			d := data.([]string)
			d[i], d[j] = d[j], d[i]
		})

		p.Equal(tc.target, tc.arg1)
	}
}

func TestSortTestSuite(t *testing.T) {
	p := &sortTestSuite{}
	suite.Run(t, p)
}
