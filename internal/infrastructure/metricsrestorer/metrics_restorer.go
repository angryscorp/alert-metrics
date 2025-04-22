package metricsrestorer

import (
	"encoding/json"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/rs/zerolog"
	"os"
)

type MetricsRestorer struct {
	fileStoragePath string
	logger          zerolog.Logger
}

func New(fileStoragePath string, logger zerolog.Logger) MetricsRestorer {
	return MetricsRestorer{
		fileStoragePath: fileStoragePath,
		logger:          logger,
	}
}

func (m MetricsRestorer) Restore() *[]domain.Metric {
	data, err := os.ReadFile(m.fileStoragePath)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to read metrics file")
		return nil
	}

	var result []domain.Metric
	err = json.Unmarshal(data, &result)
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to unmarshal metrics file")
		return nil
	}

	return &result
}
