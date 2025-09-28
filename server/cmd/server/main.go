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
	// 環境変数 GRPC_ADDR を取得。空白または未設定ならデフォルトアドレス "0.0.0.0:50051" を使用
	addr := os.Getenv("GRPC_ADDR")
	if strings.TrimSpace(addr) == "" {
		addr = "0.0.0.0:50051"
	}

	// 環境変数 LOG_LEVEL を取得。空白または未設定なら "info" レベルを使用
	logLevel := os.Getenv("LOG_LEVEL")
	if strings.TrimSpace(logLevel) == "" {
		logLevel = "info"
	}

	// ロガーを初期化（指定されたログレベルで zap.Logger を生成）
	logger := log.New(logLevel)
	defer logger.Sync() // ログのフラッシュ（バッファを出力）

	// TCPリスナーを指定アドレスで作成。失敗した場合はログを出力して終了
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	// gRPCサーバーを作成し、UnaryInterceptorとしてロギングインターセプターを設定
	grpcServer := grpc.NewServer(
		// single request and single response
		grpc.UnaryInterceptor(grpcTransport.LoggingInterceptor(logger)),
	)

	// MemoService を gRPC サーバーに登録
	grpcTransport.RegisterMemoService(grpcServer)

	// サーバー起動ログを出力
	logger.Info("server started", zap.String("addr", ":50051"))

	// gRPCサーバーを起動し、リクエストの受け付けを開始。失敗した場合はログを出力して終了
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
