package metricdumper

import (
	"encoding/json"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/rs/zerolog"
	"io"
	"time"
)

var _ domain.MetricStorage = (*MetricStorageDumper)(nil)

type MetricStorageDumper struct {
	metricStorage          domain.MetricStorage
	writeIntervalInSeconds time.Duration
	writer                 io.Writer
	logger                 zerolog.Logger
}

func New(
	metricStorage domain.MetricStorage,
	writeIntervalInSeconds time.Duration,
	writer io.Writer,
	logger zerolog.Logger,
) *MetricStorageDumper {
	dumper := &MetricStorageDumper{
		metricStorage:          metricStorage,
		writeIntervalInSeconds: writeIntervalInSeconds,
		writer:                 writer,
		logger:                 logger,
	}

	go dumper.saveCurrentMetrics()

	return dumper
}

func (m MetricStorageDumper) GetMetric(metricType domain.MetricType, metricName string) (domain.Metric, bool) {
	return m.metricStorage.GetMetric(metricType, metricName)
}

func (m MetricStorageDumper) GetAllMetrics() []domain.Metric {
	return m.metricStorage.GetAllMetrics()
}

func (m MetricStorageDumper) UpdateMetric(metric domain.Metric) error {
	err := m.metricStorage.UpdateMetric(metric)
	if err != nil {
		return err
	}

	if m.writeIntervalInSeconds == 0 {
		m.write(metric)
	}

	return nil
}

func (m MetricStorageDumper) saveCurrentMetrics() {
	allMetrics := m.metricStorage.GetAllMetrics()
	for _, metric := range allMetrics {
		m.write(metric)
	}

	if m.writeIntervalInSeconds > 0 {
		time.Sleep(m.writeIntervalInSeconds)
		go m.saveCurrentMetrics()
	}
}

func (m MetricStorageDumper) write(metric domain.Metric) {
	data, err := json.Marshal(metric)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to marshal metric")
		return
	}

	_, err = m.writer.Write(data)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to write metric")
		return
	}

	_, err = m.writer.Write([]byte("\n"))
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to write separator")
		return
	}
}
