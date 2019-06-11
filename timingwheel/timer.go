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

import (
	"container/list"
	"sync/atomic"
	"unsafe"
)

// TimerTask represent single task. When expires, the given
// task will been executed.
type TimerTask struct {
	// expiration of the task
	expiration int64

	// task handler
	f Handler

	// the bucket pointer that holds the TimerTask list
	b unsafe.Pointer

	e *list.Element
}

func (t *TimerTask) bucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.b))
}

func (t *TimerTask) setBucket(b *bucket) bool {
	old := atomic.LoadPointer(&t.b)
	return atomic.CompareAndSwapPointer(&t.b, old, unsafe.Pointer(b))
}

// Stop the timer task from fire, return true if the timer is stopped success,
// or false if the timer has already expired or been stopped.
func (t *TimerTask) Stop() bool {
	stopped := false
	for b := t.bucket(); b != nil; b = t.bucket() {
		stopped = b.Remove(t)
	}
	return stopped
}
