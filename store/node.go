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
	"errors"
	"path"
	"sort"
)

// inode is basic element in the store system
type inode struct {
	Path  string
	Value string

	Parent   *inode
	Children map[string]*inode // for directory

	// A reference to the store this inode is attached to
	store *defaultFileSystemStore
}

func newFileInode(store *defaultFileSystemStore, nodePath string, value string, parent *inode) (*inode, error) {
	return &inode{
		Path:   nodePath,
		Value:  value,
		Parent: parent,
		store:  store,
	}, nil
}

func newDirInode(store *defaultFileSystemStore, nodePath string, parent *inode) (*inode, error) {
	return &inode{
		Path:     nodePath,
		Parent:   parent,
		Children: make(map[string]*inode),
		store:    store,
	}, nil
}

func (n *inode) IsHidden() bool {
	_, name := path.Split(n.Path)
	return name[0] == '.'
}

func (n *inode) IsDir() bool {
	return n.Children != nil
}

func (n *inode) Read() {

}

func (n *inode) Write() {

}

func (n *inode) List() ([]*inode, error) {
	if !n.IsDir() {
		return nil, errors.New("Not Dir")
	}

	nodes := make([]*inode, len(n.Children))

	i := 0
	for _, node := range n.Children {
		nodes[i] = node
		i++
	}

	return nodes, nil
}

// GetChild returns the child inode under the directory inode
// If current inode is file, returns error
// If child not exists, returns error
func (n *inode) GetChild(name string) (*inode, error) {
	if !n.IsDir() {
		return nil, errors.New("Not Dir")
	}

	child, ok := n.Children[name]
	if ok {
		return child, nil
	}

	return nil, errors.New("File not Exists")
}

// Add function adds a inode to the directory inode
// If current inode is not directory, returns error
// If same name already exists under the directory, returns error
func (n *inode) Add(child *inode) error {
	if !n.IsDir() {
		return errors.New("Not Dir")
	}

	_, name := path.Split(child.Path)
	if _, ok := n.Children[name]; ok {
		return errors.New("already exists")
	}

	n.Children[name] = child
	return nil
}

func (n *inode) Remove(dir bool, recursive bool) {

}

// Repr will translate inode to Node expression
func (n *inode) Repr(recursive bool, sorted bool) *Node {
	if n.IsDir() {
		node := &Node{
			Key: n.Path,
			Dir: true,
		}

		if !recursive {
			return node
		}

		children, _ := n.List()
		node.Nodes = make(NodeArray, len(children))

		i := 0

		for _, child := range children {
			if child.IsHidden() {
				continue
			}

			node.Nodes[i] = child.Repr(recursive, sorted)
			i++
		}

		node.Nodes = node.Nodes[:i]
		if sorted {
			sort.Sort(node.Nodes)
		}

		return node
	}

	value := n.Value
	node := &Node{
		Key:   n.Path,
		Value: &value,
	}
	return node
}

func (n *inode) Clone() {

}

// Node is the external representation of the inode with additional fields
type Node struct {
	Key   string
	Value *string
	Dir   bool
	Nodes NodeArray
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

			n.Nodes[i] = child.Repr(recursive, sorted)
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
