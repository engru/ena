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

package delayqueue

import (
	"context"
	"sync"
	"testing"
)

func Benchmark_DelayQueue(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	dq := New(0)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		dq.Poll(ctx)
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-dq.Chan():
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		dq.Offer(i, defaultTimer.Now()+int64(i))
	}
	cancel()
	wg.Wait()
}

func Benchmark_DelayQueue_Parallel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	dq := New(0)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		dq.Poll(ctx)
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-dq.Chan():
			}
		}
	}()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			dq.Offer(0, defaultTimer.Now())
		}
	})

	cancel()
	wg.Wait()
}
