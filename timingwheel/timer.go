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
	"time"
)

// timerTask represent single task. When expires, the given
// task will been executed.
type timerTask struct {
	// d is the duration of timertask
	d time.Duration
	// expiration of the task
	expiration int64

	// the id of timertask, unique
	id uint64

	// task handler
	f Handler

	// sign of whether the timertask has stopped,
	// 1: stopped
	// 0: non stopped
	stopped uint32

	// the bucket pointer that holds the TimerTask list
	b *bucket
	w *timingWheel

	e *list.Element
}

func (t *timerTask) bucket() *bucket {
	return t.b
}

func (t *timerTask) setBucket(b *bucket) bool {
	t.b = b
	return true
}

// Stop the timer task from fire, return true if the timer is stopped success,
// or false if the timer has already expired or been stopped.
func (t *timerTask) Stop() bool {
	if atomic.LoadUint32(&t.stopped) == 1 {
		return true
	}

	return t.w.StopFunc(t) || atomic.LoadUint32(&t.stopped) == 1
}
