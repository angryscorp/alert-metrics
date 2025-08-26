package server

import (
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/angryscorp/alert-metrics/internal/domain"
	grpcmetrics "github.com/angryscorp/alert-metrics/internal/grpc/metrics"
)

// GRPCServer wraps the gRPC server with additional functionality
type GRPCServer struct {
	server *grpc.Server
	logger zerolog.Logger
}

func NewGRPCServer(storage domain.MetricStorage, logger zerolog.Logger) *GRPCServer {
	var opts []grpc.ServerOption

	opts = append(opts, grpc.UnaryInterceptor(loggingInterceptor(logger)))

	grpcServer := grpc.NewServer(opts...)
	metricsServer := NewMetricsServer(storage, logger)

	grpcmetrics.RegisterMetricsServiceServer(grpcServer, metricsServer)

	return &GRPCServer{
		server: grpcServer,
		logger: logger,
	}
}

func (gs *GRPCServer) Run(address string, shutdownCh <-chan struct{}) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	gs.logger.Info().Str("address", address).Msg("starting gRPC server")

	errCh := make(chan error, 1)
	go func() {
		if err := gs.server.Serve(listener); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-shutdownCh:
		gs.logger.Info().Msg("gracefully stopping gRPC server")
		gs.server.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}
