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

package store

import (
	"container/list"
	"sync"
)

// watcherHub contains all subscribed watchers
type watcherHub struct {
	count uint64 // current number of watchers

	mutex         sync.Mutex
	watchers      map[string]*list.List
	ResultHistory *ResultHistory
}

func newWatchHub(capacity int) *watcherHub {
	return &watcherHub{
		watchers:      make(map[string]*list.List),
		ResultHistory: newResultHistory(capacity),
	}
}

func (h *watcherHub) watch(key string, recursive bool, stream bool, index uint64, storeIndex uint64) (Watcher, error) {
    return nil, nil
}
