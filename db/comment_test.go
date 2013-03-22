package db

import (
	"fmt"
	"testing"
	"time"
)

func generateUserAndPost(pers *Persister, i int64) (*User, *Post) {
	var user = pers.NewUser("John", time.Now().UTC(), -5, "john@smith.com")
	_ = user.Save()
	var authUser = pers.NewUser("Antoine", time.Now().UTC(), -5, "antoine@grondin.com")
	var author = pers.NewAuthor(
		fmt.Sprintf("aybabtme#%d", i),
		authUser)
	_ = author.Save()
	var post = pers.NewPost(author.Id(),
		fmt.Sprintf("Title #%d", i),
		fmt.Sprintf("Content #%d", i),
		fmt.Sprintf("ImageUrl #%d", i),
		time.Now().UTC())
	return user, post
}

func TestNewComment(t *testing.T) {
	newComment(t, setupSQLitePersist())
	//newComment(t, setupPGPersist())
}

func newComment(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	var user, post = generateUserAndPost(pers, 0)
	var comment = pers.NewComment(
		user.Id(),
		post.Id(),
		"I really love your new pointless post on Justin Bieber",
		time.Now().UTC())

	if comment == nil {
		t.Error("Receive a nil comment")
	}
}

func TestSaveComment(t *testing.T) {
	saveComment(t, setupSQLitePersist())
	//saveComment(t, setupPGPersist())
}

func saveComment(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	var user, post = generateUserAndPost(pers, 0)
	var comment = pers.NewComment(
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
	destroyComment(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reason
	// destroyComment(t, setupPGPersist())
}

func destroyComment(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	for i := int64(1); i < 10; i++ {

		var user, post = generateUserAndPost(pers, 0)
		var expected = pers.NewComment(
			user.Id(),
			post.Id(),
			fmt.Sprintf("I really love your new pointless post #%d", i),
			time.Now().UTC())

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := pers.FindCommentById(id)

		if actual != nil {
			t.Error("Comment shouldnt exist in DB after destroy")
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

func TestFindByIdComment(t *testing.T) {
	findByIdComment(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reasons
	//findByIdComment(t, setupPGPersist())
}

func findByIdComment(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()
	for i := int64(1); i < 10; i++ {
		var user, post = generateUserAndPost(pers, 0)
		var expected = pers.NewComment(
			user.Id(),
			post.Id(),
			fmt.Sprintf("I really love your new pointless post #%d", i),
			time.Now().UTC())

		expected.Save()

		actual, err := pers.FindCommentById(expected.Id())

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
	findAllComment(t, setupSQLitePersist())
	//findAllComment(t, setupPGPersist())
}

func findAllComment(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()
	var commentCount = int64(10)

	for i := int64(1); i <= commentCount; i++ {
		var user, post = generateUserAndPost(pers, 0)
		var comment = pers.NewComment(
			user.Id(),
			post.Id(),
			fmt.Sprintf("I really love your new pointless post #%d", i),
			time.Now().UTC())

		comment.Save()
	}

	comments, err := pers.FindAllComments()
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
	commentIdIncrements(t, setupSQLitePersist())
	// TODO PG doesnt work
	// idIncrements(t, setupPGPersist())
}

func commentIdIncrements(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	for i := int64(1); i < 10; i++ {
		var user, post = generateUserAndPost(pers, 0)
		var comment = pers.NewComment(
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
