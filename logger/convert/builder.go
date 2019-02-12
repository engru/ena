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
	"fmt"
)

var (
	fieldBuilders = map[string]FieldConverterBuilder{}
)

// Builder build layout strings to Converter
type Builder interface {
	// Build construct Converter from layout pattern string
	Build(layout string) (Converter, error)
}

// FieldConverterBuilder build FieldConverer
type FieldConverterBuilder interface {
	Build(layoutField *LayoutField) (FieldConverter, error)
	Key() string
}

// defBuilderImpl implement Builder interface
type defBuilderImpl struct {
}

// NewBuilder return Builder
func NewBuilder() Builder {
	return &defBuilderImpl{}
}

func (b *defBuilderImpl) Build(layout string) (Converter, error) {
	fields, err := ParseToLayoutField(layout)
	if err != nil {
		return nil, err
	}

	fieldConverters := make([]FieldConverter, 0, len(fields))
	for _, field := range fields {
		switch {
		case field.Type == LayoutFieldText:
			fieldConverters = append(fieldConverters, &textConverter{
				text: field.Value,
			})
		case field.Type == LayoutFieldConverter:
			if builder, ok := fieldBuilders[field.Value]; ok {
				var converter FieldConverter
				converter, err = builder.Build(field)
				if err != nil {
					return nil, err
				}
				fieldConverters = append(fieldConverters, converter)
			} else {
				return nil, fmt.Errorf("unsupported field key: %s", field.Value)
			}
		default:
			return nil, fmt.Errorf("unexpected field type: %v", field.Type)
		}
	}

	return NewConverter(fieldConverters), nil
}

// AddFieldBuilder add user field converter builder, if key conflict it will overwrite old field converter builder
func AddFieldBuilder(f FieldConverterBuilder) {
	fieldBuilders[f.Key()] = f
}

func init() {
	AddFieldBuilder(&dateBuilder{})
	AddFieldBuilder(&levelBuilder{})
	AddFieldBuilder(&packageBuilder{})
	AddFieldBuilder(&fileBuilder{})
	AddFieldBuilder(&methodBuilder{})
	AddFieldBuilder(&lineBuilder{})
	AddFieldBuilder(&messageBuilder{})
}
