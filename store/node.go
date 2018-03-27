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
	"fmt"
	"sort"
)

// Node is the external representation of the inode with additional fields
type Node struct {
	Key   string
	Value *string
	Dir   bool
	Nodes NodeArray
}

func (n Node) String() string {
	var value string
	if n.Value != nil {
		value = *n.Value
	}

	return fmt.Sprintf("Node(Key=%s, Value=%s, Dir=%v)", n.Key, value, n.Dir)
}

func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}

	nn := &Node{
		Key: n.Key,
		Dir: n.Dir,
	}

	if n.Value != nil {
		nn.Value = &(*n.Value)
	}

	if n.Nodes != nil {
		nn.Nodes = nn.Nodes.Clone()
	}

	return nn
}

func (n *Node) loadFromInode(in *inode, recursive bool, sorted bool) {
	if in.IsDir() {
		n.Dir = true

		children, _ := in.List()
		n.Nodes = make(NodeArray, len(children))

		i := 0

		for _, child := range children {
			if child.IsHidden() {
				continue
			}

			// n.Nodes[i] = child.Repr(recursive, sorted)
			n.Nodes[i] = inodeToNode(child, recursive, sorted)
			i++
		}

		// slice down length
		n.Nodes = n.Nodes[:i]

		if sorted {
			sort.Sort(n.Nodes)
		}
	} else {
		n.Value = &in.Value
	}
}

// NodeArray is list of Node
type NodeArray []*Node

func (na NodeArray) Len() int {
	return len(na)
}

func (na NodeArray) Less(i int, j int) bool {
	return na[i].Key < na[j].Key
}

func (na NodeArray) Swap(i int, j int) {
	na[i], na[j] = na[j], na[i]
}

func (na NodeArray) Clone() NodeArray {
	nodes := make(NodeArray, na.Len())
	for i, n := range na {
		nodes[i] = n.Clone()
	}

	return nodes
}
