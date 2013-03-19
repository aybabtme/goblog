package db

import (
	"testing"
)

func TestNewPost(t *testing.T) {
	post := NewPost("Antoine", "I really love hamburgers.")

	if post == nil {
		t.Errorf("NewPost=nil, want non-nil")
	}
}
