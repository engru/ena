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

package builder

import (
	"fmt"
	"strings"

	"github.com/lsytj0413/ena/logger/convert/parser"
	"github.com/lsytj0413/ena/logger/convert/parser/antlr"
	"github.com/lsytj0413/ena/logger/convert/parser/art"
)

// AvailableTokenizers is the name list of supported tokenizers
var AvailableTokenizers = []string{
	art.TokenizerName,
	antlr.TokenizerName,
}

// AvailableTokenizersDescription is string readable description for tokenizer list
var AvailableTokenizersDescription string

// Builder for build tokenizer from name
type Builder interface {

	// Build will return tokenizer instance from name
	Build(tokenizerName string) (parser.Tokenizer, error)
}

// NewBuilder return Builder instance
func NewBuilder() Builder {
	return &builder{}
}

// builder implement the Builder
type builder struct{}

// Build implement the builder.Build
func (b *builder) Build(tokenizerName string) (parser.Tokenizer, error) {
	switch tokenizerName {
	case art.TokenizerName:
		return createArtTokenizer()
	case antlr.TokenizerName:
		return createAntlrTokenizer()
	}

	return nil, fmt.Errorf("Invalid Tokenizer Name[%v], Available: %v", tokenizerName, AvailableTokenizersDescription)
}

var createAntlrTokenizer = func() (parser.Tokenizer, error) {
	return antlr.NewTokenizer()
}

var createArtTokenizer = func() (parser.Tokenizer, error) {
	return art.NewTokenizer()
}

func init() {
	AvailableTokenizersDescription = "[" + strings.Join(AvailableTokenizers, ", ") + "]"
}
