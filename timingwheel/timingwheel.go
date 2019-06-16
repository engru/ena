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

// event is the representation of value in the wch
type event struct {
	// the identify of the event
	Type eventType

	t *timerTask
}

func (tw *timingWheel) Start() {
	tw.wg.Wrap(func() {
		tw.dq.Poll(tw.ctx)
	})

	addOrRun := func(t *timerTask) {
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
					if stopped {
						e.t.stopped = 1
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

func (tw *timingWheel) AfterFunc(d time.Duration, f Handler) (TimerTask, error) {
	wid := atomic.AddUint64(&tw.wid, 1)

	t := &timerTask{
		expiration: timeToMs(time.Now().Add(d)),
		f:          f,
		id:         wid,
		w:          tw,
	}

	// TODO(lsytj0413): deal the err
	outch, err := tw.wt.Register(wid)
	if err != nil {
		return nil, err
	}
	tw.wch <- event{
		Type: eventAddNew,
		t:    t,
	}

	v := <-outch
	return v.(*timerTask), nil
}

// TODO(lsytj0413): implement it
func (tw *timingWheel) TickFunc(d time.Duration, f Handler) (TimerTask, error) {
	return nil, nil
}

func (tw *timingWheel) StopFunc(t *timerTask) (bool, error) {
	outch, err := tw.wt.Register(t.id)
	if err != nil {
		return false, err
	}

	tw.wch <- event{
		Type: eventDelete,
		t:    t,
	}

	v := <-outch
	return v.(bool), nil
}
