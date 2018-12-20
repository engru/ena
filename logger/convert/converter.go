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

package convert

import (
	"strings"
	"time"
)

// Converter convert fields to string
type Converter interface {
	Convert(entry *Entry) string
}

// Entry for convert input
type Entry struct {
	Time    time.Time
	Level   string
	Package string
	File    string
	Method  string
	Line    string
	Message string
	Data    map[string]string
}

// defConverterImpl implememnt Convert interface
type defConverterImpl struct {
	fields []FieldConverter
}

func (c *defConverterImpl) Convert(entry *Entry) string {
	buffer := strings.Builder{}
	for _, f := range c.fields {
		buffer.WriteString(f.Convert(entry))
	}
	return buffer.String()
}

// NewConverter construct Converter
func NewConverter(fields []FieldConverter) Converter {
	return &defConverterImpl{
		fields: fields,
	}
}
