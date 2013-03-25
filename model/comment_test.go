package model

import (
	"fmt"
	"testing"
	"time"
)

func generateUserAndPost(conn *DBConnection, i int64) (*User, *Post) {
	// i + 1 to cover when i = 0
	user := generateUser(conn, (i+1)*100000000)
	user.Save()
	authUser := generateUser(conn, i)
	var author = conn.NewAuthor(authUser)
	_ = author.Save()
	var post = conn.NewPost(author,
		fmt.Sprintf("Title #%d", i),
		fmt.Sprintf("Content #%d", i),
		fmt.Sprintf("ImageUrl #%d", i),
		time.Now().UTC())
	_ = post.Save()
	return user, post
}

func TestNewComment(t *testing.T) {
	newComment(t, setupPGConnection())
}

func newComment(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	var user, post = generateUserAndPost(conn, 0)
	var comment = conn.NewComment(
		user.Id(),
		post.Id(),
		"I really love your new pointless post on Justin Bieber",
		time.Now().UTC())

	if comment == nil {
		t.Error("Receive a nil comment")
	}
}

func TestSaveComment(t *testing.T) {
	saveComment(t, setupPGConnection())
}

func saveComment(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	var user, post = generateUserAndPost(conn, 0)
	var comment = conn.NewComment(
		user.Id(),
		post.Id(),
		"I really love your new pointless post on Justin Bieber",
		time.Now().UTC())

	if comment == nil {
		t.Error("Receive a nil comment")
	}

	if comment.Id() != -1 {
		t.Error("Id should be of -1 at this point")
	}

	if err := comment.Save(); err != nil {
		t.Error("Save failed", err)
	}

	if comment.Id() != 1 {
		t.Error("Id should be 1 at this point")
	}
}

func TestDestroyComment(t *testing.T) {
	destroyComment(t, setupPGConnection())
}

func destroyComment(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	for i := int64(1); i < int64(10); i++ {

		var user, post = generateUserAndPost(conn, i)
		var expected = conn.NewComment(
			user.Id(),
			post.Id(),
			fmt.Sprintf("I really love your new pointless post #%d", i),
			time.Now().UTC())

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := conn.FindCommentById(id)

		if actual != nil {
			t.Error("Comment shouldnt exist in DB after destroy")
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

func TestFindByIdComment(t *testing.T) {
	findByIdComment(t, setupPGConnection())
}

func findByIdComment(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	for i := int64(1); i < 10; i++ {
		var user, post = generateUserAndPost(conn, i)
		var expected = conn.NewComment(
			user.Id(),
			post.Id(),
			fmt.Sprintf("I really love your new pointless post #%d", i),
			time.Now().UTC())

		expected.Save()

		actual, err := conn.FindCommentById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying comment %d: %v", i, err)
			return
		}

		if actual.Content() != expected.Content() {
			t.Errorf("Expected <%s> but was <%s>\n",
				expected.Content(), actual.Content())
			return
		}
	}
}

func TestFindAllComment(t *testing.T) {
	findAllComment(t, setupPGConnection())
}

func findAllComment(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	var commentCount = int64(10)

	for i := int64(1); i <= commentCount; i++ {
		var user, post = generateUserAndPost(conn, i)
		var comment = conn.NewComment(
			user.Id(),
			post.Id(),
			fmt.Sprintf("I really love your new pointless post #%d", i),
			time.Now().UTC())

		comment.Save()
	}

	comments, err := conn.FindAllComments()
	if err != nil {
		t.Errorf("Couldn't query comments although just saved %d",
			commentCount)
	}

	if comments == nil {
		t.Errorf("Saved %d comments but query returns none",
			commentCount)
	}

	if int64(len(comments)) != commentCount {
		t.Errorf("Saved and expected <%d> posts, was <%d>",
			commentCount, len(comments))
	}

	for idx, comment := range comments {
		if comment.Id() != int64(idx)+int64(1) {
			t.Errorf("Expected <%d> but was <%d>", idx+1, comment.Id())
		}
	}
}

func TestCommentIdIncrements(t *testing.T) {
	commentIdIncrements(t, setupPGConnection())
}

func commentIdIncrements(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	for i := int64(1); i < int64(10); i++ {
		var user, post = generateUserAndPost(conn, i)
		var comment = conn.NewComment(
			user.Id(),
			post.Id(),
			fmt.Sprintf("I really love your new pointless post #%d", i),
			time.Now().UTC())

		if comment.Id() != -1 {
			t.Error("Id should be of -1 at this point")
		}

		comment.Save()

		if comment.Id() != i {
			t.Errorf("Id expected %d but was %d", i, comment.Id())
		}
	}
}
