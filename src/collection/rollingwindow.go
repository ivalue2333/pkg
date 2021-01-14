package collection

import (
	"fmt"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
)

type (
	// RollingWindowOption let callers customize the RollingWindow.
	RollingWindowOption func(rollingWindow *RollingWindow)

	// RollingWindow defines a rolling window to calculate the events in buckets with time interval.
	RollingWindow struct {
		lock          sync.RWMutex
		size          int // 有多少个小窗口
		window        *window
		interval      time.Duration // 每个窗口的时间大小，例如 100 ms，表示一个小窗口记录 100ms 内的计数
		offset        int           // offset 表示当前应该更新或者读取哪个 bucket
		ignoreCurrent bool
		lastTime      time.Duration // start time of the last bucket
	}
)

// NewRollingWindow returns a RollingWindow that with size buckets and time interval,
// use opts to customize the RollingWindow.
func NewRollingWindow(size int, interval time.Duration, opts ...RollingWindowOption) *RollingWindow {
	if size < 1 {
		panic("size must be greater than 0")
	}

	w := &RollingWindow{
		size:     size,
		window:   newWindow(size),
		interval: interval,
		lastTime: timex.Now(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Add adds value to current bucket.
func (rw *RollingWindow) Add(v float64) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	// 确定当前应该是哪个小窗口（offset）， 顺便更新最近一次访问时间，方便下次计算 span
	rw.updateOffset()
	// 在这个小窗口计数
	rw.window.add(rw.offset, v)
}

// Reduce runs fn on all buckets, ignore current bucket if ignoreCurrent was set.
func (rw *RollingWindow) Reduce(fn func(b *Bucket)) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int
	span := rw.span()
	// ignore current bucket, because of partial data
	// diff 实际上是 rw.size -1 or rw.size - span
	if span == 0 && rw.ignoreCurrent {
		diff = rw.size - 1
	} else {
		diff = rw.size - span
	}
	if diff > 0 {
		// 计算当前 offset 下标
		offset := (rw.offset + span + 1) % rw.size

		// span 这几个 bucket 就不参与计算，因为已经过期了： 为了实现这个效果， 需要计算新的 rw.offset 和总共的 bucket的个数, rw.size - span
		fmt.Println("span", span, "last.offset", rw.offset, "cur.offset", offset, "count", diff)

		rw.window.reduce(offset, diff, fn)
	}
}

// 上次计数到当次，经过了 几个 rw.internal
func (rw *RollingWindow) span() int {
	offset := int(timex.Since(rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	} else {
		return rw.size
	}
}

func (rw *RollingWindow) updateOffset() {
	span := rw.span()
	if span <= 0 {
		return
	}

	// reset expired buckets
	for i := 0; i < span; i++ {
		//fmt.Println("reset", (rw.offset+i+1)%rw.size, "old offset", rw.offset, span, (rw.offset + span) % rw.size)
		rw.window.resetBucket((rw.offset + i + 1) % rw.size)
	}

	// 当前应该更新的 bucket 的下标 = 上次的下标 + 两次的间隔
	rw.offset = (rw.offset + span) % rw.size
	now := timex.Now()
	// align to interval time boundary
	rw.lastTime = now - (now-rw.lastTime)%rw.interval
}

// Bucket defines the bucket that holds sum and num of additions.
type Bucket struct {
	Sum   float64
	Count int64
}

func (b *Bucket) add(v float64) {
	b.Sum += v
	b.Count++
}

func (b *Bucket) reset() {
	b.Sum = 0
	b.Count = 0
}

type window struct {
	buckets []*Bucket
	size    int
}

func newWindow(size int) *window {
	buckets := make([]*Bucket, size)
	for i := 0; i < size; i++ {
		buckets[i] = new(Bucket)
	}
	return &window{
		buckets: buckets,
		size:    size,
	}
}

func (w *window) add(offset int, v float64) {
	w.buckets[offset%w.size].add(v)
}

func (w *window) reduce(start, count int, fn func(b *Bucket)) {
	for i := 0; i < count; i++ {
		fn(w.buckets[(start+i)%w.size])
	}
}

func (w *window) resetBucket(offset int) {
	w.buckets[offset%w.size].reset()
}

func (w *window) print() {
	buckets := make([]Bucket, len(w.buckets))
	for i, x := range w.buckets {
		buckets[i] = *x
	}
	fmt.Println(buckets)
}

// IgnoreCurrentBucket lets the Reduce call ignore current bucket.
func IgnoreCurrentBucket() RollingWindowOption {
	return func(w *RollingWindow) {
		w.ignoreCurrent = true
	}
}
