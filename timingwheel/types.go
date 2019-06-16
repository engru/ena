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

// Handler for function execution
type Handler func()

// TimingWheel is an interface for implementation.
type TimingWheel interface {
	// Start starts the current timing wheel
	Start()

	// Stop stops the current timing wheel. If there is any timer's task being running, the stop
	// will not wait for complete.
	// TODO(lsytj0413): should pass ctx from Start?
	Stop()

	// AfterFunc will call the Handler in its own goroutine after the duration elapse.
	// It return an Timer that can use to cancel the Handler.
	AfterFunc(d time.Duration, f Handler) (TimerTask, error)

	// TickFunc will call the Handler in its own goroutine after the duration elapse tick.
	// It reutrn an Timer that can use to cancel the Handler.
	TickFunc(d time.Duration, f Handler) (TimerTask, error)
}

// TimerTask is an interface for task implementation.
type TimerTask interface {
	Stop() (bool, error)
}
