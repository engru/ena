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
	"time"

	"github.com/lsytj0413/ena/delayqueue"
)

// newWheel is the interval implement of creates time wheel instance.
// it always used when add timertask into timing wheel for creates overflow timingwheel.
func newWheel(tickMs int64, wheelSize int64, startMs int64) *wheel {
	buckets := make([]*bucket, wheelSize)
	for i := range buckets {
		buckets[i] = newBucket()
	}

	return &wheel{
		tick:        tickMs,
		wheelSize:   wheelSize,
		interval:    tickMs * wheelSize,
		currentTime: truncate(startMs, tickMs),
		buckets:     buckets,
	}
}

type wheel struct {
	// tick is the interval of every bucket representation in milliseconds,
	// it the min expire unit in the timmingWheel, so it always be time.Milliseconds in the first
	// layer timingWheel.
	tick int64

	// wheelSize is the bucket count of layer
	wheelSize int64

	// interval is the count of tick*wheelSize, it's the interval of all this layer can representation
	interval int64

	// currentTime is the current time in milliseconds
	currentTime int64

	// buckets is the array of bucket, the len is wheelSize
	buckets []*bucket

	// overflowWheel is the high-layer timing wheel
	overflowWheel *wheel
}

func (w *wheel) addOrRun(t *timerTask, dq delayqueue.DelayQueue) {
	if !w.add(t, dq) {
		now := time.Now()

		// the timertask already expired, wo we run execute the timer's task in its own goroutine.
		defaultExecutor(t.f, now)

		if t.t == taskTick && t.stopped == 0 {
			// the timertask is tick func, and haven't been stopped, reinsert it
			t.expiration = timeToMs(now.Add(t.d))
			w.addOrRun(t, dq)
		}
	}
}

func (w *wheel) add(t *timerTask, dq delayqueue.DelayQueue) bool {
	switch {
	case t.expiration < w.currentTime+w.tick:
		// if the timertask is in the first bucket, we treat it as expired.
		return false
	case t.expiration < w.currentTime+w.interval:
		// the timertask is in current layer wheel

		// vid is the multiple of expireation and tick,
		// EX: the tick is 2ms, and expiration is 9ms, so the vid will be 4
		vid := t.expiration / w.tick

		// b is the bucket witch the timertask should put in
		// EX: the tick is 2ms, and expiration is 9ms, and wheelSize is 5,
		// it should put in the 5th bucket
		b := w.buckets[vid%w.wheelSize]
		b.Add(t)

		if b.SetExpiration(vid * w.tick) {
			// the bucket expiration has been changed, we enqueue it into the delayqueue
			// EX: the wheel has been advanced, and the bucket is reused after flush
			// 1. flush
			//   a. after the flush, the expiration will been set to -1, so it will always been updated
			// 2. advanced
			//   a. currentTime is 90, tick is 10, and wheelsize is 3, interval is 30
			//   b. timer task with expiration(101) is add, put in the bucket 1, bucket expiration is 100, and offered
			//   b. the currentTime is advanced to 100
			//   c. timer task with expiration(101) is add, put in the bucket 1, bucket expiration is 100, no offered
			//   e. timer task with expiration(121) is add, put in the bucket 0, bucket expiration is 120, and offered
			// the currentTime doesn't affect witch bucket put in and the bucket expiration, the same range expiration will
			// always been put in the same bucket.
			// but before advance the 121 will addto overflowwheel, and addto currentwheel after advanced. the two bucket has same
			// expiration(120).
			dq.Offer(b, b.Expiration())
		}
		return true
	default:
		if w.overflowWheel == nil {
			w.overflowWheel = newWheel(w.interval, w.wheelSize, w.currentTime)
		}
		return w.overflowWheel.add(t, dq)
	}
}

func (w *wheel) advanceClock(expiration int64) {
	if expiration >= w.currentTime+w.tick {
		w.currentTime = truncate(expiration, w.tick)

		if w.overflowWheel != nil {
			w.overflowWheel.advanceClock(w.currentTime)
		}
	}
}
