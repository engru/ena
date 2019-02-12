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

type dateBuilder struct{}

func (b *dateBuilder) Key() string {
	return string(FieldKeyDate)
}

func (b *dateBuilder) Build(layoutField *LayoutField) (FieldConverter, error) {
	return &dateConverter{
		timestampFormat: defaultTimestampFormat,
	}, nil
}

type levelBuilder struct{}

func (b *levelBuilder) Key() string {
	return string(FieldKeyLevel)
}

func (b *levelBuilder) Build(layoutField *LayoutField) (FieldConverter, error) {
	return &levelConverter{}, nil
}

type packageBuilder struct{}

func (b *packageBuilder) Key() string {
	return string(FieldKeyPackage)
}

func (b *packageBuilder) Build(layoutField *LayoutField) (FieldConverter, error) {
	return &packageConverter{}, nil
}

type fileBuilder struct{}

func (b *fileBuilder) Key() string {
	return string(FieldKeyFile)
}

func (b *fileBuilder) Build(layoutField *LayoutField) (FieldConverter, error) {
	return &fileConverter{}, nil
}

type methodBuilder struct{}

func (b *methodBuilder) Key() string {
	return string(FieldKeyMethod)
}

func (b *methodBuilder) Build(layoutField *LayoutField) (FieldConverter, error) {
	return &methodConverter{}, nil
}

type lineBuilder struct{}

func (b *lineBuilder) Key() string {
	return string(FieldKeyLine)
}

func (b *lineBuilder) Build(layoutField *LayoutField) (FieldConverter, error) {
	return &lineConverter{}, nil
}

type messageBuilder struct{}

func (b *messageBuilder) Key() string {
	return string(FieldKeyMessage)
}

func (b *messageBuilder) Build(layoutField *LayoutField) (FieldConverter, error) {
	return &messageConverter{}, nil
}
