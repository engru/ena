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

	"github.com/lsytj0413/ena/conc"
	"github.com/lsytj0413/ena/conc/wait"
	"github.com/lsytj0413/ena/delayqueue"
)

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
	t := newWheel(tickMs, wheelSize, startMs)

	tw := &timingWheel{
		dq:  delayqueue.New(int(wheelSize)),
		w:   t,
		wt:  wait.New(),
		wch: make(chan event, int(wheelSize)*100), // the channel is bufferd, could change to unbufferd?
	}

	ctx, cancel := context.WithCancel(context.Background())
	tw.ctx = ctx
	tw.cancel = cancel
	return tw, nil
}

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

// timingWheel is an implemention of TimingWheel
type timingWheel struct {
	// the first layer wheel
	w *wheel

	// the Wait to get reponse when call AfterFunc
	wt wait.Wait

	// the Wait register id, incr
	wid uint64

	// wch is the channel which TimerTask putin when call AfterFunc
	wch chan event

	// dq is the queue of bucket expiration
	dq delayqueue.DelayQueue

	// wg for wait sub goroutine
	wg conc.WaitGroupWrapper

	// ctx to cancel sub goroutine
	ctx    context.Context
	cancel func()
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

// eventType is the representation of event, such as AddNew, RePost
type eventType = string

// eventAddNew is the identify when timertask is add from AfterFunc
var eventAddNew eventType = "AddNew"

// eventDelete is the identify
var eventDelete eventType = "Delete"

// event is the representation of value in the wch
type event struct {
	// the identify of the event
	Type eventType

	t *TimerTask
}

func (tw *timingWheel) Start() {
	tw.wg.Wrap(func() {
		tw.dq.Poll(tw.ctx)
	})

	addOrRun := func(t *TimerTask) {
		tw.w.addOrRun(t, tw.dq)
	}

	tw.wg.Wrap(func() {
		for {
			select {
			case elem := <-tw.dq.Chan():
				b := elem.(*bucket)
				tw.w.advanceClock(b.Expiration())

				// TODO(lsytj0413): race condition here?
				b.Flush(addOrRun)
			case e := <-tw.wch:
				switch e.Type {
				case eventAddNew:
					// an timer task is add from AfterFunc/StopFunc
					addOrRun(e.t)
					tw.wt.Trigger(e.t.id, e.t)
				case eventDelete:
					stopped := false
					for b := e.t.bucket(); b != nil; b = e.t.bucket() {
						stopped = b.remove(e.t)
					}
					tw.wt.Trigger(e.t.id, stopped)
				}
			case <-tw.ctx.Done():
				return
			}
		}
	})
}

func (tw *timingWheel) Stop() {
	tw.cancel()
	tw.wg.Wait()
}

func (tw *timingWheel) AfterFunc(d time.Duration, f Handler) *TimerTask {
	wid := atomic.AddUint64(&tw.wid, 1)

	t := &TimerTask{
		expiration: timeToMs(time.Now().Add(d)),
		f:          f,
		id:         wid,
		w:          tw,
	}

	// TODO(lsytj0413): deal the err
	outch, _ := tw.wt.Register(wid)
	tw.wch <- event{
		Type: eventAddNew,
		t:    t,
	}

	v := <-outch
	return v.(*TimerTask)
}

func (tw *timingWheel) StopFunc(t *TimerTask) bool {
	outch, _ := tw.wt.Register(t.id)
	tw.wch <- event{
		Type: eventDelete,
		t:    t,
	}

	v := <-outch
	return v.(bool)
}

func (w *wheel) addOrRun(t *TimerTask, dq delayqueue.DelayQueue) {
	if !w.add(t, dq) {
		// the timertask already expired, wo we run execute the timer's task in its own goroutine.
		go t.f()
	}
}

func (w *wheel) add(t *TimerTask, dq delayqueue.DelayQueue) bool {
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
