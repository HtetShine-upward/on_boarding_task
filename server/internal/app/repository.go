package app

import (
	"sync" // 排他制御のための同期パッケージ

	memo "example.com/memo/api/memo/v1" // memoパッケージのインポート（Memo構造体を使用）
)

// MemoRepositoryインターフェース：メモの保存・取得・一覧・削除の操作を定義
type MemoRepository interface {
	Save(m *memo.Memo)                // メモを保存
	Get(id string) (*memo.Memo, bool) // IDでメモを取得（見つからない場合はfalse）
	List() []*memo.Memo               // 全メモを一覧取得
	Delete(id string)                 // メモを削除
}

// inMemoryMemoRepository構造体：メモをメモリ上に保存する実装
type inMemoryMemoRepository struct {
	mu    sync.RWMutex          // 読み書き用の排他制御（複数読み取り、単一書き込み）
	store map[string]*memo.Memo // メモをIDで管理するマップ
}

// NewMemoRepository関数：新しいメモリベースのリポジトリを作成して返す
func NewMemoRepository() MemoRepository {
	return &inMemoryMemoRepository{
		store: make(map[string]*memo.Memo), // storeを初期化
	}
}

// Saveメソッド：メモを保存（書き込みロックを使用）
func (r *inMemoryMemoRepository) Save(m *memo.Memo) {
	r.mu.Lock()         // 書き込みロックを取得
	defer r.mu.Unlock() // 関数終了時にロック解除
	r.store[m.Id] = m   // メモをIDで保存
}

// Getメソッド：IDでメモを取得（読み取りロックを使用）
func (r *inMemoryMemoRepository) Get(id string) (*memo.Memo, bool) {
	r.mu.RLock()         // 読み取りロックを取得
	defer r.mu.RUnlock() // 関数終了時にロック解除
	m, ok := r.store[id] // メモを取得
	return m, ok         // メモと存在フラグを返す
}

// Listメソッド：すべてのメモを一覧で取得（読み取りロックを使用）
func (r *inMemoryMemoRepository) List() []*memo.Memo {
	r.mu.RLock()         // 読み取りロックを取得
	defer r.mu.RUnlock() // 関数終了時にロック解除
	var memos []*memo.Memo
	for _, m := range r.store { // store内のすべてのメモをループ
		memos = append(memos, m) // スライスに追加
	}
	return memos // メモ一覧を返す
}

// Deleteメソッド：指定されたIDのメモを削除（書き込みロックを使用）
func (r *inMemoryMemoRepository) Delete(id string) {
	r.mu.Lock()         // 書き込みロックを取得
	defer r.mu.Unlock() // 関数終了時にロック解除
	delete(r.store, id) // メモを削除
}
