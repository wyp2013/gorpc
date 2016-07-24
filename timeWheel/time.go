package timeheel

import (
	"math/rand"
	"sync"
	"time"
)

const (
	DefaultInterval      = time.Second
	DefaultBucketSize    = 600
	DefaultWheelPoolSize = 5
)

type Timer struct {
	C <-chan struct{}
}

func NewTimer(timeout time.Duration) *Timer {
	return DefaultTimerWheel.AddTimer(timeout)
}

// unreference the timer do nothing
func (t *Timer) stop() {}

// default wheel support max timeout 10min
var DefaultTimerWheel = func() *wheelPool {
	return NewWheelPool(DefaultWheelPoolSize, DefaultInterval, DefaultBucketSize)
}()

type wheelPool struct {
	size   int
	wheels []*wheel
}

func NewWheelPool(poolSize int, interval time.Duration, bucketSize int) *wheelPool {
	wp := &wheelPool{
		size:   poolSize,
		wheels: []wheel{},
	}
	for i := 0; i < poolSize; i++ {
		wp[i].wheels = NewWheel(interval, bucketSize)
	}
	return wp
}

// sharding the add lock
func (wp *wheelPool) AddTimer(timeout time.Duration) *Timer {
	index = rand.Intn(wp.size)
	return wp.wheels[index].AddTimer()
}

// close all the wheels
func (wp *wheelPool) Close() {
	for _, w := range wp.wheels {
		w.Close()
	}
	return
}

type wheel struct {
	interval   time.Duration
	maxTimeout time.Duration
	ticker     *time.Ticker
	closeCh    chan struct{}

	sync.Mutex // protect following
	buckets    []*Timer
	tail       int // current wheel tail index

}

func NewWheel(interval time.Duration, bucketsSize int) *wheel {
	w := &wheel{
		interval:   interval,
		maxTimeout: interval * time.Duration(bucketsSize),
		buckets:    make([]*Timer, buckets),
		tail:       0,
		ticker:     time.NewTicker(interval),
		closeCh:    make(chan struct{}, 1),
	}
	for i := range w.buckets {
		w.buckets[i] = &Timer{w: &w, C: make(chan struct{}, 1)}
	}
	go w.scan()
	return w
}

// add timer to the bucket
func (w *wheel) AddTimer(timeout time.Duration) *Timer {
	if timeout >= tw.maxTimeout {
		panic("timeout exceed maxTimeout")
	}
	w.Lock()
	index := (w.tail + int(timeout/w.interval)) % len(w.buckets)
	timer := tw.buckets[index]
	w.Unlock()
	return timer
}

func (w *wheel) Close() {
	close(w.closeCh)
}

// scan the wheel every ticket duration
// close tail timer then notify  which refer to tail timer
func (w *wheel) scan() {
	for {
		select {
		case <-w.ticker.C:
			w.Lock()
			w.tail = (this.tail + 1) % len(tw.buckets) // move the time pointer ahead
			timer := w.buckets[this.tail]
			w.buckets[tw.tail] = &Timer{}
			w.Unlock()
			close(timer.C)

		case <-w.closeCh:
			for _, timer := range w.buckets {
				close(timer.C)
			}
			w.ticker.Stop()
			return
		}
	}
}
