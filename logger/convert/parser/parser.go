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

package parser

// TokenType represent type of Token
type TokenType = string

// Token for string group
type Token struct {
	// Type is the identify of tokentype
	Type TokenType

	// Value is the text of token
	Value string
}

const (
	// TokenTypePercent for '%'
	TokenTypePercent = '%'

	// TokenTypeText for text
	TokenTypeText = "text"

	// TokenTypeOpenBrace for option start
	TokenTypeOpenBrace = "{"

	// TokenTypeCloseBrace for option end
	TokenTypeCloseBrace = "}"

	// TokenTypeComma for split options
	TokenTypeComma = ","
)

// Tokens is array of Token
type Tokens = []Token

// Tokenizer for lexical parser,
type Tokenizer interface {
	Name() string
	Parse(pattern string) (Tokens, error)
}

// Compiler defines for compile pattern to fields
type Compiler interface {
	// Compile will parse the pattern to fields
	Compile(pattern string) (Fields, error)
}

// Field for layout field type and value
type Field struct {
	// Type is the type of field
	Type FieldType

	// Modifier represent modifier option for field, only seted when type is converter
	Modifier string

	// Value of field text
	Value string

	// Option represent converter options for field, only seted when type is converter
	Option []string
}

// FieldType for Value type
type FieldType = string

const (
	// FieldText for text layout field
	FieldText FieldType = "text"
	// FieldConverter for converter layout field
	FieldConverter FieldType = "converter"
)

// Fields is array of Field
type Fields = []Field
