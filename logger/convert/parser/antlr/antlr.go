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

package antlr

import (
	"fmt"

	"github.com/lsytj0413/ena/logger/convert/parser"
)

const (
	// TokenizerName defines of antlr
	TokenizerName = "antlr"
)

type antlrTokenizer struct {
}

func (p *antlrTokenizer) Name() string {
	return TokenizerName
}

func (p *antlrTokenizer) Parse(pattern string) (parser.Tokens, error) {
	return nil, fmt.Errorf("NotImplement")
}

// NewTokenizer construct antlr tokenizer parser instance
func NewTokenizer() (parser.Tokenizer, error) {
	return &antlrTokenizer{}, nil
}
