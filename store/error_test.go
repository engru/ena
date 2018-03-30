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
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type errorTestSuite struct {
	suite.Suite
}

func (s *errorTestSuite) TestNewError() {
	for k, v := range errorsMessage {
		e := NewError(k, v)
		s.Equal(k, e.ErrorCode)
		s.Equal(v, e.Message)
		s.Equal(v, e.Cause)
	}
}

func (s *errorTestSuite) TestNewErrorUnkownCode() {
	code := 0
	cause := "Unknown"

	e := NewError(code, cause)
	s.Equal(code, e.ErrorCode)
	s.Equal("", e.Message)
	s.Equal(cause, e.Cause)
}

func (s *errorTestSuite) TestJSONString() {
	e := NewError(EcodeNotDir, "TestJSONString")
	str := e.JSONString()

	str2, err := json.Marshal(e)
	s.NoError(err)
	s.Equal(str2, str)
}

func (s *errorTestSuite) TestJSONStringError() {
	marshal = func(interface{}) ([]byte, error) {
		return nil, errors.New("Error Marshal failed")
	}
	defer func() {
		marshal = json.Marshal
	}()

	e := NewError(EcodeNotDir, "TestJSONString")
	str := e.JSONString()

	err := &Error{
		ErrorCode: 1,
		Message:   "Error Marshal failed",
		Cause:     e.Error(),
	}

	str2, err2 := json.Marshal(err)
	s.NoError(err2)
	s.Equal(str2, str)
}

func TestErrorTestSuite(t *testing.T) {
	s := &errorTestSuite{}
	suite.Run(t, s)
}
