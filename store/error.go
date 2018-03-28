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

import "encoding/json"

// Error is store package error message define
type Error struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	Cause     string `json:"cause,omitempty"`
}

const (
	// EcodeUnknown is unknown error info
	EcodeUnknown = 10009999
	// EcodeNotFile errors for operate on dir but file is required
	EcodeNotFile = 10000001
	// EcodeNotDir errors for operate on file but dir is required
	EcodeNotDir = 10000002
	// EcodeFileNotExists errors for operate on file but doesn't exists
	EcodeFileNotExists = 10000003
	// EcodeFileExists errors for Add file but already exists
	EcodeFileExists = 10000004
	// EcodeDirNotEmpty errors for Remove directory but directory has child etc
	EcodeDirNotEmpty = 10000005
)

var errorsMessage = map[int]string{
	EcodeUnknown:       "Unknown Error",
	EcodeNotFile:       "Target is Not File",
	EcodeNotDir:        "Target is Not Dir",
	EcodeFileNotExists: "Target file is not exists",
	EcodeFileExists:    "Target file is exists",
}

// NewError construct a Error struct and return it
func NewError(errorCode int, cause string) *Error {
	return &Error{
		ErrorCode: errorCode,
		Message:   errorsMessage[errorCode],
		Cause:     cause,
	}
}

// Error is for the error interface
func (e Error) Error() string {
	return e.Message + " (" + e.Cause + ")"
}

// JSONString returns the JSON format message
func (e Error) JSONString() string {
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(b)
}
