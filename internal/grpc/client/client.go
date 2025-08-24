package client

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/angryscorp/alert-metrics/internal/grpc/mapper"

	"github.com/angryscorp/alert-metrics/internal/domain"
	grpcmetrics "github.com/angryscorp/alert-metrics/internal/grpc/metrics"
)

const contextTimeout = 5 * time.Second

type GRPCMetricReporter struct {
	client grpcmetrics.MetricsServiceClient
	conn   *grpc.ClientConn
	logger zerolog.Logger
}

func New(address string, logger zerolog.Logger) (*GRPCMetricReporter, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := grpcmetrics.NewMetricsServiceClient(conn)

	return &GRPCMetricReporter{
		client: client,
		conn:   conn,
		logger: logger,
	}, nil
}

func (gr *GRPCMetricReporter) ReportRawMetric(metricType domain.MetricType, key string, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	req := &grpcmetrics.ReportRawMetricRequest{
		MetricType: mapper.MetricTypeToProto(metricType),
		Key:        key,
		Value:      value,
	}

	_, err := gr.client.ReportRawMetric(ctx, req)
	if err != nil {
		gr.logger.Error().Err(err).Str("key", key).Str("value", value).Msg("failed to report raw metric via gRPC")
		return
	}

	gr.logger.Debug().Str("key", key).Str("value", value).Msg("raw metric reported via gRPC")
}

func (gr *GRPCMetricReporter) ReportMetric(metric domain.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	req := &grpcmetrics.ReportMetricRequest{
		Metric: mapper.MetricToProto(metric),
	}

	_, err := gr.client.ReportMetric(ctx, req)
	if err != nil {
		gr.logger.Error().Err(err).Str("metric_id", metric.ID).Msg("failed to report metric via gRPC")
		return
	}

	gr.logger.Debug().Str("metric_id", metric.ID).Msg("metric reported via gRPC")
}

func (gr *GRPCMetricReporter) ReportBatch(metrics []domain.Metric) {
	if len(metrics) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	protoMetrics := make([]*grpcmetrics.Metric, len(metrics))
	for i, metric := range metrics {
		protoMetrics[i] = mapper.MetricToProto(metric)
	}

	req := &grpcmetrics.ReportBatchRequest{
		Metrics: protoMetrics,
	}

	_, err := gr.client.ReportBatch(ctx, req)
	if err != nil {
		gr.logger.Error().Err(err).Int("count", len(metrics)).Msg("failed to report batch via gRPC")
		return
	}

	gr.logger.Debug().Int("count", len(metrics)).Msg("batch reported via gRPC")
}

func (gr *GRPCMetricReporter) Close() error {
	if gr.conn != nil {
		return gr.conn.Close()
	}
	return nil
}
