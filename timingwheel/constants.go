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

package timingwheel

import "fmt"

var (
	// ErrInvalidTickValue is representation error of invalid tick value
	ErrInvalidTickValue = fmt.Errorf("tick must be greater than or equal to 1ms")

	// ErrInvalidWheelSize is representation error of invalid wheel size value
	ErrInvalidWheelSize = fmt.Errorf("wheel size must greater than zero")
)

// eventType is the representation of event, such as AddNew, RePost
type eventType = string

// eventAddNew is the identify when timertask is add from AfterFunc
var eventAddNew eventType = "AddNew"

// eventDelete is the identify
var eventDelete eventType = "Delete"
