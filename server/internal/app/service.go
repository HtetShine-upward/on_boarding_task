package app

import (
	"errors" // エラー定義のための標準パッケージ
	"time"   // 時刻取得のためのパッケージ

	memo "example.com/memo/api/memo/v1" // memoパッケージ（gRPCで定義されたメモ構造体）
	"github.com/google/uuid"            // UUID生成ライブラリ
)

// バリデーションや検索失敗時に使うエラー定数（Sentinel errors）
var (
	ErrInvalidTitle   = errors.New("タイトルは1〜100文字である必要があります")
	ErrInvalidContent = errors.New("コンテンツは0〜2000文字である必要があります")
	ErrMemoNotFound   = errors.New("該当するIDのメモが存在しません")
)

// MemoServiceインターフェース：メモ操作のビジネスロジックを定義
type MemoService interface {
	Create(title, content string) (*memo.Memo, error) // メモ作成
	Get(id string) (*memo.Memo, error)                // メモ取得
	List() []*memo.Memo                               // メモ一覧取得
	Delete(id string) error                           // メモ削除
}

// memoService構造体：MemoServiceインターフェースの実装。MemoRepositoryに依存。
type memoService struct {
	repo MemoRepository // データアクセス層（リポジトリ）への依存
}

// NewMemoService：MemoServiceのインスタンスを生成するコンストラクタ関数
func NewMemoService(repo MemoRepository) MemoService {
	return &memoService{repo: repo}
}

// Create：タイトルとコンテンツをバリデーションし、メモを作成・保存
func (s *memoService) Create(title, content string) (*memo.Memo, error) {
	if len(title) == 0 || len(title) > 100 {
		return nil, ErrInvalidTitle // タイトルが不正
	}
	if len(content) > 2000 {
		return nil, ErrInvalidContent // コンテンツが不正
	}

	// メモ構造体を生成（UUIDと作成時刻を付加）
	m := &memo.Memo{
		Id:        uuid.NewString(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UnixMilli(),
	}

	s.repo.Save(m) // リポジトリに保存
	return m, nil  // 作成したメモを返す
}

// Get：IDでメモを取得。存在しない場合はエラーを返す。
func (s *memoService) Get(id string) (*memo.Memo, error) {
	if len(id) == 0 {
		return nil, ErrInvalidTitle // 空のIDは不正
	}
	m, ok := s.repo.Get(id)
	if !ok {
		return nil, ErrMemoNotFound // メモが見つからない
	}
	return m, nil // メモを返す
}

// List：すべてのメモを取得
func (s *memoService) List() []*memo.Memo {
	return s.repo.List()
}

// Delete：IDでメモを削除。存在しない場合はエラーを返す。
func (s *memoService) Delete(id string) error {
	if len(id) == 0 {
		return ErrInvalidTitle // 空のIDは不正
	}
	_, ok := s.repo.Get(id)
	if !ok {
		return ErrMemoNotFound // メモが存在しない
	}
	s.repo.Delete(id) // メモを削除
	return nil
}
