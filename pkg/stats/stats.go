// Package stats provide a simple real-time metric registry for monitoring and observability.
package stats

type Registry interface {
	RegisterSequentialCounter(name, description string) SequentialCounter
	RegisterGaugeCounter(name, description string) GaugeCounter
	GetSequentialCounter(name string) (SequentialCounter, error)
	GetGaugeCounter(name string) (GaugeCounter, error)
	MustGetSequentialCounter(name string) SequentialCounter
	MustGetGaugeCounter(name string) GaugeCounter
}

// Counter is a metric that accumulates values monotonically.
type Counter interface {
	// Tags returns the tags of the counter.
	Tags() []string
	// Value returns the current value of the counter.
	Value(tags ...[]string) int64
	// NotifyOnChange returns a channel that is closed when the counter value changes.
	NotifyOnChange(tags ...[]string) <-chan struct{}
}

// SequentialCounter is a counter that can only be incremented.
type SequentialCounter interface {
	Counter
	// Inc increments the counter by 1 and returns the new value.
	Inc(tags ...[]string) int64
	// IncBy increments the counter by delta and returns the new value.
	IncBy(delta int64, tags ...[]string) int64
}

// GaugeCounter is a counter that can be set to arbitrary values.
type GaugeCounter interface {
	Counter
	// Set sets the counter to the given value.
	Set(value int64, tags ...[]string)
	// Inc increments the counter by 1 and returns the new value.
	Inc(tags ...[]string) int64
	// IncBy increments the counter by delta and returns the new value.
	IncBy(delta int64, tags ...[]string) int64
	// Dec decrements the counter by 1 and returns the new value.
	Dec(tags ...[]string) int64
	// DecBy decrements the counter by delta and returns the new value.
	DecBy(delta int64, tags ...[]string) int64
}
