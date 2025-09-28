package grpc

import (
	"example.com/memo/api/memo/v1"
	"example.com/memo/server/internal/app"
	"google.golang.org/grpc"
)

func RegisterMemoService(s *grpc.Server) {
	repo := app.NewMemoRepository()
	svc := app.NewMemoService(repo)
	handler := NewMemoHandler(svc)
	memo.RegisterMemoServiceServer(s, handler)
}
