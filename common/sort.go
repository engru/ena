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
	"sort"
)

// Lessor is a function witch compare data's i and j element
type Lessor = func(data interface{}, i int, j int) bool

// Swaper is a function witch swap data's i and j element
type Swaper = func(data interface{}, i, j int)

// Lener is a function witch returns data's length
type Lener = func(data interface{}) int

type sortHelper struct {
	data   interface{}
	lessor Lessor
	swaper Swaper
	lener  Lener
}

func (s sortHelper) Less(i int, j int) bool {
	return s.lessor(s.data, i, j)
}

func (s sortHelper) Swap(i int, j int) {
	s.swaper(s.data, i, j)
}

func (s sortHelper) Len() int {
	return s.lener(s.data)
}

// Sort is a helper function for sort.Sort
func Sort(data interface{}, lessor Lessor, lener Lener, swaper Swaper) {
	sorter := &sortHelper{
		data:   data,
		lessor: lessor,
		swaper: swaper,
		lener:  lener,
	}
	sort.Sort(sorter)
}
