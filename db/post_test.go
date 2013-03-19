package db

import (
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	post := NewPost("Antoine", "I really love hamburgers.", time.Now().UTC())

	if post == nil {
		t.Errorf("NewPost=nil, want non-nil")
	}
}

func TestPostId(t *testing.T) {
	post := NewPost("Jack", "This is great", time.Now().UTC())

	if post.Id() == nil {
		t.Errorf("post.Id() = nil, want non-nil")
	}

	// TODO, stub the DB to
}

func TestAuthor(t *testing.T) {

}
