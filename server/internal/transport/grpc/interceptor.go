package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type contextKey string

const RequestIDKey contextKey = "req_id"

func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		reqID := uuid.NewString()
		ctx = context.WithValue(ctx, RequestIDKey, reqID)

		resp, err := handler(ctx, req)
		latency := time.Since(start).Milliseconds()
		code := status.Code(err)

		logger.Info("RPC completed",
			zap.String("rpc", info.FullMethod),
			zap.String("req_id", reqID),
			zap.Int64("latency_ms", latency),
			zap.String("status", code.String()),
		)
		return resp, err
	}
}
