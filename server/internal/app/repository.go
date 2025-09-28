package app

import (
	"sync"

	memo "example.com/memo/api/memo/v1"
)

type MemoRepository interface {
	Save(m *memo.Memo)
	Get(id string) (*memo.Memo, bool)
	List() []*memo.Memo
	Delete(id string)
}

type inMemoryMemoRepository struct {
	mu    sync.RWMutex
	store map[string]*memo.Memo
}

func NewMemoRepository() MemoRepository {
	return &inMemoryMemoRepository{
		store: make(map[string]*memo.Memo),
	}
}

func (r *inMemoryMemoRepository) Save(m *memo.Memo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[m.Id] = m
}

func (r *inMemoryMemoRepository) Get(id string) (*memo.Memo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.store[id]
	return m, ok
}

func (r *inMemoryMemoRepository) List() []*memo.Memo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var memos []*memo.Memo
	for _, m := range r.store {
		memos = append(memos, m)
	}
	return memos
}
func (r *inMemoryMemoRepository) Delete(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, id)
}
