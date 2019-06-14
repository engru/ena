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
)

// buchet is list of timer
type bucket struct {
	// the container list of timer
	timers *list.List

	// expiration is the expire of bucket
	expiration int64
}

func newBucket() *bucket {
	return &bucket{
		timers:     list.New(),
		expiration: -1,
	}
}

func (b *bucket) Expiration() int64 {
	return atomic.LoadInt64(&b.expiration)
}

func (b *bucket) SetExpiration(v int64) bool {
	return atomic.SwapInt64(&b.expiration, v) != v
}

func (b *bucket) Add(t *timerTask) {
	e := b.timers.PushBack(t)
	t.setBucket(b)
	t.e = e
}

func (b *bucket) remove(t *timerTask) bool {
	if t.bucket() != b {
		return false
	}

	b.timers.Remove(t.e)
	t.setBucket(nil)
	t.e = nil
	return true
}

// Flush will called by bucket expire exceed, and in the reinsert function,
// the TimerTask will reinsert into the timingwheel.
func (b *bucket) Flush(reinsert func(*timerTask)) {
	e := b.timers.Front()
	for e != nil {
		n := e.Next()
		t := e.Value.(*timerTask)
		b.remove(t)
		reinsert(t)
		e = n
	}

	b.SetExpiration(-1)
}
