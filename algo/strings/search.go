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

package strings

// Find will return the index where sub in str
// -1 will return when sub not in str
func Find(str string, sub string) int {
	return _normalFind(str, sub)
}

func _normalFind(str string, sub string) int {
	for i := 0; i < len(str)-len(sub)+1; i++ {
		j := 0
		for j = 0; j < len(sub); j++ {
			if str[i+j] != sub[j] {
				break
			}
		}
		if j == len(sub) {
			return i
		}
	}

	return -1
}
