package model

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	newPost(t, setupPGConnection())
}

func newPost(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	var author = generateAuthor(conn, 0)
	var post = conn.NewPost(author.Id(),
		"My first post",
		"Hello World",
		"fake.url/to/image.jpg",
		time.Now().UTC())
	if post == nil {
		t.Error("Receive a nil post")
	}
}

func TestSavePost(t *testing.T) {
	savePost(t, setupPGConnection())
}

func savePost(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	var author = generateAuthor(conn, 0)
	var post = conn.NewPost(author.Id(),
		"My first post",
		"Hello World",
		"fake.url/to/image.jpg",
		time.Now().UTC())
	if post == nil {
		t.Error("Receive a nil post")
	}

	if post.Id() != -1 {
		t.Error("Id should be of -1 at this point")
	}

	if err := post.Save(); err != nil {
		t.Error("Save failed", err)
	}

	if post.Id() != 1 {
		t.Error("Id should be 1 at this point")
	}
}

func TestDestroyPost(t *testing.T) {
	destroyPost(t, setupPGConnection())
}

func destroyPost(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	for i := int64(1); i < 10; i++ {
		var author = generateAuthor(conn, i)
		var expected = conn.NewPost(author.Id(),
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := conn.FindPostById(id)

		if actual != nil {
			t.Error("Post shouldnt exist in DB after destroy")
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

func TestFindByIdPost(t *testing.T) {
	findByIdPost(t, setupPGConnection())
}

func findByIdPost(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	for i := int64(1); i < 10; i++ {
		var author = generateAuthor(conn, i)
		var expected = conn.NewPost(author.Id(),
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		expected.Save()

		actual, err := conn.FindPostById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying post %d: %v", i, err)
		}

		if actual.Content() != expected.Content() {
			t.Errorf("Expected <%s> but was <%s>\n", expected.Content(), actual.Content())
		}
	}
}

func TestFindAllPost(t *testing.T) {
	findAllPost(t, setupPGConnection())
}

func findAllPost(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	var postCount = int64(10)

	for i := int64(1); i <= postCount; i++ {
		var author = generateAuthor(conn, i)
		var post = conn.NewPost(author.Id(),
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		post.Save()
	}

	posts, err := conn.FindAllPosts()
	if err != nil {
		t.Errorf("Couldn't query posts although just saved %d", postCount)
	}

	if posts == nil {
		t.Errorf("Saved %d posts but query returns none", postCount)
	}

	if int64(len(posts)) != postCount {
		t.Errorf("Saved and expected <%d> posts, was <%d>",
			postCount, len(posts))
	}

	for idx, post := range posts {
		if post.Id() != int64(idx)+int64(1) {
			t.Errorf("Expected <%d> but was <%d>", idx+1, post.Id())
		}
	}
}

func TestIdIncrements(t *testing.T) {
	postIdIncrements(t, setupPGConnection())
}

func postIdIncrements(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	for i := int64(1); i < 10; i++ {
		var author = generateAuthor(conn, i)
		var post = conn.NewPost(author.Id(),
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())

		if post.Id() != -1 {
			t.Error("Id should be of -1 at this point")
		}

		post.Save()

		if post.Id() != i {
			t.Errorf("Id expected %d but was %d", i, post.Id())
		}
	}
}

func TestFindAllPostComments(t *testing.T) {
	findAllPostComments(t, setupPGConnection())
}

func findAllPostComments(t *testing.T, connist *DBConnection) {
	defer connist.DeleteConnection()
	var commentCount = 10

	var user, post = generateUserAndPost(connist, 0)
	var expected []Comment
	for i := 0; i < commentCount; i++ {
		var comment = connist.NewComment(user.Id(), post.Id(),
			fmt.Sprintf("I agree times %d thousand", i),
			time.Now().UTC())
		comment.Save()
		expected = append(expected, *comment)
	}

	actual, err := post.Comments()
	if err != nil {
		t.Errorf("Post had %d comments but an error was returned when queried.",
			commentCount)
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
