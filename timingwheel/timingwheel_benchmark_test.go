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
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/lsytj0413/ena/conc"
)

func Benchmark_TimingWheel_AfterFunc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg conc.WaitGroupWrapper
	tw, err := NewTimingWheel(WithTickDuration(time.Millisecond), WithSize(20))
	if err != nil {
		b.FailNow()
	}
	tw.Start()

	wg.Wrap(func() {
		<-ctx.Done()
		tw.Stop()
	})

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		tw.AfterFunc(
			time.Duration(rand.Intn(300))*time.Millisecond,
			func(time.Time) {},
		)
	}

	cancel()
	wg.Wait()
}

func Benchmark_TimingWheel_AfterFunc_Parallel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg conc.WaitGroupWrapper
	tw, err := NewTimingWheel(WithTickDuration(time.Millisecond), WithSize(20))
	if err != nil {
		b.FailNow()
	}
	tw.Start()

	wg.Wrap(func() {
		<-ctx.Done()
		tw.Stop()
	})

	rand.Seed(time.Now().UnixNano())

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tw.AfterFunc(
				time.Duration(rand.Intn(300))*time.Millisecond,
				func(time.Time) {},
			)
		}
	})

	cancel()
	wg.Wait()
}

func Benchmark_TimingWheel_TickFunc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg conc.WaitGroupWrapper
	tw, err := NewTimingWheel(WithTickDuration(time.Millisecond), WithSize(20))
	if err != nil {
		b.FailNow()
	}
	tw.Start()

	wg.Wrap(func() {
		<-ctx.Done()
		tw.Stop()
	})

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		t, _ := tw.TickFunc(
			time.Duration(rand.Intn(300)+1)*time.Millisecond,
			func(time.Time) {},
		)
		t.Stop()
	}

	cancel()
	wg.Wait()
}

func Benchmark_TimingWheel_TickFunc_Parallel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg conc.WaitGroupWrapper
	tw, err := NewTimingWheel(WithTickDuration(time.Millisecond), WithSize(20))
	if err != nil {
		b.FailNow()
	}
	tw.Start()

	wg.Wrap(func() {
		<-ctx.Done()
		tw.Stop()
	})

	rand.Seed(time.Now().UnixNano())

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			t, _ := tw.TickFunc(
				time.Duration(rand.Intn(300)+1)*time.Millisecond,
				func(time.Time) {},
			)
			t.Stop()
		}
	})

	cancel()
	wg.Wait()
}

func Benchmark_TimingWheel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg conc.WaitGroupWrapper
	tw, err := NewTimingWheel(WithTickDuration(time.Millisecond), WithSize(20))
	if err != nil {
		b.FailNow()
	}
	tw.Start()

	wg.Wrap(func() {
		<-ctx.Done()
		tw.Stop()
	})

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		t, _ := tw.TickFunc(
			time.Duration(rand.Intn(300)+1)*time.Millisecond,
			func(time.Time) {},
		)
		t.Stop()

		tw.AfterFunc(
			time.Duration(rand.Intn(300)+1)*time.Millisecond,
			func(time.Time) {},
		)
	}

	cancel()
	wg.Wait()
}

func Benchmark_TimingWheel_Parallel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg conc.WaitGroupWrapper
	tw, err := NewTimingWheel(WithTickDuration(time.Millisecond), WithSize(20))
	if err != nil {
		b.FailNow()
	}
	tw.Start()

	wg.Wrap(func() {
		<-ctx.Done()
		tw.Stop()
	})

	rand.Seed(time.Now().UnixNano())

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			t, _ := tw.TickFunc(
				time.Duration(rand.Intn(300)+1)*time.Millisecond,
				func(time.Time) {},
			)
			t.Stop()

			tw.AfterFunc(
				time.Duration(rand.Intn(300)+1)*time.Millisecond,
				func(time.Time) {},
			)
		}
	})

	cancel()
	wg.Wait()
}
