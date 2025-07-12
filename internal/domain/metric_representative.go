package domain

import (
	"slices"
)

// MetricRepresentative represents a simplified version of a metric with type, name, and value for easier manipulation.
type MetricRepresentative struct {
	Type  MetricType
	Name  string
	Value string
}

func (m MetricRepresentative) String() string {
	return m.Name + " (" + string(m.Type) + ") = " + m.Value
}

type MetricRepresentatives []MetricRepresentative

// NewMetricRepresentatives converts a slice of Metric into MetricRepresentatives for simplified representation and manipulation.
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

// SortByName sorts MetricRepresentatives by their Name field in ascending lexicographical order. Returns the sorted slice.
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
