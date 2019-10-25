# 定时器和时间轮系列(一): 初识

在实现需求的过程中, 经常会遇到如下的一类需求:

- 在一个间隔时间之后做某事: 例如在最后一次消息发送的5分钟之后, 断开连接
- 在一个间隔时间之后不停的做某事: 例如每隔5分钟之后刷新内存中的缓存

使用定时器可以方便的实现上述功能. 定时器是一种结构, 它的主要作用是在一个给定的时间间隔之后, 调用一个给定的回调函数或者发出一个信号, 应用可以在回调函数或信号处理函数中实现相应的业务逻辑.

## 普通用法

在 Go 中, 标准库 **time** 包提供了一些基本的定时器相关操作, 常见的用法如下:

```
func delayOnce() {
	n := time.Now()
	fmt.Println("delayOnce start, ", n)

	// delay 1 second
	<-time.After(time.Second)
	fmt.Println("Cost ", time.Since(n))
}

func delayTicker() {
	n := time.Now()
	fmt.Println("delayTicker start, ", n)

	// ticker 1 second, 3 times
	t := time.NewTicker(time.Second)
	for i := 0; i < 3; i++ {
		<-t.C
		fmt.Println("Tick ", time.Since(n))
	}

	t.Stop()
	fmt.Println("Cost ", time.Since(n))
}
```

[示例代码1](https://play.golang.org/p/iotO9UeHxT4)

可以看到, 有几种方式来使用定时器:

1. 通过 time.NewTimer 创建一个定时器, 这样可以在循环中对定时器进行复用, 降低 runtime 的压力
2. 通过 time.After 得到一个 channel, 当该 channel 可读时即定时器到期触发
3. 通过 time.AfterFunc 在指定时间间隔后运行一个回调函数

上述三种都是一次性的定时器, 还有一种持续性的定时器(ticker):

1. 通过 time.NewTicker 创建一个定时器, 该定时器会周期性的向 channel 中发送信号(如果 channel 中还有未读取的信号则直接丢弃当次信号)

## 如何实现

Go 中的定时器代码经过多个版本的演进, 到现在正在开发中的 Go1.14 为止, 主要有三次大的变更.

### Go1.10 之前

在这个版本的实现中, Go 将所有的定时器都放在一个[最小堆](https://en.wikipedia.org/wiki/Binary_heap) 中, 并且在内部会启动一个 goroutine 持续的检查堆顶定时器是否已经到期, 如果到期则触发对应的回调函数.

创建定时器并添加到最小堆的主要代码如下:

```
// Add a timer to the heap and start or kick timerproc if the new timer is
// earlier than any of the others.
// Timers are locked.
func addtimerLocked(t *timer) {
    // ...
	t.i = len(timers.t)
	timers.t = append(timers.t, t)
	siftupTimer(t.i)
	if t.i == 0 {
		// siftup moved to top: new earliest deadline.
		if timers.sleeping {
			timers.sleeping = false
			notewakeup(&timers.waitnote)
		}
        // ...
	}
	if !timers.created {
		timers.created = true
		go timerproc()
	}
}
```

1. 将定时器放到一个内部的切片中
2. 调用 siftupTimer 调整最小堆的结构, 时间复杂度为 O(lg^n)
3. 如果当前的定时器是最接近的, 则唤醒等待的 goroutine
4. 如果是第一个定时器, 则启动内部的 goroutine

内部 goroutine 的循环代码如下:

```
// Timerproc runs the time-driven events.
// It sleeps until the next event in the timers heap.
// If addtimer inserts a new earlier event, it wakes timerproc early.
func timerproc() {
	timers.gp = getg()
	for {
        // ...
		for {
            // ...
			t := timers.t[0]
			delta = t.when - now
			if delta > 0 {
				break
			}
			if t.period > 0 {
				// leave in heap but adjust next time to fire
				t.when += t.period * (1 + -delta/t.period)
				siftdownTimer(0)
			} else {
				// remove from heap
				last := len(timers.t) - 1
				if last > 0 {
					timers.t[0] = timers.t[last]
					timers.t[0].i = 0
				}
				timers.t[last] = nil
				timers.t = timers.t[:last]
				if last > 0 {
					siftdownTimer(0)
				}
				t.i = -1 // mark as removed
			}
            // ...
			f(arg, seq)
		}
        // ...

		// At least one timer pending. Sleep until then.
		timers.sleeping = true
		timers.sleepUntil = now + delta
		notetsleepg(&timers.waitnote, delta)
	}
}
```

1. 获取堆顶的定时器, 如果到期了则触发回调, 并且如果该定时器是持续的则更新下次到期时间, 并调整最小堆; 如果不是则移除该定时器, 同时也调整一次最小堆
2. 没到期则等待被唤醒, 或者指定的时间间隔到达

### Go1.10 ~ Go1.13

在上面版本的实现里所有的定时器都放在一个最小堆中, 这样就有几个显而易见的缺点: 

1. 当有多个 P 同时运行的时候, 容易造成锁竞争(所有对最小堆的操作都是有锁保护的), 降低吞吐
2. 当有很多的定时器存在时, 最小堆的插入/删除效率也会降低

在这个阶段, 运行时针对上述的缺点进行了修改, 最主要的方法有如下两点:

1. 将所有定时器分布到 64 个最小堆中, 减小每个堆的数据量
2. 插入定时器时用 P 的 id 将其分布到不同的最小堆, 这样插入时就可以降低锁竞争

最主要的插入定时器的代码如下:

```
func (t *timer) assignBucket() *timersBucket {
	id := uint8(getg().m.p.ptr().id) % timersLen
	t.tb = &timers[id].timersBucket
	return t.tb
}

func addtimer(t *timer) {
	tb := t.assignBucket()
	lock(&tb.lock)
	ok := tb.addtimerLocked(t)
	unlock(&tb.lock)
	if !ok {
		badTimer()
	}
}
```

### Maybe Go1.14 及之后

上述版本在多 GPU 系统的性能仍然不够好, 主要是不是 cpu-scale 的, 见 [runtime: timer doesn't scale on multi-CPU systems with a lot of timers](https://github.com/golang/go/issues/15133).

Go 仍然在对定时器相关的代码进行进一步优化, 可以见 [runtime: make timers faster](https://github.com/golang/go/issues/6239). 主要的思路是:

1. 将每个定时器直接绑定到 P 上, 这样可以直接随着 P 扩展
2. 不再采用最小堆, 直接利用 netpoller 来让定时器的到期后直接得到通知

上述的优化仍然在开发的过程中, 可能会在 Go1.14 上发布, 也可能会在更后面的版本.

## 时间轮

### 为何要引入时间轮

1. 在定时器的数量增长到百万级之后, 基于最小堆实现的定时器的性能会显著降低, 需要一种更高效的实现
2. 在有些场景下的使用不是很方便

例如, 服务器维护有对客户端的连接, 并且定时在连接中发送心跳来确保连接的可用性, 一个普遍的实现方式如下:

```
func onConnect(ctx context.Context, i int) {
	t := time.NewTicker(time.Second)
	n := time.Now()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done")
			return
		case <-t.C:
			fmt.Printf("Tick[%v] %v\n", i, time.Since(n))
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		go onConnect(ctx, i)
	}

	<-ctx.Done()
	time.Sleep(time.Second) // wait all sub goroutine exit, should use WaitGroup
}
```

[示例代码2](https://play.golang.org/p/UaVsfXJDRfS)

但是这种方式中, 每一个连接就需要新增一个 goroutine, 并且对 goroutine 的清理也会比较复杂. 所以, 如果有一个独立的 goroutine 能够对这些定时任务进行触发, 操作上会方便很多.

### 时间轮

在 Kafka 中, 使用一个叫做 [Hashed and Hierarchical Timing Wheels, 层级时间轮](http://www.cs.columbia.edu/~nahum/w6998/papers/ton97-timing-wheels.pdf) 的高性能时间轮数据结构, 实现了自己的[时间轮](https://www.confluent.io/blog/apache-kafka-purgatory-hierarchical-timing-wheels/).

一个时间轮就是一个定时器容器, 该容器可以高效的管理定时器. 思路如下:

1. 轮盘上有多个插槽
2. 每个定时器都放置到合适的插槽中
3. 每次轮询时直接获取最早的插槽中的定时器并触发即可

在层级时间轮中, 将插槽分为多个层次, 每一层的时间轮的插槽范围都会扩大, 例如:

1. 第一层时间轮有20个插槽, 每个插槽为1秒, 那么第二层时间轮每个插槽为20秒, 第三层为400秒, 依次类推, 除第一层外都是按需创建
2. 当一个10秒的定时器插入时放置到第一层时间轮中, 100秒的定时器则放置到第二层时间轮
3. 随着时间的流逝, 高层时间轮中的定时任务会降级重新插入低层的时间轮, 直到触发位置
4. 每个插槽共享一个触发时间, 这样可以显著降低需要触发的事件的个数

一个示意图如下:

[Hashed and Hierarchical Timing Wheels](https://www.confluent.io/wp-content/uploads/2016/08/TimingWheels2.png)

## 总结

本篇文章为系列的第一篇, 主要介绍了 Go 标准库中对定时器的处理方式, 以及对层级时间轮进行了大致介绍. 后续文章会进一步介绍一个层级时间轮的简单实现, 并如何进行简单的测试与调优.

## 参考资料

- [How Do They Do It: Timers in Go](https://blog.gopheracademy.com/advent-2016/go-timers/)
- [Hashed and Hierarchical Timing Wheels, 层级时间轮](http://www.cs.columbia.edu/~nahum/w6998/papers/ton97-timing-wheels.pdf)
- [Kafka Purgatory](https://www.confluent.io/blog/apache-kafka-purgatory-hierarchical-timing-wheels/)
- [Timing Wheel](https://github.com/lsytj0413/ena/tree/master/timingwheel)