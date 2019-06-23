// MIT License

// Copyright (c) 2019 soren yang

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package convert

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	parse func(string) ([]*LayoutField, error)
)

func defaultParser(layout string) ([]*LayoutField, error) {
	fields := make([]*LayoutField, 0)

	// TODO: merge same type layout field
	n, mode := bytes.Buffer{}, 0
	modesp := map[int]func(int) (int, int, error){
		0: func(i int) (next int, m int, err error) {
			next = i + 1
			if layout[i] != '%' {
				n.WriteByte(layout[i])
				return
			}

			m = 1
			if n.Len() != 0 {
				fields = append(fields, &LayoutField{
					Type:  LayoutFieldText,
					Value: n.String(),
				})
				n.Reset()
			}
			return
		},
		1: func(i int) (next int, m int, err error) {
			next = i + 1
			switch {
			case layout[i] == '%':
				m = 0
				fields = append(fields, &LayoutField{
					Type:  LayoutFieldText,
					Value: "%",
				})
				return
			case IsAlpha(layout[i]):
				n.WriteByte(layout[i])
				m = 2
				return
			default:
				return -1, -1, fmt.Errorf("unexpected char[%v] after %% at %d", layout[i], i)
			}
		},
		2: func(i int) (next int, m int, err error) {
			if IsAlpha(layout[i]) {
				n.WriteByte(layout[i])
				next = i + 1
				return
			}

			m = 0
			fields = append(fields, &LayoutField{
				Type:  LayoutFieldConverter,
				Value: n.String(),
			})
			n.Reset()
			next = i
			return
		},
	}

	var err error
	for i := 0; i < len(layout); {
		i, mode, err = modesp[mode](i)
		if err != nil {
			return nil, err
		}
	}

	if mode == 1 {
		return nil, errors.New("unexpected % at end of pattern")
	}

	if n.Len() != 0 {
		if mode == 0 {
			fields = append(fields, &LayoutField{
				Type:  LayoutFieldText,
				Value: n.String(),
			})
		} else if mode == 2 {
			fields = append(fields, &LayoutField{
				Type:  LayoutFieldConverter,
				Value: n.String(),
			})
		}
	}

	return fields, nil
}

func init() {
	parse = defaultParser
}
