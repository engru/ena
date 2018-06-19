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

package algo

import (
	"errors"
	"hash/fnv"
)

// HashCarp will choice an endpoints for key use CARP algorithm
func HashCarp(key string, endpoints []string) (idx int, err error) {
	if 0 == len(endpoints) {
		return -1, errors.New("endpoints length must be greater than 0")
	}

	if 1 == len(endpoints) {
		return 0, nil
	}

	hashedArr := make([]uint64, len(endpoints))
	for i, v := range endpoints {
		hashedArr[i] = hashUint64(key + v)
	}

	min := hashedArr[0]
	for i, v := range hashedArr {
		if v < min {
			idx = i
			min = v
		}
	}
	return
}

func hashUint64(value string) uint64 {
	a := fnv.New64()
	a.Write([]byte(value))
	return a.Sum64()
}
