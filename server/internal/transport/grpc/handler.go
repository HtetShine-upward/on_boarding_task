package grpc

import (
	"context"

	"example.com/memo/api/memo/v1"
	"example.com/memo/server/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MemoHandler struct {
	memo.UnimplementedMemoServiceServer
	svc app.MemoService
}

func NewMemoHandler(svc app.MemoService) *MemoHandler {
	return &MemoHandler{svc: svc}
}

func (h *MemoHandler) CreateMemo(ctx context.Context, req *memo.CreateMemoRequest) (*memo.CreateMemoResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	default:
	}

	m, err := h.svc.Create(req.Title, req.Content)
	if err != nil {
		switch err {
		case app.ErrInvalidTitle, app.ErrInvalidContent:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &memo.CreateMemoResponse{Memo: m}, nil
}

func (h *MemoHandler) GetMemo(ctx context.Context, req *memo.GetMemoRequest) (*memo.GetMemoResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	default:
	}

	m, err := h.svc.Get(req.Id)
	if err != nil {
		if err == app.ErrMemoNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &memo.GetMemoResponse{Memo: m}, nil
}

func (h *MemoHandler) ListMemos(ctx context.Context, req *memo.ListMemosRequest) (*memo.ListMemosResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	default:
	}

	memos := h.svc.List()
	return &memo.ListMemosResponse{Memos: memos}, nil
}

func (h *MemoHandler) DeleteMemo(ctx context.Context, req *memo.DeleteMemoRequest) (*memo.DeleteMemoResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	default:
	}

	err := h.svc.Delete(req.Id)
	if err != nil {
		if err == app.ErrMemoNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &memo.DeleteMemoResponse{}, nil
}
