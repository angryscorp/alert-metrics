package metricstorage

import (
	"context"
	"encoding/json"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type FileMetricStorage struct {
	storage            domain.MetricStorage
	logger             zerolog.Logger
	writeInterval      time.Duration
	fileStoragePath    string
	fileStorageContext context.Context
}

var _ domain.MetricStorage = (*FileMetricStorage)(nil)

func NewFileMetricStorage(
	storage domain.MetricStorage,
	logger zerolog.Logger,
	writeInterval time.Duration,
	fileStoragePath string,
	shouldRestore bool,
) *FileMetricStorage {
	m := &FileMetricStorage{
		storage:            storage,
		logger:             logger,
		writeInterval:      writeInterval,
		fileStoragePath:    fileStoragePath,
		fileStorageContext: context.TODO(),
	}

	if shouldRestore {
		m.RestoreFromFile()
	}

	go m.saveCurrentMetrics()

	return m
}

func (s FileMetricStorage) RestoreFromFile() {
	data, err := os.ReadFile(s.fileStoragePath)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to read metrics file")
		return
	}

	var metrics []domain.Metric
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to unmarshal metrics file")
		return
	}

	for _, v := range metrics {
		if err := s.storage.UpdateMetric(s.fileStorageContext, v); err != nil {
			s.logger.Error().Err(err).Msg("failed to update metric")
			return
		}
	}
}

func (s FileMetricStorage) GetAllMetrics(ctx context.Context) []domain.Metric {
	return s.storage.GetAllMetrics(ctx)
}

func (s FileMetricStorage) GetMetric(ctx context.Context, metricType domain.MetricType, metricName string) (domain.Metric, bool) {
	return s.storage.GetMetric(ctx, metricType, metricName)
}

func (s FileMetricStorage) UpdateMetric(ctx context.Context, metric domain.Metric) error {
	err := s.storage.UpdateMetric(ctx, metric)
	if err != nil {
		return err
	}

	if s.writeInterval == 0 {
		s.saveCurrentMetrics()
	}

	return nil
}

func (s FileMetricStorage) UpdateMetrics(ctx context.Context, metrics []domain.Metric) error {
	err := s.storage.UpdateMetrics(ctx, metrics)
	if err != nil {
		return err
	}

	if s.writeInterval == 0 {
		s.saveCurrentMetrics()
	}

	return nil
}

func (s FileMetricStorage) Ping(ctx context.Context) error {
	return nil
}

func (s FileMetricStorage) saveCurrentMetrics() {
	allMetrics := s.storage.GetAllMetrics(s.fileStorageContext)
	s.writeToFile(allMetrics)

	if s.writeInterval > 0 {
		time.Sleep(s.writeInterval)
		go s.saveCurrentMetrics()
	}
}

func (s FileMetricStorage) writeToFile(metrics []domain.Metric) {
	data, err := json.Marshal(metrics)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal metrics")
		return
	}

	file, err := os.Create(s.fileStoragePath)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create metrics file")
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = file.Write(data)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to write metrics to file")
		return
	}
}
