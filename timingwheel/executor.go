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

import "time"

// executor is an redefine for task execute
type executor = func(Handler, time.Time)

// taskExecutor is the default implement to run Handler, will execute the Handler in other goroutine
func taskExecutor(f Handler, ct time.Time) {
	go f(ct)
}

var (
	// defaultExecutor is the executor for task, it will execute the task in it's own goroutine
	defaultExecutor = taskExecutor
)

// blockExecutor will run task in current goruntine, for test usage
func blockExecutor(f Handler, ct time.Time) {
	f(ct)
}
