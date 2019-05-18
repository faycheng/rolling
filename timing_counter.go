package rolling

import (
	"fmt"
	"time"
)

var _ Metric = &rollingCounter{}
var _ Aggregation = &rollingCounter{}

type TimingCounter interface {
	Metric
	Aggregation
	// Reduce applies the reduction function to all buckets within the window.
	Reduce(func(Iterator) float64) float64
}

// TimingCounterOpts contains the arguments for creating TimingCounter.
type TimingCounterOpts struct {
	Size           int
	BucketDuration time.Duration
}

type timingCounter struct {
	policy *TimingPolicy
}

func NewTimingCounter(opts TimingCounterOpts) TimingCounter {
	window := NewWindow(WindowOpts{Size: opts.Size})
	policy := NewTimingPolicy(window, TimingPolicyOpts{BucketDuration: opts.BucketDuration})
	return &timingCounter{
		policy: policy,
	}
}

func (t *timingCounter) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("rolling: cannot decrease in value. val: %d", val))
	}
	t.policy.Add(float64(val))
}

func (t *timingCounter) Reduce(f func(Iterator) float64) float64 {
	return t.policy.Reduce(f)
}

func (t *timingCounter) Avg() float64 {
	return t.policy.Reduce(Avg)
}

func (t *timingCounter) Min() float64 {
	return t.policy.Reduce(Min)
}

func (t *timingCounter) Max() float64 {
	return t.policy.Reduce(Max)
}

func (t *timingCounter) Sum() float64 {
	return t.policy.Reduce(Sum)
}

func (t *timingCounter) Value() int64 {
	return int64(t.Sum())
}
