package domain

import (
	"slices"
)

type MetricRepresentative struct {
	Type  MetricType
	Name  string
	Value string
}

func (m MetricRepresentative) String() string {
	return m.Name + " (" + string(m.Type) + ") = " + m.Value
}

type MetricRepresentatives []MetricRepresentative

func NewMetricRepresentatives(metrics []Metric) MetricRepresentatives {
	res := make(MetricRepresentatives, len(metrics))
	for i, metric := range metrics {
		res[i] = MetricRepresentative{
			Type:  metric.MType,
			Name:  metric.ID,
			Value: metric.StringValue(),
		}
	}
	return res
}

func (m MetricRepresentatives) SortByName() MetricRepresentatives {
	slices.SortFunc(m, func(a, b MetricRepresentative) int {
		if a.Name > b.Name {
			return 1
		} else if a.Name < b.Name {
			return -1
		}
		return 0
	})
	return m
}
