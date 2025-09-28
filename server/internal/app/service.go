package app

import (
	"errors"
	"time"

	memo "example.com/memo/api/memo/v1"
	"github.com/google/uuid"
)

// Sentinel errors for validation and lookup failures.
var (
	ErrInvalidTitle   = errors.New("タイトルは1〜100文字である必要があります")
	ErrInvalidContent = errors.New("コンテンツは0〜2000文字である必要があります")
	ErrMemoNotFound   = errors.New("該当するIDのメモが存在しません")
)

// MemoService defines the interface for memo operations.
type MemoService interface {
	Create(title, content string) (*memo.Memo, error)
	Get(id string) (*memo.Memo, error)
	List() []*memo.Memo
	Delete(id string) error
}

// memoService implements MemoService using a MemoRepository.
type memoService struct {
	repo MemoRepository
}

// NewMemoService creates a new MemoService.
func NewMemoService(repo MemoRepository) MemoService {
	return &memoService{repo: repo}
}

// Create validates input and stores a new memo.
func (s *memoService) Create(title, content string) (*memo.Memo, error) {
	if len(title) == 0 || len(title) > 100 {
		return nil, ErrInvalidTitle
	}
	if len(content) > 2000 {
		return nil, ErrInvalidContent
	}

	m := &memo.Memo{
		Id:        uuid.NewString(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UnixMilli(),
	}
	s.repo.Save(m)
	return m, nil
}

// Get retrieves a memo by ID or returns ErrMemoNotFound.
func (s *memoService) Get(id string) (*memo.Memo, error) {
	if len(id) == 0 {
		return nil, ErrInvalidTitle
	}
	m, ok := s.repo.Get(id)
	if !ok {
		return nil, ErrMemoNotFound
	}
	return m, nil
}

// List returns all memos.
func (s *memoService) List() []*memo.Memo {
	return s.repo.List()
}

// Delete removes a memo by ID or returns ErrMemoNotFound.
func (s *memoService) Delete(id string) error {
	if len(id) == 0 {
		return ErrInvalidTitle
	}
	_, ok := s.repo.Get(id)
	if !ok {
		return ErrMemoNotFound
	}
	s.repo.Delete(id)
	return nil
}
