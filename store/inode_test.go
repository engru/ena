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
	"testing"

	"github.com/stretchr/testify/suite"
)

type inodeTestSuite struct {
	suite.Suite
}

func (s *inodeTestSuite) TestIsHidden() {
	values := map[string]struct {
		node     *inode
		isHidden bool
	}{
		"1": {node: newFileInode(nil, "/", "", nil), isHidden: false},
		"2": {node: newFileInode(nil, "/xx", "", nil), isHidden: false},
		"3": {node: newFileInode(nil, "/.key", "", nil), isHidden: true},
		"4": {node: newFileInode(nil, "/key/xx", "", nil), isHidden: false},
		"5": {node: newFileInode(nil, "/key/xx/.ky", "", nil), isHidden: true},
		"6": {node: newDirInode(nil, "/", nil), isHidden: false},
		"7": {node: newDirInode(nil, "/.k", nil), isHidden: true},
		"8": {node: newDirInode(nil, "/.k/v", nil), isHidden: false},
		"9": {node: newDirInode(nil, "/.k/xx/.ee", nil), isHidden: true},
	}

	for _, v := range values {
		s.Equal(v.isHidden, v.node.IsHidden())
	}
}

func (s *inodeTestSuite) TestIsDir() {
	values := map[string]struct {
		node     *inode
		isHidden bool
	}{
		"1": {node: newFileInode(nil, "/", "", nil), isHidden: false},
		"2": {node: newFileInode(nil, "/xx", "", nil), isHidden: false},
		"3": {node: newFileInode(nil, "/.key", "", nil), isHidden: false},
		"4": {node: newFileInode(nil, "/key/xx", "", nil), isHidden: false},
		"5": {node: newFileInode(nil, "/key/xx/.ky", "", nil), isHidden: false},
		"6": {node: newDirInode(nil, "/", nil), isHidden: true},
		"7": {node: newDirInode(nil, "/.k", nil), isHidden: true},
		"8": {node: newDirInode(nil, "/.k/v", nil), isHidden: true},
		"9": {node: newDirInode(nil, "/.k/xx/.ee", nil), isHidden: true},
	}

	for _, v := range values {
		s.Equal(v.isHidden, v.node.IsDir())
	}
}

func TestInodeTestSuite(t *testing.T) {
	s := &inodeTestSuite{}
	suite.Run(t, s)
}
