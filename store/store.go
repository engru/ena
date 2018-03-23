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

// FileSystemStore defines a filesystem like kv store
type FileSystemStore interface {
	// Get nodePath node infomation
	Get(nodePath string, recursive bool, sorted bool)
	// Set value to nodePath
	Set(nodePath string, dir bool, value string)
	// Update value to nodePath
	Update(nodePath string, newValue string)
	// Create nodePath with value
	Create(nodePath string, dir bool, value string)
	// Delete nodePath
	Delete(nodePath string, dir bool, recursive bool)
}

// defaultFileSystemStore implemented FileSystemStore interface
type defaultFileSystemStore struct {
}
