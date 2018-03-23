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

func (n *inode) List() {

}

func (n *inode) GetChild(name string) {

}

func (n *inode) Add(child *inode) {

}

func (n *inode) Remove(dir bool, recursive bool) {

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
