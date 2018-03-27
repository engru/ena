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
	"path"
	"strings"
)

const (
	root = "/"
)

func key(nodePath string) string {
	return path.Clean(path.Join(root, nodePath))
}

func isRoot(nodePath string) bool {
	return nodePath == root
}

func components(nodePath string) []string {
	return strings.Split(nodePath, root)
}
