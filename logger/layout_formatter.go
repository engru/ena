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

package logger

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

// LayoutFormatter ...
// %d dateformat
// %level
// %M method name
// %L line
// %msg
// %p package name
type LayoutFormatter struct {
	pattern string
	c       []Converter
}

// NewLayoutFormatter ...
func NewLayoutFormatter(pattern string) (*LayoutFormatter, error) {
	f := &LayoutFormatter{
		pattern: pattern,
	}

	tokens, err := parse(pattern)
	if err != nil {
		return nil, err
	}

	for _, t := range tokens {
		switch t.t {
		case "t":
			f.c = append(f.c, &textConverter{v: t.v})
		case "c":
			c, exists := defaultConverterMap[t.v]
			if !exists {
				return nil, fmt.Errorf("unsupported token value: %s", t.v)
			}
			f.c = append(f.c, c)
		default:
			return nil, errors.New("unsupported token type")
		}
	}

	return f, nil
}

// Format ...
func (l *LayoutFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := bytes.Buffer{}
	for _, c := range l.c {
		b.WriteString(c.Convert(entry))
	}

	return b.Bytes(), nil
}
