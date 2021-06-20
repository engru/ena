// MIT License

// Copyright (c) 2018 soren yang

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

// Package helper contains common convert function
package helper

import (
	"encoding/json"

	"github.com/lsytj0413/ena/conversion"
)

// nolint
// Convert_map_string_string_to_string convert *map[string]string -> string
func Convert_map_string_string_to_string(in *map[string]string, out *string, _ conversion.Scope) error {
	if len(*in) == 0 {
		*out = ""
		return nil
	}

	data, err := json.Marshal(*in)
	if err != nil {
		return err
	}

	*out = string(data)
	return nil
}

func init() {
	conversion.DefaultConverter.RegisterConversionFunc((*map[string]string)(nil), (*string)(nil), func(in interface{}, out interface{}, scope conversion.Scope) error {
		return Convert_map_string_string_to_string(in.(*map[string]string), out.(*string), scope)
	})
}
