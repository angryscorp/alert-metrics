package metricsrestorer

import (
	"encoding/json"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"os"
)

type MetricsRestorer struct {
	fileStoragePath string
}

func New(fileStoragePath string) MetricsRestorer {
	return MetricsRestorer{fileStoragePath: fileStoragePath}
}

func (m MetricsRestorer) Restore() (*[]domain.Metric, error) {
	data, err := os.ReadFile(m.fileStoragePath)
	if err != nil {
		return nil, err
	}

	var result []domain.Metric
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
