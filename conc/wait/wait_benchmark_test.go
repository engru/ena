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

package wait

import (
	"testing"
	"time"
)

func Benchmark_Wait(b *testing.B) {
	w := New()

	for i := 0; i < b.N; i++ {
		w.Register(uint64(i))
		w.Trigger(uint64(i), nil)
	}
}

func Benchmark_Wait_Parallel(b *testing.B) {
	w := New()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			i := uint64(time.Now().UnixNano())
			w.Register(i)
			w.Trigger(i, nil)
		}
	})
}
