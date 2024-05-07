package memdb

import (
	"testing"
)

func TestStore_Posts(t *testing.T) {
	store, _ := New(10)

	resp, err := store.Posts("3", 1)

	if err != nil {
		t.Error("erorr in memdb")
	}

	if len(resp.News) != 1 {
		t.Error("wrong get search posts")
	}
}

func TestStore_PostsPages(t *testing.T) {
	store, _ := New(10)

	resp, err := store.Posts("Title", 1)

	if err != nil {
		t.Error("erorr in memdb")
	}

	if len(resp.News) != 5 {
		t.Error("wrong get search posts")
	}

	if resp.Pages.Total != 2 {
		t.Error("wrong get search posts")
	}

	if resp.Pages.Current != 1 {
		t.Error("wrong get search posts")
	}
}
