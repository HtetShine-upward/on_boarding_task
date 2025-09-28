package app_test

import (
	"fmt"
	"testing"

	"example.com/memo/server/internal/app"
)

func TestMemoService(t *testing.T) {

	repo := app.NewMemoRepository()
	service := app.NewMemoService(repo)

	// Test Create
	memo1, err := service.Create("Test Title", "Test Content")

	if err != nil {
		t.Fatalf("Failed to create memo: %v", err)
	}

	if memo1.Title != "Test Title" || memo1.Content != "Test Content" {
		t.Errorf("Unexpected memo created: %+v", memo1)
	}
	fmt.Printf("Created memo: %+v\n", memo1)

	// Test Get
	fetchedMemo, err := service.Get(memo1.Id)
	if err != nil {
		t.Fatalf("Failed to get memo: %v", err)
	}
	if fetchedMemo.Id != memo1.Id {
		t.Errorf("Fetched memo ID mismatch: got %s, want %s", fetchedMemo.Id, memo1.Id)
	}
	fmt.Printf("Fetched memo: %+v\n", fetchedMemo)

	// Test List
	memos := service.List()
	if len(memos) != 1 {
		t.Errorf("Expected 1 memo, got %d", len(memos))
	}
	fmt.Printf("List of memos: %+v\n", memos)

	//create another memo
	memo2, err := service.Create("Another Title", "Another Content")
	if err != nil {
		t.Fatalf("Failed to create second memo: %v", err)
	}
	fmt.Printf("Created second memo: %+v\n", memo2)
	memos = service.List()
	if len(memos) != 2 {
		t.Errorf("Expected 2 memos, got %d", len(memos))
	}
	fmt.Printf("List of memos: %+v\n", memos)
	// Test Delete
	err = service.Delete(memo1.Id)
	if err != nil {
		t.Fatalf("Failed to delete memo: %v", err)
	}
	fmt.Printf("Deleted memo with ID: %s\n", memo1.Id)

	// Verify deletion
	_, err = service.Get(memo1.Id)
	if err != app.ErrMemoNotFound {
		t.Errorf("Expected ErrMemoNotFound after deletion, got %v", err)
	}
	// Test validation
	_, err = service.Create("", "No Title")
	if err != app.ErrInvalidTitle {
		t.Errorf("Expected ErrInvalidTitle for empty title, got %v", err)
	}
	_, err = service.Create("A", string(make([]byte, 2001)))
	if err != app.ErrInvalidContent {
		t.Errorf("Expected ErrInvalidContent for too long content, got %v", err)
	}
	_, err = service.Get("")
	if err != app.ErrInvalidTitle {
		t.Errorf("Expected ErrInvalidTitle for empty ID, got %v", err)
	}
	_, err = service.Get("invalid-id")
	if err != app.ErrMemoNotFound {
		t.Errorf("Expected ErrMemoNotFound for non-existent ID, got %v", err)
	}

}
