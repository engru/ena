# 定时器和时间轮系列(二): 优先级队列和延时队列

在上一篇文章中, 我们对标准库中的定时器相关功能进行了熟悉, 并且简单介绍了一种高性能的层级时间轮. 在本篇文章中, 将开始讲解如何参考 [Kafka Purgatory](https://www.confluent.io/blog/apache-kafka-purgatory-hierarchical-timing-wheels/) 去实现一个简单的层级时间轮.

## 整体结构

一个层级时间轮的实现中, 主要包括三个部分: 优先级队列, 延时队列以及时间轮自身. 它们三者的关系如下:

[Hierarchical Timing Wheel](./images/timingwheel.png)

其中:

+ 一个层级时间轮中有多层的时间轮, 每个轮盘上有多个插槽, 每个插槽都包含有相应时间范围内的定时器
+ 一个层级时间轮中有一个延时队列, 如果一个插槽中包含有定时器, 则该插槽的触发时间会被添加到延时队列, 并在该时间到期时从队列中弹出
+ 延时队列会包含有多个触发时间, 所有的触发时间都通过优先级队列进行排序

当一个定时任务被放置到层级时间轮之中时, 流程如下:

1. 按照任务的触发时间添加到某一层的时间轮的插槽中, 如果时间轮不存在则创建.
2. 如果插槽的触发时间发生了变化, 则把对应的插槽添加到延时队列
3. 延时队列在等待对应的触发时间之后, 将改元素弹出
4. 层级时间轮收到插槽弹出的时间之后, 遍历该插槽中的定时任务, 如果到达触发时间则触发改任务, 否则重新插入到层级时间轮

## 优先级队列

优先级队列是基于最小堆的实现, 最小堆是一种经过排序的二叉树结构, 它的特点是其中任何一个非终端节点的数据值均不会大于其左子节点和右子节点的数据值. 这样可以很容易的证明, 其中最小的值一定在根节点元素.

基于最小堆我们可以很方便的实现优先级队列, 只需要将优先级作为节点的数据值即可. 

[Priority Queue](https://he-s3.s3.amazonaws.com/media/uploads/2270b0f.jpg)

在标准库的 **container/heap** 包中, 实现了多个用于辅助最小堆实现的函数, 利用它们, 一个非常简易的优先级队列实现如下:

```
// Element for priority queue element
type Element struct {
	// Value for element
	Value interface{}

	// priority of the element, make it private to avoid change from element
	priority int64

	// index of the element in the slice
	index int

	// pq is the refer of priority queue
	pq *priorityQueue
}

// ...

// PriorityQueue for priority queue trait
type PriorityQueue interface {
	// Add element to the PriorityQueue, it will return the element witch been added
	Add(v interface{}, priority int64) *Element

	// Peek return the lowest priority element
	Peek() *Element

	// Pop return the lowest priority element and remove it
	Pop() *Element

	// Remove will remove the element from the priority queue
	Remove(v *Element) error

	// Update the element in the priority queue with the new priority
	Update(v *Element, priority int64) error

	// Size return the element size of queue
	Size() int
}

// elems is an slice of Elements, it implement the heap.Interface
type elems = []*Element

// heapi implement the heap.Interface
type heapi struct {
	pq *priorityQueue
}

// ...

// Push the value at the end of slice, implement heap.Push
func (h *heapi) Push(x interface{}) {
	h.pq.e = append(h.pq.e, x.(*Element))
}

// Pop the value at the last position of slice, implement heap.Pop
func (h *heapi) Pop() interface{} {
	old := h.pq.e
	n := len(old)
	if n == 0 {
		return nil
	}

	// set element to nil for GC
	x := old[n-1]
	old[n-1] = nil
	h.pq.e = old[0 : n-1]
	return x
}

// priorityQueue is a implement by min heap, the 0th element is the lowest value
type priorityQueue struct {
	e elems
	h *heapi
}

// NewPriorityQueue construct a PriorityQueue
func NewPriorityQueue(size int) PriorityQueue {
	pq := &priorityQueue{
		e: make(elems, 0, size),
	}
	pq.h = &heapi{
		pq: pq,
	}
	return pq
}

// Add element to the PriorityQueue, it will return the element witch been added
func (pq *priorityQueue) Add(x interface{}, priority int64) *Element {
	e := &Element{
		Value:    x,
		priority: priority,
		index:    len(pq.e),
		pq:       pq,
	}
	heap.Push(pq.h, e)
	return e
}

// Pop return the lowest priority element and remove it
func (pq *priorityQueue) Pop() *Element {
	if len(pq.e) == 0 {
		return nil
	}

	x := heap.Pop(pq.h)
	e := x.(*Element)
	e.index = -1
	e.pq = nil
	return e
}

// Remove will remove the element from the priority queue
func (pq *priorityQueue) Remove(e *Element) error {
	if e.pq != pq {
		return fmt.Errorf("PriorityQueue.Remove: QueueMatchFailed: Element[%v], Queue[%v]", e.pq, pq)
	}

	if e.index < 0 || e.index >= len(pq.e) {
		return fmt.Errorf("PriorityQueue.Remove: OutOfIndex: Index[%v], Len[%v]", e.index, len(pq.e))
	}
	if e.priority != pq.e[e.index].priority {
		return fmt.Errorf("PriorityQueue.Remove: PriorityMatchFailed: Element[%v], Queue[%v]", e.priority, pq.e[e.index].priority)
	}

	heap.Remove(pq.h, e.index)
	e.index = -1
	e.pq = nil
	return nil
}
```

1. 将所有元素放置在一个切片中, 模拟一个二叉树的结构
2. 将元素添加到切片之后, 调用 **heap.Push** 来完成最小堆结构的调整
3. 将元素从切片中弹出或者移除元素之后, 调用 **heap.Pop** 或者 **heap.Remove** 来完成最小堆结构的调整

## 延时队列

在有一个优先级队列的实现之后, 可以在这个结构之上快速的实现一个延时队列, 延时队列的接口定义如下:

```
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

	// Size return the element count in the queue
	Size() int
}
```

其中最重要的是 **Offer** 以及 **Poll** 两个函数, 以及 **Chan** 函数返回的通道. 我们可以通过 **Offer** 将延时元素放入队列, 并且通过 **Poll** 函数来启动延时队列, 最后从 **Chan** 函数返回的通道中读取触发的元素.

接下来我们看看 **Offer** 函数的实现:

```
// Offer implement the DelayQueue.Offer
func (q *delayQueue) Offer(element interface{}, expireation int64) {
	_push := func() (*Element, int) {
		q.mu.Lock()
		defer q.mu.Unlock()

		e := q.pq.Add(element, expireation)
		return e, e.Index()
	}
	_, index := _push()

	// there is no concurrent protection, EX:
	// 1. goroutine1 add element with expireation 100
	// 2. goroutine2 add element with expireation 50
	// 3. the both goroutine get the element index 0
	// 4. goroutine2 cas the sleeping state to 0, and send the wakeup signal
	// 5. pool wakeup and update the fired point, cas the sleeping state to 1
	// 6. goroutine1 cas the sleeping state to 0, and send the wakeup signal
	// 7. pool wakeup and update the fired point, cas the sleeping state to 1
	// because the pool always update the fired point to the min expireation, so there is no problem(always update to 50)
	if index == 0 {
		// the element is the first element(with the earliest expireation), we
		// need week up the Pool loop to update the fired point
		if atomic.CompareAndSwapInt32(&q.sleeping, 1, 0) {
			// if we change the sleeping state from sleep to weekup success, send the signal to wakepupC
			q.wakeupC <- struct{}{}
		}
	}
}
```

首先将元素添加到优先级队列中(根据触发时间进行排序), 并且如果是队列中的最优先的元素(具有最小的触发时间), 则通过一个内部的通道唤醒等待中的 goroutine.

接下来看 **Poll** 函数的实现(具体内容在 **pollImpl** 函数中):

```
// Poll implement the DelayQueue.Pool
func (q *delayQueue) Poll(ctx context.Context) {
	defer func() {
		// reset the state to wakeup
		atomic.StoreInt32(&q.sleeping, 0)
	}()

	// an infinite loop
	// 1. wakeup at the min expiration
	// 2. send to the C
	for poll(ctx, q) {
	}
}

var (
	poll func(ctx context.Context, q *delayQueue) bool
)

// the inner implement of poll, split from Poll for test
// return true if been wakeup or fired, false to shutdown the loop
func pollImpl(ctx context.Context, q *delayQueue) bool {
	n := q.T.Now()

	q.mu.Lock()
	item := q.pq.Peek()
	if item == nil || item.Priority() > n {
		// No item left, change the sleeping state to 1
		atomic.StoreInt32(&q.sleeping, 1)
	}
	q.mu.Unlock()

	// we have got the min expiration item, it maybe nil for empty pq
	if item == nil {
		// wait for wakeup (new item Offer into the queue)
		select {
		case <-ctx.Done():
			return false
		case <-q.wakeupC:
			return true
		}
	}

	// have item, wait for the fired point
	delta := item.Priority() - n
	if delta <= 0 {
		// the item need fired, send the value to the output channel
		select {
		case q.C <- item.Value:
			// the element is fired
			q.mu.Lock()
			_ = q.pq.Remove(item)
			q.mu.Unlock()
			return true
		case <-ctx.Done():
			return false
		}
	}

	// the item is pending, wait for fired or new min element add
	select {
	case <-q.wakeupC:
		return true
	case <-time.After(time.Duration(delta) * time.Millisecond):
		// we doesn't fired the item at there, go to next loop and the item will been fired because delta <= 0
		if atomic.SwapInt32(&q.sleeping, 0) == 0 {
			// if the old state is wakeup, the maybe an signal in wakeupC,
			// so we drain it the unblock the caller
			select {
			case <-q.wakeupC:
			default:
			}
		}
		return true
	case <-ctx.Done():
		return false
	}
}
```

整个 **Poll** 函数是一个无限循环, 只有当 **ctx** 被取消时候才会退出. 在循环中主要有以下几部分的逻辑:

1. 从优先级队列中取出队列头部的元素, 如果队列中元素为空则睡眠在内部通道上, 等待被 **Offer** 函数唤醒
2. 如果该元素已经到达触发时间, 则发送到返回通道中, 并且从优先级队列中移除该元素, 表示该元素已经触发
3. 如果该元素还未到达触发时间, 则通过 **time.After** 等待对应的超时时间, 然后进入下一轮循环; 或者被 **Offer** 中添加的新元素唤醒, 进入下一轮循环.

## 总结结构

在本篇文章中我们从层级时间轮的整体结构, 以及优先级队列和延时队列的简单实现中进行了一番探索. 在后续文章中会对层级时间轮自身如何实现进行进一步梳理.

## 参考资料

- [Hashed and Hierarchical Timing Wheels, 层级时间轮](http://www.cs.columbia.edu/~nahum/w6998/papers/ton97-timing-wheels.pdf)
- [Kafka Purgatory](https://www.confluent.io/blog/apache-kafka-purgatory-hierarchical-timing-wheels/)
- [DelayQueue](https://github.com/lsytj0413/ena/tree/master/delayqueue)