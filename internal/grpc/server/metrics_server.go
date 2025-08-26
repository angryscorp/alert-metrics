package server

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/angryscorp/alert-metrics/internal/grpc/mapper"

	"github.com/angryscorp/alert-metrics/internal/domain"
	grpcmetrics "github.com/angryscorp/alert-metrics/internal/grpc/metrics"
)

type MetricsServer struct {
	grpcmetrics.UnimplementedMetricsServiceServer
	storage domain.MetricStorage
	logger  zerolog.Logger
}

var _ grpcmetrics.MetricsServiceServer = (*MetricsServer)(nil)

func NewMetricsServer(storage domain.MetricStorage, logger zerolog.Logger) *MetricsServer {
	return &MetricsServer{
		storage: storage,
		logger:  logger,
	}
}

func (s *MetricsServer) ReportRawMetric(ctx context.Context, req *grpcmetrics.ReportRawMetricRequest) (*grpcmetrics.Empty, error) {
	s.logger.Debug().
		Str("key", req.Key).
		Str("value", req.Value).
		Str("type", req.MetricType.String()).
		Msg("received raw metric via gRPC")

	metricType := mapper.MetricTypeToDomain(req.MetricType)

	metric, err := domain.NewMetrics(string(metricType), req.Key, req.Value)
	if err != nil {
		s.logger.Error().Err(err).
			Str("key", req.Key).
			Str("value", req.Value).
			Msg("failed to create metric")
		return &grpcmetrics.Empty{}, err
	}

	if err := s.storage.UpdateMetric(ctx, *metric); err != nil {
		s.logger.Error().Err(err).
			Str("key", req.Key).
			Msg("failed to update raw metric")
		return &grpcmetrics.Empty{}, err
	}

	return &grpcmetrics.Empty{}, nil
}

func (s *MetricsServer) ReportMetric(ctx context.Context, req *grpcmetrics.ReportMetricRequest) (*grpcmetrics.Empty, error) {
	s.logger.Debug().
		Str("metric_id", req.Metric.Id).
		Str("type", req.Metric.Type.String()).
		Msg("received metric via gRPC")

	metric := mapper.MetricToDomain(req.Metric)

	if err := s.storage.UpdateMetric(ctx, metric); err != nil {
		s.logger.Error().Err(err).
			Str("metric_id", req.Metric.Id).
			Msg("failed to update metric")
		return &grpcmetrics.Empty{}, err
	}

	return &grpcmetrics.Empty{}, nil
}

func (s *MetricsServer) ReportBatch(ctx context.Context, req *grpcmetrics.ReportBatchRequest) (*grpcmetrics.Empty, error) {
	s.logger.Debug().
		Int("count", len(req.Metrics)).
		Msg("received batch via gRPC")

	metrics := make([]domain.Metric, len(req.Metrics))
	for i, protoMetric := range req.Metrics {
		metrics[i] = mapper.MetricToDomain(protoMetric)
	}

	if err := s.storage.UpdateMetrics(ctx, metrics); err != nil {
		s.logger.Error().Err(err).
			Int("count", len(metrics)).
			Msg("failed to update batch")
		return &grpcmetrics.Empty{}, err
	}

	return &grpcmetrics.Empty{}, nil
}
