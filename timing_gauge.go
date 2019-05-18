package rolling

import (
	"fmt"
	"time"
)

var _ Metric = &rollingCounter{}
var _ Aggregation = &rollingCounter{}

type TimingGauge interface {
	Metric
	Aggregation
	// Reduce applies the reduction function to all buckets within the window.
	Reduce(func(Iterator) float64) float64
}

// TimingGaugeOpts contains the arguments for creating TimingGauge.
type TimingGaugeOpts struct {
	Size           int
	BucketDuration time.Duration
}

type timingGauge struct {
	policy *TimingPolicy
}

func NewTimingGauge(opts TimingGaugeOpts) TimingGauge {
	window := NewWindow(WindowOpts{Size: opts.Size})
	policy := NewTimingPolicy(window, TimingPolicyOpts{BucketDuration: opts.BucketDuration})
	return &timingGauge{
		policy: policy,
	}
}

func (t *timingGauge) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("rolling: cannot decrease in value. val: %d", val))
	}
	t.policy.Add(float64(val))
}

func (t *timingGauge) Reduce(f func(Iterator) float64) float64 {
	return t.policy.Reduce(f)
}

func (t *timingGauge) Avg() float64 {
	return t.policy.Reduce(Avg)
}

func (t *timingGauge) Min() float64 {
	return t.policy.Reduce(Min)
}

func (t *timingGauge) Max() float64 {
	return t.policy.Reduce(Max)
}

func (t *timingGauge) Sum() float64 {
	return t.policy.Reduce(Sum)
}

func (t *timingGauge) Value() int64 {
	return int64(t.Sum())
}
