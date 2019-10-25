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

package art

import (
	"bytes"
	"fmt"

	"github.com/lsytj0413/ena/logger/convert/parser"
)

const (
	// TokenizerName defines of art
	TokenizerName = "art"
)

type artTokenizer struct {
	state   tokenizerState
	pos     int
	len     int
	pattern string

	tokenStream parser.Tokens
	buf         bytes.Buffer
}

func (p *artTokenizer) Name() string {
	return TokenizerName
}

func (p *artTokenizer) Parse(pattern string) (parser.Tokens, error) {
	p.initialize(pattern)

	var step func(byte) error
	var err error
	for p.pos < p.len {
		c := pattern[p.pos]
		p.pos++

		switch p.state {
		case literalState:
			step = p.onLiteralState
		}

		err = step(c)
		if err != nil {
			return nil, err
		}
	}

	return p.tokenStream, nil
}

func (p *artTokenizer) onLiteralState(c byte) error {
	switch {
	case c == escapeChar:
		return p.escape()
	case c == percentChar:
	default:
		p.buf.WriteByte(c)
	}

	return nil
}

func (p *artTokenizer) escape() error {
	if p.pos >= p.len {
		return fmt.Errorf("unexpected char '\\' at end of pattern")
	}

	next := p.pattern[p.pos]
	p.pos++
	switch next {
	case percentChar, '\\':
		p.buf.WriteByte(next)
	case 'r':
		p.buf.WriteByte('\r')
	case 't':
		p.buf.WriteByte('\t')
	case 'n':
		p.buf.WriteByte('\n')
	default:
		return fmt.Errorf("unexpected char %v after '\\' at index %v", next, p.pos-1)
	}

	return nil
}

func (p *artTokenizer) initialize(pattern string) {
	p.state = literalState
	p.pos = 0
	p.len = len(pattern)
	p.pattern = pattern
	p.tokenStream = make(parser.Tokens, 0)
	p.buf = bytes.Buffer{}
}

// NewTokenizer construct art parser tokenizer instance
func NewTokenizer() (parser.Tokenizer, error) {
	p := &artTokenizer{}
	p.initialize("")
	return p, nil
}
