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
	"strings"
	"sync"
)

// FileSystemStore defines a filesystem like kv store
type FileSystemStore interface {
	// Get nodePath node infomation
	Get(nodePath string, recursive bool, sorted bool) (*Result, error)
	// Set value to nodePath
	Set(nodePath string, dir bool, value string) (*Result, error)
	// Update value to nodePath
	Update(nodePath string, newValue string) (*Result, error)
	// Create nodePath with value
	Create(nodePath string, dir bool, value string) (*Result, error)
	// Delete nodePath
	Delete(nodePath string, dir bool, recursive bool) (*Result, error)
}

// defaultFileSystemStore implemented FileSystemStore interface
type defaultFileSystemStore struct {
	Root *inode
	lock sync.RWMutex
}

// NewFileSystemStore creates a FileSystemStore with root directories
func NewFileSystemStore() FileSystemStore {
	s := new(defaultFileSystemStore)
	s.Root, _ = newDirInode(s, "/", nil)

	return s
}

// Get returns Node
// If recursive is true, it will return all the content under the node path
func (s *defaultFileSystemStore) Get(
	nodePath string,
	recursive bool,
	sorted bool) (*Result, error) {
	var err error

	s.lock.RLock()
	defer s.lock.RUnlock()

	defer func() {
		if err == nil {
			fmt.Printf("Get %s success", nodePath)
			return
		}

		fmt.Printf("Get %s failed, %v", nodePath, err)
	}()

	n, err := s.get(nodePath)
	if err != nil {
		return nil, err
	}

	r := newResult(Get, nodePath)
	r.CurrNode.loadFromInode(n, recursive, sorted)
	return r, nil
}

func (s *defaultFileSystemStore) Set(
	nodePath string,
	dir bool,
	value string) (*Result, error) {
	var err error

	s.lock.RLock()
	defer s.lock.RUnlock()

	defer func() {
		if err == nil {
			fmt.Printf("Set %s success", nodePath)
			return
		}

		fmt.Printf("Set %s failed, %v", nodePath, err)
	}()

	// First, get prevNode Value
	_, err = s.get(nodePath)
	if err != nil {
		return nil, err
	}

	e, err := s.create(nodePath, dir, value)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *defaultFileSystemStore) Update(
	nodePath string,
	newValue string) (*Result, error) {
	return nil, errors.New("Not Implement")
}

func (s *defaultFileSystemStore) Create(
	nodePath string,
	dir bool,
	value string) (*Result, error) {
	return nil, errors.New("Not Implement")
}

func (s *defaultFileSystemStore) create(nodePath string, dir bool, value string) (*Result, error) {
	nodePath = path.Clean(path.Join("/", nodePath))

	dirName, nodeName := path.Split(nodePath)

	d, err := s.walk(dirName, s.checkDir)

	if err != nil {
		return nil, err
	}

	r := newResult(Create, nodePath)
	eNode := r.CurrNode

	n, _ := d.GetChild(nodeName)
	if n != nil {
		return nil, errors.New("inode exsits")
	}

	if dir {
		valueCopy := value
		eNode.Value = &valueCopy

		n, _ = newFileInode(s, nodePath, value, d)
	} else {
		eNode.Dir = true
		n, _ = newDirInode(s, nodePath, d)
	}

	d.Add(n)

	return r, nil
}

// checkDir will check dirName under parent
// If is directory, return inode
// If does not exsits, create a new directory and return inode
// If is file, return an error
func (s *defaultFileSystemStore) checkDir(parent *inode, dirName string) (*inode, error) {
	node, ok := parent.Children[dirName]
	if ok {
		if node.IsDir() {
			return node, nil
		}

		return nil, errors.New("Not a dir")
	}

	n, _ := newDirInode(s, path.Join(parent.Path, dirName), parent)
	parent.Children[dirName] = n
	return n, nil
}

func (s *defaultFileSystemStore) Delete(
	nodePath string,
	dir bool,
	recursive bool) (*Result, error) {
	return nil, errors.New("Not Implement")
}

func (s *defaultFileSystemStore) get(nodePath string) (*inode, error) {
	nodePath = path.Clean(path.Join("/", nodePath))

	walkFunc := func(parent *inode, name string) (*inode, error) {
		if !parent.IsDir() {
			return nil, errors.New("Not Dir")
		}

		child, ok := parent.Children[name]
		if ok {
			return child, nil
		}

		return nil, errors.New("Key Not Found")
	}

	f, err := s.walk(nodePath, walkFunc)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *defaultFileSystemStore) walk(nodePath string, walkFunc func(prev *inode, component string) (*inode, error)) (*inode, error) {
	components := strings.Split(nodePath, "/")

	curr := s.Root
	var err error

	for i := 1; i < len(components); i++ {
		if len(components[i]) == 0 {
			return curr, nil
		}

		curr, err = walkFunc(curr, components[i])
		if err != nil {
			return nil, err
		}
	}

	return curr, nil
}
