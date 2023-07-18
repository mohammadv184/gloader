// Package stats provide a simple real-time metric registry for monitoring and observability.
package stats

import (
	"fmt"
)

type Stats struct {
	registry Registry
}

func New() *Stats {
	return &Stats{
		registry: NewDefaultRegistry(),
	}
}

func (s *Stats) WithRegistry(registry Registry) *Stats {
	s.registry = registry
	return s
}

func (s *Stats) Registry() Registry {
	return s.registry
}

func (s *Stats) RegisterSequentialCounter(name, description string) SequentialCounter {
	m := NewDefaultSequentialCounter()
	s.registry.RegisterMetric(name, description, m)
	return m
}

func (s *Stats) RegisterGaugeCounter(name, description string) GaugeCounter {
	m := NewDefaultGaugeCounter()
	s.registry.RegisterMetric(name, description, m)
	return m
}

func (s *Stats) GetSequentialCounter(name string) (SequentialCounter, error) {
	m, err := s.registry.GetMetric(name)
	if err != nil {
		return nil, err
	}

	if m, ok := m.(SequentialCounter); ok {
		return m, nil
	}
	return nil, fmt.Errorf("stats: metric %s is not a SequentialCounter", name)
}

func (s *Stats) GetGaugeCounter(name string) (GaugeCounter, error) {
	m, err := s.registry.GetMetric(name)
	if err != nil {
		return nil, err
	}

	if m, ok := m.(GaugeCounter); ok {
		return m, nil
	}
	return nil, fmt.Errorf("stats: metric %s is not a GaugeCounter", name)
}

func (s *Stats) MustGetSequentialCounter(name string) SequentialCounter {
	m, err := s.registry.GetMetric(name)
	if err != nil {
		panic(err)
	}

	if m, ok := m.(SequentialCounter); ok {
		return m
	}
	panic(fmt.Errorf("stats: metric %s is not a SequentialCounter", name))
}

func (s *Stats) MustGetGaugeCounter(name string) GaugeCounter {
	m, err := s.registry.GetMetric(name)
	if err != nil {
		panic(err)
	}

	if m, ok := m.(GaugeCounter); ok {
		return m
	}
	panic(fmt.Errorf("stats: metric %s is not a GaugeCounter", name))
}
