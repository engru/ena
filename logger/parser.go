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
)

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func parse(pattern string) ([]*token, error) {
	tokens := make([]*token, 0)

	n, mode := bytes.Buffer{}, 0
	for i := 0; i < len(pattern); i++ {
		switch mode {
		case 0: // 初始状态
			if pattern[i] == '%' {
				mode = 1
				if n.Len() != 0 {
					tokens = append(tokens, &token{t: "t", v: n.String()})
					n.Reset()
				}
			} else {
				n.WriteByte(pattern[i])
			}
		case 1: // 上一个为 %, 进入 layout 匹配
			if pattern[i] == '%' {
				// 直接结束
				tokens = append(tokens, &token{t: "t", v: "%"})
				mode = 0
			} else if isAlpha(pattern[i]) {
				n.WriteByte(pattern[i])
				mode = 2
			} else {
				return nil, fmt.Errorf("unexpected char after %% at %d", i)
			}
		case 2:
			if isAlpha(pattern[i]) {
				n.WriteByte(pattern[i])
			} else {
				mode = 0
				tokens = append(tokens, &token{t: "c", v: n.String()})
				n.Reset()
				i--
			}
		}
	}

	if mode == 1 {
		return nil, errors.New("unexpected % at end of pattern")
	}

	if n.Len() != 0 {
		if mode == 0 {
			tokens = append(tokens, &token{t: "t", v: n.String()})
		} else if mode == 2 {
			tokens = append(tokens, &token{t: "c", v: n.String()})
		}
	}

	return tokens, nil
}
