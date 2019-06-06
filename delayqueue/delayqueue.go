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

// Package delayqueue describe an delayqueue implemention.
package delayqueue

import (
	"context"
)

// DelayQueue is an blocking queue of *Delay* elements, the element
// can only been taken when its delay has expired. The head of the queue
// is the element whose delay expired most recent in the queue.
type DelayQueue interface {
	// Offer insert the element into the current DelayQueue,
	// if the expiration is blow the current min expiration, the item will
	// been fired first.
	Offer(elem interface{}, expireation int64)

	// Poll starts an infinite loop, it will continually waits for an element to
	// been fired, and send the element to the output Chan.
	Poll(ctx context.Context)

	// Chan return the output chan, when the element is fired the element
	// will send to the channel.
	Chan() <-chan interface{}
}

// delayQueue implement the DelayQueue interface
type delayQueue struct {
	C chan interface{}
	T Timer
}

// New construct a DelayQueue with the initial size
func New(size int) DelayQueue {
	return NewWithTimer(size, defaultTimer)
}

// NewWithTimer construct a DelayQueue with the initial size and Timer
func NewWithTimer(size int, t Timer) DelayQueue {
	return &delayQueue{
		C: make(chan interface{}),
		T: t,
	}
}

// Offer implement the DelayQueue.Offer
func (q *delayQueue) Offer(element interface{}, expireation int64) {

}

// Poll implement the DelayQueue.Pool
func (q *delayQueue) Poll(ctx context.Context) {

}

// Chan implement the DelayQueue.Chan
func (q *delayQueue) Chan() <-chan interface{} {
	return q.C
}
