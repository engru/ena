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

import (
	"fmt"
	"math"
	"sync/atomic"
)

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

	segments []Segment
}

// NewMap will construct a Map instance
func NewMap(concurrency uint32, pairRedistributor PairRedistributor) (Map, error) {
	if concurrency <= 0 {
		return nil, newInvalidParamError(fmt.Sprintf("concurrency should in range of [1, %d]", MaxConcurrency))
	}
	if concurrency > MaxConcurrency {
		return nil, newInvalidParamError(fmt.Sprintf("concurrency should in range of [1, %d]", MaxConcurrency))
	}

	m := &defMap{}
	m.concurrency = concurrency
	m.segments = make([]Segment, concurrency)
	for i := 0; i < int(concurrency); i++ {
		m.segments[i] = newSegment(int(DefaultBucketNumber), pairRedistributor)
	}

	return m, nil
}

func (m *defMap) Concurrency() uint32 {
	return m.concurrency
}

func (m *defMap) Put(key string, v interface{}) (bool, error) {
	p, err := newPair(key, v)
	if err != nil {
		return false, err
	}

	s := m.findSegment(p.Hash())
	ok, err := s.Put(p)
	if ok {
		atomic.AddUint64(&m.total, 1)
	}
	return ok, err
}

func (m *defMap) Get(key string) (interface{}, bool) {
	keyHash := hash(key)
	s := m.findSegment(keyHash)
	pair := s.GetWithHash(key, keyHash)
	if pair == nil {
		return nil, false
	}

	return pair.Value(), true
}

func (m *defMap) Delete(key string) bool {
	s := m.findSegment(hash(key))
	if s.Delete(key) {
		atomic.AddUint64(&m.total, ^uint64(0))
		return true
	}

	return false
}

func (m *defMap) Len() uint64 {
	return atomic.LoadUint64(&m.total)
}

func (m *defMap) findSegment(keyHash uint64) Segment {
	if m.concurrency == 1 {
		return m.segments[0]
	}

	var keyHash32 uint32
	if keyHash > math.MaxUint32 {
		keyHash32 = uint32(keyHash >> 32)
	} else {
		keyHash32 = uint32(keyHash)
	}

	return m.segments[uint32(keyHash32>>16)%(m.concurrency-1)]
}
