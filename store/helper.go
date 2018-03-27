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

import "sort"

// translate inode to Node expression
// If recursive is true, will translate children inode
// If sorted is true, will sort the children
func inodeToNode(n *inode, recursive bool, sorted bool) *Node {
	if !n.IsDir() {
		// not a directory, translate the current inode and return
		value := n.Value
		return &Node{
			Key:   n.Path,
			Value: &value,
		}
	}

	// directory
	node := &Node{
		Key: n.Path,
		Dir: true,
	}

	if !recursive {
		// just return the current node
		return node
	}

	children, _ := n.List()
	node.Nodes = make(NodeArray, len(children))

	// translate children inode
	i := 0
	for _, child := range children {
		if child.IsHidden() {
			continue
		}

		node.Nodes[i] = inodeToNode(child, recursive, sorted)
		i++
	}

	// slice, because mybe hidden node exists
	node.Nodes = node.Nodes[:i]
	if sorted {
		sort.Sort(node.Nodes)
	}
	return node
}
