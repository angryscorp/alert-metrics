package server

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func loggingInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		logger.Info().
			Str("method", info.FullMethod).
			Msg("gRPC request started")

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Dur("duration", duration).
				Msg("gRPC request failed")
		} else {
			logger.Info().
				Str("method", info.FullMethod).
				Dur("duration", duration).
				Msg("gRPC request completed")
		}

		return resp, err
	}
}
