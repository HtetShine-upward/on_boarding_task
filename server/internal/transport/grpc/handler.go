package grpc

import (
	"context" // コンテキストによるタイムアウトやキャンセル管理

	"example.com/memo/api/memo/v1"
	"example.com/memo/server/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MemoHandler は gRPC の MemoServiceServer を実装する構造体
type MemoHandler struct {
	memo.UnimplementedMemoServiceServer                 // gRPCのインターフェースを継承（必須）
	svc                                 app.MemoService // ビジネスロジック層への依存（DI）
}

// NewMemoHandler は MemoHandler を初期化するコンストラクタ関数
func NewMemoHandler(svc app.MemoService) *MemoHandler {
	return &MemoHandler{svc: svc} // 依存注入されたサービスを保持
}

// CreateMemo はメモ作成リクエストを処理する
func (h *MemoHandler) CreateMemo(ctx context.Context, req *memo.CreateMemoRequest) (*memo.CreateMemoResponse, error) {
	// リクエストがタイムアウトしていないか確認
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	default:
	}

	// サービス層にメモ作成を依頼
	m, err := h.svc.Create(req.Title, req.Content)
	if err != nil {
		// バリデーションエラーの場合は InvalidArgument を返す
		switch err {
		case app.ErrInvalidTitle, app.ErrInvalidContent:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			// その他のエラーは Internal エラーとして返す
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// 成功した場合は作成されたメモをレスポンスとして返す
	return &memo.CreateMemoResponse{Memo: m}, nil
}

// GetMemo は指定されたIDのメモを取得する
func (h *MemoHandler) GetMemo(ctx context.Context, req *memo.GetMemoRequest) (*memo.GetMemoResponse, error) {
	// リクエストのタイムアウト確認
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	default:
	}

	// サービス層からメモを取得
	m, err := h.svc.Get(req.Id)
	if err != nil {
		// メモが存在しない場合は NotFound を返す
		if err == app.ErrMemoNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		// その他のエラーは Internal エラーとして返す
		return nil, status.Error(codes.Internal, err.Error())
	}

	// メモをレスポンスとして返す
	return &memo.GetMemoResponse{Memo: m}, nil
}

// ListMemos はすべてのメモを一覧で取得する
func (h *MemoHandler) ListMemos(ctx context.Context, req *memo.ListMemosRequest) (*memo.ListMemosResponse, error) {
	// リクエストのタイムアウト確認
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	default:
	}

	// サービス層からメモ一覧を取得
	memos := h.svc.List()
	return &memo.ListMemosResponse{Memos: memos}, nil
}

// DeleteMemo は指定されたIDのメモを削除する
func (h *MemoHandler) DeleteMemo(ctx context.Context, req *memo.DeleteMemoRequest) (*memo.DeleteMemoResponse, error) {
	// リクエストのタイムアウト確認
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	default:
	}

	// サービス層に削除処理を依頼
	err := h.svc.Delete(req.Id)
	if err != nil {
		// メモが存在しない場合は NotFound を返す
		if err == app.ErrMemoNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		// その他のエラーは Internal エラーとして返す
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 削除成功時は空のレスポンスを返す
	return &memo.DeleteMemoResponse{}, nil
}
