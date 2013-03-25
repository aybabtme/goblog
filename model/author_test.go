package model

import (
	"fmt"
	"testing"
	"time"
)

func generateAuthor(conn *DBConnection, i int64) *Author {
	user := generateUser(conn, i)
	author := conn.NewAuthor(user)
	_ = author.Save()
	return author
}

func TestNewAuthor(t *testing.T) {
	newAuthor(t, setupPGConnection())
}

func newAuthor(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	user := generateUser(conn, 0)

	var author = conn.NewAuthor(
		user)

	if author == nil {
		t.Error("Receive a nil author")
		return
	}
}

func TestSaveAuthor(t *testing.T) {
	saveAuthor(t, setupPGConnection())
}

func saveAuthor(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	user := generateUser(conn, 0)

	var author = conn.NewAuthor(
		user)

	if author == nil {
		t.Error("Receive a nil author")
		return
	}

	if author.Id() != -1 {
		t.Error("Id should be of -1 at this point")
		return
	}

	if err := author.Save(); err != nil {
		t.Error("Save failed", err)
		return
	}

	if author.Id() != 1 {
		t.Error("Id should be 1 at this point")
		return
	}
}

func TestDestroyAuthor(t *testing.T) {
	destroyAuthor(t, setupPGConnection())
}

func destroyAuthor(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	for i := int64(1); i < 10; i++ {
		user := generateUser(conn, i)

		var expected = conn.NewAuthor(
			user)

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := conn.FindAuthorById(id)

		if actual != nil {
			t.Error("Author shouldnt exist in DB after destroy")
			return
		}

		if err == nil {
			t.Error("An error should have been raised")
			return
		}

	}
}

func TestFindByIdAuthor(t *testing.T) {
	findByIdAuthor(t, setupPGConnection())
}

func findByIdAuthor(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	for i := int64(1); i < 10; i++ {
		user := generateUser(conn, i)

		var expected = conn.NewAuthor(
			user)
		expected.Save()

		actual, err := conn.FindAuthorById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying author %d: %v", i, err)
			return
		}

		if actual.Id() != expected.Id() {
			t.Errorf("Expected <%d> but was <%d>\n",
				expected.Id(), actual.Id())
			return
		}
	}
}

func TestFindAllAuthor(t *testing.T) {
	findAllAuthor(t, setupPGConnection())
}

func findAllAuthor(t *testing.T, conn *DBConnection) {

	defer conn.DeleteConnection()
	var authorCount = int64(10)

	for i := int64(1); i <= authorCount; i++ {
		user := generateUser(conn, i)

		var author = conn.NewAuthor(user)
		author.Save()
	}

	authors, err := conn.FindAllAuthors()
	if err != nil {
		t.Errorf("Couldn't query authors although just saved %d",
			authorCount)
		return
	}

	if authors == nil {
		t.Errorf("Saved %d authors but query returns none",
			authorCount)
		return
	}

	if int64(len(authors)) != authorCount {
		t.Errorf("Saved and expected <%d> posts, was <%d>",
			authorCount, len(authors))
		return
	}

	for idx, author := range authors {
		if author.Id() != int64(idx)+int64(1) {
			t.Errorf("Expected <%d> but was <%d>", idx+1, author.Id())

			return
		}
	}

}

func TestAuthorIdIncrements(t *testing.T) {
	authorIdIncrements(t, setupPGConnection())
}

func authorIdIncrements(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	for i := int64(1); i < 10; i++ {
		user := generateUser(conn, i)

		var author = conn.NewAuthor(user)

		if author.Id() != -1 {
			t.Errorf("Id should be of -1 at this point, but is <%d>",
				author.Id())
			return
		}

		author.Save()

		if author.Id() != i {
			t.Errorf("Id expected %d but was %d", i, author.Id())
			return
		}
	}
}

func TestDeleteUserCascadesToAuthor(t *testing.T) {
	deleteUserCascadesToAuthor(t, setupPGConnection())
}

func deleteUserCascadesToAuthor(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	user := generateUser(conn, 0)
	var author = conn.NewAuthor(user)
	_ = author.Save()

	var copyUser, _ = conn.FindUserById(author.User().Id())

	if copyUser == nil {
		t.Error("Copyuser is nil")
		return
	}

	copyUser.Destroy()

	var _, err = conn.FindAuthorById(author.Id())
	if err == nil {
		t.Error("User was deleted but Author could still be found.")
		return
	}

}

func TestFindAllAuthorPosts(t *testing.T) {
	findAllAuthorPosts(t, setupPGConnection())
}

func findAllAuthorPosts(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	var postCount = int64(10)

	user := generateUser(conn, 0)
	var author = conn.NewAuthor(user)
	_ = author.Save()

	var ghostPosts, err = author.Posts()
	if len(ghostPosts) != 0 {
		t.Errorf("Author has no posts yet but query returned posts with len(%d).",
			len(ghostPosts))
		return
	}

	var expected []Post
	for i := int64(0); i < postCount; i++ {
		var post = conn.NewPost(author.Id(),
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())

		_ = post.Save()
		expected = append(expected, *post)
	}

	actual, err := author.Posts()
	if err != nil {
		t.Errorf("Author has %d posts but an error was returned when queried.",
			postCount)
		return
	}

	if len(actual) != len(expected) {
		t.Errorf("Expected <len(%d)> but was <len(%d)>", len(expected), len(actual))
		return
	}

	for i := 0; i < len(expected); i++ {
		if expected[i].Content() != actual[i].Content() {
			t.Errorf("Compare #%d, expected <%s> but was <%s>",
				i, expected[i].Content(), actual[i].Content())
			return
		}
	}
}
