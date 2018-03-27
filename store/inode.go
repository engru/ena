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
	"fmt"
	"path"
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

func (n inode) String() string {
	return fmt.Sprintf("inode(Path=%s, Value=%s)", n.Path, n.Value)
}

func newFileInode(store *defaultFileSystemStore, nodePath string, value string, parent *inode) *inode {
	return &inode{
		Path:   nodePath,
		Value:  value,
		Parent: parent,
		store:  store,
	}
}

func newDirInode(store *defaultFileSystemStore, nodePath string, parent *inode) *inode {
	return &inode{
		Path:     nodePath,
		Parent:   parent,
		Children: make(map[string]*inode),
		store:    store,
	}
}

func (n *inode) IsHidden() bool {
	_, name := path.Split(n.Path)
	return name[0] == '.'
}

func (n *inode) IsDir() bool {
	return n.Children != nil
}

// Read function gets the value of the node
// If node is a directory, fail
func (n *inode) Read() (string, error) {
	if n.IsDir() {
		return "", errors.New("Not File")
	}

	return n.Value, nil
}

// Write function set the value of the node to the given value
// If node is a directory, fail
func (n *inode) Write(value string) error {
	if n.IsDir() {
		return errors.New("Not file")
	}

	n.Value = value
	return nil
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

// Remove function remove the node
func (n *inode) Remove(dir bool, recursive bool) error {
	if !n.IsDir() {
		_, name := path.Split(n.Path)
		if n.Parent != nil && n.Parent.Children[name] == n {
			delete(n.Parent.Children, name)
		}

		return nil
	}

	if !dir {
		return errors.New("Not File")
	}

	if len(n.Children) != 0 && !recursive {
		return errors.New("Dir Not Empty")
	}

	for _, child := range n.Children {
		child.Remove(true, true)
	}

	// Delete self
	_, name := path.Split(n.Path)
	if n.Parent != nil && n.Parent.Children[name] == n {
		delete(n.Parent.Children, name)
	}
	return nil
}

// Repr will translate inode to Node expression
// func (n *inode) Repr(recursive bool, sorted bool) *Node {
// 	if n.IsDir() {
// 		node := &Node{
// 			Key: n.Path,
// 			Dir: true,
// 		}

// 		if !recursive {
// 			return node
// 		}

// 		children, _ := n.List()
// 		node.Nodes = make(NodeArray, len(children))

// 		i := 0

// 		for _, child := range children {
// 			if child.IsHidden() {
// 				continue
// 			}

// 			node.Nodes[i] = child.Repr(recursive, sorted)
// 			i++
// 		}

// 		node.Nodes = node.Nodes[:i]
// 		if sorted {
// 			sort.Sort(node.Nodes)
// 		}

// 		return node
// 	}

// 	value := n.Value
// 	node := &Node{
// 		Key:   n.Path,
// 		Value: &value,
// 	}
// 	return node
// }

// Clone function clone the node recursively and return the new node
// If the node is a directory, it will clone all the content under this directory
// If the node is a file, it will clone the file
func (n *inode) Clone() *inode {
	if !n.IsDir() {
		return newFileInode(n.store, n.Path, n.Value, n.Parent)
	}

	d := newDirInode(n.store, n.Path, n.Parent)
	for key, child := range n.Children {
		d.Children[key] = child.Clone()
	}

	return d
}
