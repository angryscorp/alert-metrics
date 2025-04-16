package domain

import (
	"errors"
	"strconv"
)

type Metrics struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
}

func NewMetrics(metricType string, metricName string, value string) (*Metrics, error) {
	mType, err := NewMetricType(metricType)
	if err != nil {
		return nil, err
	}

	result := Metrics{
		ID:    metricName,
		MType: mType,
	}

	switch mType {
	case MetricTypeCounter:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, errors.New("invalid counter value")
		}
		result.Delta = &v

	case MetricTypeGauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, errors.New("invalid gauge value")
		}
		result.Value = &v

	default:
		return nil, errors.New("unsupported metric type")
	}

	return &result, nil
}

func (m Metrics) StringValue() string {
	switch m.MType {
	case MetricTypeGauge:
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)

	case MetricTypeCounter:
		return strconv.FormatInt(*m.Delta, 10)

	default:
		return ""
	}
}
