package rolling

import (
	"sync"
	"time"
)

type TimingPolicy struct {
	mu             sync.RWMutex
	size           int
	window         *Window
	offset         int
	stime          time.Time
	bucketDuration time.Duration
}

type TimingPolicyOpts struct {
	BucketDuration time.Duration
}

func NewTimingPolicy(window *Window, opts TimingPolicyOpts) *TimingPolicy {
	return &TimingPolicy{
		window:         window,
		size:           window.Size(),
		offset:         -1,
		bucketDuration: opts.BucketDuration,
	}
}

func (p *TimingPolicy) add(f func(offset int, val float64), val float64) {
	p.mu.Lock()
	if p.stime.IsZero() {
		p.stime = time.Now()
	}
	offset := int(time.Since(p.stime) / p.bucketDuration)
	f(offset, val)
	p.mu.Unlock()
}

func (p *TimingPolicy) Append(val float64) {
	p.add(p.window.Append, val)
}

func (p *TimingPolicy) Add(val float64) {
	p.add(p.window.Add, val)
}

// Reduce applies the reduction function to all buckets within the window.
func (p *TimingPolicy) Reduce(f func(Iterator) float64) (val float64) {
	p.mu.RLock()
	val = f(p.window.Iterator(0, p.size))
	p.mu.RUnlock()
	return val
}
