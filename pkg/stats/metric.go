package stats

import (
	"sync/atomic"
)

// Metric is a metric that accumulates values monotonically.
type Metric interface {
	// Tags returns the tags of the Metric.
	Tags() []string
	// Value returns the current value of the counter.
	Value(tags ...string) int64
	// NotifyOnChange receives a channel that will be notified when the value of the counter changes.
	NotifyOnChange(ch chan<- any, tags ...string)
}

// SequentialCounter is a counter that can only be incremented.
type SequentialCounter interface {
	Metric
	// Inc increments the counter by 1 and returns the new value.
	Inc(tags ...string) int64
	// IncBy increments the counter by delta and returns the new value.
	IncBy(delta int64, tags ...string) int64
}

// GaugeCounter is a counter that can be set to arbitrary values.
type GaugeCounter interface {
	Metric
	// Set sets the counter to the given value.
	Set(value int64, tags ...string)
	// Inc increments the counter by 1 and returns the new value.
	Inc(tags ...string) int64
	// IncBy increments the counter by delta and returns the new value.
	IncBy(delta int64, tags ...string) int64
	// Dec decrements the counter by 1 and returns the new value.
	Dec(tags ...string) int64
	// DecBy decrements the counter by delta and returns the new value.
	DecBy(delta int64, tags ...string) int64
}

type DefaultSequentialCounter struct {
	count          *uint64
	tagCountersMap map[string]SequentialCounter
	notifyCh       chan<- any
}

func (c *DefaultSequentialCounter) Tags() []string {
	tags := make([]string, 0, len(c.tagCountersMap))
	for tag := range c.tagCountersMap {
		tags = append(tags, tag)
	}
	return tags
}

func (c *DefaultSequentialCounter) Value(tags ...string) int64 {
	if len(tags) > 0 {
		var count int64
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultSequentialCounter()
			}
			count += c.tagCountersMap[tag].Value()
		}
		return count
	}

	return int64(atomic.LoadUint64(c.count))
}

func (c *DefaultSequentialCounter) NotifyOnChange(ch chan<- any, tags ...string) {
	if len(tags) > 0 {
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultSequentialCounter()
			}
			c.tagCountersMap[tag].NotifyOnChange(ch)
		}
		return
	}

	c.notifyCh = ch
}

func (c *DefaultSequentialCounter) Inc(tags ...string) int64 {
	if len(tags) > 0 {
		var count int64
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultSequentialCounter()
			}
			count += c.tagCountersMap[tag].Inc()
		}
		return count
	}

	cv := int64(atomic.AddUint64(c.count, 1))
	c.notify(c.notifyCh)
	return cv
}

func (c *DefaultSequentialCounter) IncBy(delta int64, tags ...string) int64 {
	if len(tags) > 0 {
		var count int64
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultSequentialCounter()
			}
			count += c.tagCountersMap[tag].IncBy(delta)
		}
		return count
	}

	cv := int64(atomic.AddUint64(c.count, uint64(delta)))
	c.notify(c.notifyCh)
	return cv
}

func (c *DefaultSequentialCounter) notify(ch chan<- any) {
	if ch == nil {
		return
	}

	select {
	case ch <- "notify":
	default:
	}
}

func NewDefaultSequentialCounter() *DefaultSequentialCounter {
	return &DefaultSequentialCounter{
		count:          new(uint64),
		tagCountersMap: make(map[string]SequentialCounter),
	}
}

type DefaultGaugeCounter struct {
	count          *int64
	tagCountersMap map[string]GaugeCounter
	notifyCh       chan<- any
}

func (c *DefaultGaugeCounter) Tags() []string {
	tags := make([]string, 0, len(c.tagCountersMap))
	for tag := range c.tagCountersMap {
		tags = append(tags, tag)
	}
	return tags
}

func (c *DefaultGaugeCounter) Value(tags ...string) int64 {
	if len(tags) > 0 {
		var count int64
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultGaugeCounter()
			}
			count += c.tagCountersMap[tag].Value()
		}
		return count
	}

	return atomic.LoadInt64(c.count)
}

func (c *DefaultGaugeCounter) NotifyOnChange(ch chan<- any, tags ...string) {
	if len(tags) > 0 {
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultGaugeCounter()
			}
			c.tagCountersMap[tag].NotifyOnChange(ch)
		}
		return
	}

	c.notifyCh = ch
}

func (c *DefaultGaugeCounter) Set(value int64, tags ...string) {
	if len(tags) > 0 {
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultGaugeCounter()
			}
			c.tagCountersMap[tag].Set(value)
		}
	}

	atomic.StoreInt64(c.count, value)
	c.notify(c.notifyCh)
}

func (c *DefaultGaugeCounter) IncBy(delta int64, tags ...string) int64 {
	if len(tags) > 0 {
		var count int64
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultGaugeCounter()
			}
			count += c.tagCountersMap[tag].IncBy(delta)
		}
		return count
	}

	cv := atomic.AddInt64(c.count, delta)
	c.notify(c.notifyCh)
	return cv
}

func (c *DefaultGaugeCounter) DecBy(delta int64, tags ...string) int64 {
	if len(tags) > 0 {
		var count int64
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultGaugeCounter()
			}
			count += c.tagCountersMap[tag].DecBy(delta)
		}
		return count
	}

	cv := atomic.AddInt64(c.count, -delta)
	c.notify(c.notifyCh)
	return cv
}

func (c *DefaultGaugeCounter) Inc(tags ...string) int64 {
	if len(tags) > 0 {
		var count int64
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultGaugeCounter()
			}
			count += c.tagCountersMap[tag].Inc()
		}
		return count
	}

	cv := atomic.AddInt64(c.count, 1)
	c.notify(c.notifyCh)
	return cv
}

func (c *DefaultGaugeCounter) Dec(tags ...string) int64 {
	if len(tags) > 0 {
		var count int64
		for _, tag := range tags {
			if _, isExists := c.tagCountersMap[tag]; !isExists {
				c.tagCountersMap[tag] = NewDefaultGaugeCounter()
			}
			count += c.tagCountersMap[tag].Dec()
		}
		return count
	}

	cv := atomic.AddInt64(c.count, -1)
	c.notify(c.notifyCh)
	return cv
}

func (c *DefaultGaugeCounter) notify(ch chan<- any) {
	if ch == nil {
		return
	}

	select {
	case ch <- "notify":
	default:
	}
}

func NewDefaultGaugeCounter() *DefaultGaugeCounter {
	return &DefaultGaugeCounter{
		count:          new(int64),
		tagCountersMap: make(map[string]GaugeCounter),
	}
}
