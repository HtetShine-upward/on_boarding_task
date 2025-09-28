package main

import (
	"net"
	"os"
	"strings"

	"example.com/memo/server/internal/log"
	grpcTransport "example.com/memo/server/internal/transport/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	// 空文字ならデフォルト値を使う
	addr := os.Getenv("GRPC_ADDR")
	if strings.TrimSpace(addr) == "" {
		addr = "0.0.0.0:50051"
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if strings.TrimSpace(logLevel) == "" {
		logLevel = "info"
	}

	logger := log.New(logLevel)
	defer logger.Sync()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcTransport.LoggingInterceptor(logger)),
	)

	grpcTransport.RegisterMemoService(grpcServer)

	logger.Info("server started", zap.String("addr", ":50051"))

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
