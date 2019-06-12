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

// Package timingwheel implementation of Hierarchical Timing Wheels.
package timingwheel

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/lsytj0413/ena/conc"
	"github.com/lsytj0413/ena/delayqueue"
)

// Handler for function execution
type Handler func()

// TimingWheel is an interface for implementation.
type TimingWheel interface {
	// Start starts the current timing wheel
	Start()

	// Stop stops the current timing wheel. If there is any timer's task being running, the stop
	// will not wait for complete.
	Stop()

	// AfterFunc will call the Handler in its own goroutine after the duration elapse.
	// It return an Timer that can use to cancel the Handler.
	AfterFunc(d time.Duration, f Handler) *TimerTask
}

// NewTimingWheel creates an instance of TimingWheel with the given tick and wheelSize.
func NewTimingWheel(tick time.Duration, wheelSize int64) (TimingWheel, error) {
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		return nil, ErrInvalidTickValue
	}
	if wheelSize <= 0 {
		return nil, ErrInvalidWheelSize
	}

	startMs := timeToMs(time.Now())
	t := newTimingWheel(tickMs, wheelSize, startMs, delayqueue.New(int(wheelSize)))

	ctx, cancel := context.WithCancel(context.Background())
	t.ctx = ctx
	t.cancel = cancel
	return t, nil
}

// newTimingWheel is the interval implement of creates timewheel instance.
// it always used when add timertask into timingwheel for creates overflow timingwheel.
func newTimingWheel(tickMs int64, wheelSize int64, startMs int64, dq delayqueue.DelayQueue) *timingWheel {
	buckets := make([]*bucket, wheelSize)
	for i := range buckets {
		buckets[i] = newBucket()
	}

	return &timingWheel{
		tick:        tickMs,
		wheelSize:   wheelSize,
		interval:    tickMs * wheelSize,
		currentTime: truncate(startMs, tickMs),
		buckets:     buckets,
		dq:          dq,
	}
}

// timingWheel is an implemention of TimingWheel
type timingWheel struct {
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
	overflowWheel unsafe.Pointer

	// dq is the queue of bucket expiration
	dq delayqueue.DelayQueue

	// wg for wait sub goroutine
	wg conc.WaitGroupWrapper

	// ctx to cancel sub goroutine
	ctx    context.Context
	cancel func()
}

func (w *timingWheel) Start() {
	w.wg.Wrap(func() {
		w.dq.Poll(w.ctx)
	})

	w.wg.Wrap(func() {
		for {
			select {
			case elem := <-w.dq.Chan():
				b := elem.(*bucket)
				w.advanceClock(b.Expiration())
				b.Flush(w.addOrRun)
			case <-w.ctx.Done():
				return
			}
		}
	})
}

func (w *timingWheel) Stop() {
	w.cancel()
	w.wg.Wait()
}

func (w *timingWheel) AfterFunc(d time.Duration, f Handler) *TimerTask {
	t := &TimerTask{
		expiration: timeToMs(time.Now().Add(d)),
		f:          f,
	}
	w.addOrRun(t)
	return t
}

func (w *timingWheel) addOrRun(t *TimerTask) {
	if !w.add(t) {
		// the timertask already expired, wo we run execute the timer's taks in its own goroutine.
		go t.f()
	}
}

func (w *timingWheel) add(t *TimerTask) bool {
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
			w.dq.Offer(b, b.Expiration())
		}
		return true
	default:
		overflowWheel := atomic.LoadPointer(&w.overflowWheel)
		if overflowWheel == nil {
			atomic.CompareAndSwapPointer(
				&w.overflowWheel,
				nil,
				unsafe.Pointer(newTimingWheel(
					w.interval,
					w.wheelSize,
					w.currentTime,
					w.dq,
				)),
			)
			overflowWheel = atomic.LoadPointer(&w.overflowWheel)
		}
		return (*timingWheel)(overflowWheel).add(t)
	}
}

func (w *timingWheel) advanceClock(expiration int64) {
	if expiration >= w.currentTime+w.tick {
		w.currentTime = truncate(expiration, w.tick)

		overflowWheel := atomic.LoadPointer(&w.overflowWheel)
		if overflowWheel != nil {
			(*timingWheel)(overflowWheel).advanceClock(w.currentTime)
		}
	}
}
