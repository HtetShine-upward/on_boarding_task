package grpc

import (
	"example.com/memo/api/memo/v1"
	"example.com/memo/server/internal/app"
	"google.golang.org/grpc"
)

// RegisterMemoService は gRPC サーバーに MemoService を登録する関数。
// この関数は、リポジトリ、サービス、ハンドラーを初期化し、gRPC にサービスを登録する責任を持つ。
func RegisterMemoService(s *grpc.Server) {
	// MemoRepository を初期化（データアクセス層）
	repo := app.NewMemoRepository()

	// MemoService を初期化（ビジネスロジック層）
	svc := app.NewMemoService(repo)

	// gRPC ハンドラーを初期化（RPCリクエストをサービスに橋渡しする役割）
	handler := NewMemoHandler(svc)

	// gRPC サーバーに MemoService を登録
	memo.RegisterMemoServiceServer(s, handler)
}
