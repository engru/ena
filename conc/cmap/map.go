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

// Package cmap defines a Goroutine Safe Map Container
package cmap

// Map is interface define for Concurrency Safe Map
type Map interface {
	Concurrency() uint32
	Put(key string, element interface{}) (bool, error)
	Get(key string) (interface{}, bool)
	Delete(key string) bool
	Len() uint64
}

type defMap struct {
	concurrency uint32
	total       uint64
}
