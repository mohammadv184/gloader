package stats

import "fmt"

type Registry interface {
	RegisterMetric(name, description string, metric Metric)
	GetMetric(name string) (Metric, error)
	MustGetMetric(name string) Metric
	Description(name string) (string, error)
}

type DefaultRegistry struct {
	metrics        map[string]Metric
	descriptionMap map[string]string
}

func NewDefaultRegistry() *DefaultRegistry {
	return &DefaultRegistry{
		metrics:        make(map[string]Metric),
		descriptionMap: make(map[string]string),
	}
}

func (r *DefaultRegistry) RegisterMetric(name, description string, metric Metric) {
	r.metrics[name] = metric
	r.descriptionMap[name] = description
}

func (r *DefaultRegistry) GetMetric(name string) (Metric, error) {
	metric, isExists := r.metrics[name]
	if !isExists {
		return nil, fmt.Errorf("metric %s does not exist", name)
	}
	return metric, nil
}

func (r *DefaultRegistry) MustGetMetric(name string) Metric {
	metric, err := r.GetMetric(name)
	if err != nil {
		panic(err)
	}
	return metric
}

func (r *DefaultRegistry) Description(name string) (string, error) {
	description, isExists := r.descriptionMap[name]
	if !isExists {
		return "", fmt.Errorf("metric %s does not exist", name)
	}
	return description, nil
}
