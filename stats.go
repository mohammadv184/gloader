package gloader

import (
	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/pkg/stats"
)

type MetricKey string

func (m MetricKey) String() string {
	return string(m)
}

const (
	MetricBufferSizeBytes            MetricKey = "buffer.size.bytes"
	MetricBufferLengthRows           MetricKey = "buffer.length.rows"
	MetricBufferTotalWriteLengthRows MetricKey = "buffer.totalWriteLength.rows"
	MetricBufferTotalReadLengthRows  MetricKey = "buffer.totalReadLength.rows"
)

type BufferObserverAdapter struct {
	s      *stats.Stats
	dcName string
}

var _ data.BufferObserver = &BufferObserverAdapter{}

func (b *BufferObserverAdapter) SizeChanged(size uint64) {
	b.s.MustGetGaugeCounter(MetricBufferLengthRows.String()).(stats.GaugeCounter).Set(int64(size), b.dcName)
}

func (b *BufferObserverAdapter) LengthChanged(l uint64) {
	b.s.MustGetGaugeCounter(MetricBufferLengthRows.String()).(stats.GaugeCounter).Set(int64(l), b.dcName)
}

func (b *BufferObserverAdapter) Write(n int) {
	b.s.MustGetSequentialCounter(MetricBufferTotalWriteLengthRows.String()).IncBy(int64(n), b.dcName)
}

func (b *BufferObserverAdapter) Read(n int) {
	b.s.MustGetSequentialCounter(MetricBufferTotalReadLengthRows.String()).IncBy(int64(n), b.dcName)
}

func NewBufferObserverAdapter(s *stats.Stats, dcName string) *BufferObserverAdapter {
	s.MustGetGaugeCounter(MetricBufferLengthRows.String()).Set(0, dcName)
	s.MustGetGaugeCounter(MetricBufferSizeBytes.String()).Set(0, dcName)

	return &BufferObserverAdapter{
		s:      s,
		dcName: dcName,
	}
}

func NewStats() *stats.Stats {
	s := stats.New()

	s.RegisterGaugeCounter(
		MetricBufferSizeBytes.String(),
		"buffer size in bytes",
	)

	s.RegisterGaugeCounter(
		MetricBufferLengthRows.String(),
		"buffer size in rows",
	)

	s.RegisterSequentialCounter(
		MetricBufferTotalWriteLengthRows.String(),
		"total buffer writes in rows",
	)

	s.RegisterSequentialCounter(
		MetricBufferTotalReadLengthRows.String(),
		"total buffer reads in rows",
	)

	return s
}
