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
	metricStorage domain.MetricStorage
	writeInterval time.Duration
	writer        io.Writer
	logger        zerolog.Logger
}

func New(
	metricStorage domain.MetricStorage,
	writeInterval time.Duration,
	writer io.Writer,
	logger zerolog.Logger,
) *MetricStorageDumper {
	dumper := &MetricStorageDumper{
		metricStorage: metricStorage,
		writeInterval: writeInterval,
		writer:        writer,
		logger:        logger,
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

	if m.writeInterval > 0 {
		go func() {
			time.Sleep(m.writeInterval)
			go m.saveCurrentMetrics()
		}()
	}

	return nil
}

func (m MetricStorageDumper) saveCurrentMetrics() {
	allMetrics := m.metricStorage.GetAllMetrics()
	m.write(allMetrics)

	if m.writeInterval > 0 {
		time.Sleep(m.writeInterval)
		go m.saveCurrentMetrics()
	}
}

func (m MetricStorageDumper) write(metrics []domain.Metric) {
	data, err := json.Marshal(metrics)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to marshal metrics")
		return
	}

	_, err = m.writer.Write(data)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to write metrics")
		return
	}
}
