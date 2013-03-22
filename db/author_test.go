package db

import (
	"fmt"
	"testing"
	"time"
)

func TestNewAuthor(t *testing.T) {
	newAuthor(t, setupSQLitePersist())
	//newAuthor(t, setupPGPersist())
}

func newAuthor(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	var user = persist.NewUser("Author Antoine",
		time.Now().UTC(), -5, "a@g.com")

	var author = persist.NewAuthor(
		"aybabtme",
		user)

	if author == nil {
		t.Error("Receive a nil author")
	}
}

func TestSaveAuthor(t *testing.T) {
	saveAuthor(t, setupSQLitePersist())
	//saveAuthor(t, setupPGPersist())
}

func saveAuthor(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	var user = persist.NewUser("Author Antoine",
		time.Now().UTC(), -5, "a@g.com")

	var author = persist.NewAuthor(
		"aybabtme",
		user)

	if author == nil {
		t.Error("Receive a nil author")
	}

	if author.Id() != -1 {
		t.Error("Id should be of -1 at this point")
	}

	if err := author.Save(); err != nil {
		t.Error("Save failed", err)
	}

	if author.Id() != 1 {
		t.Error("Id should be 1 at this point")
	}
}

func TestDestroyAuthor(t *testing.T) {
	destroyAuthor(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reason
	// destroyAuthor(t, setupPGPersist())
}

func destroyAuthor(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	for i := int64(1); i < 10; i++ {
		var user = pers.NewUser(
			fmt.Sprintf("Author Antoine #%d", i),
			time.Now().UTC(),
			-5,
			fmt.Sprintf("a%d@b.com", i))

		var expected = pers.NewAuthor(
			fmt.Sprintf("aybabtme #%d", i),
			user)

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := pers.FindAuthorById(id)

		if actual != nil {
			t.Error("Author shouldnt exist in DB after destroy")
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

func TestFindByIdAuthor(t *testing.T) {
	findByIdAuthor(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reasons
	//findByIdAuthor(t, setupPGPersist())
}

func findByIdAuthor(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	for i := int64(1); i < 10; i++ {
		var user = persist.NewUser(
			fmt.Sprintf("Author Antoine #%d", i),
			time.Now().UTC(),
			-5,
			fmt.Sprintf("a%d@b.com", i))

		var expected = persist.NewAuthor(
			fmt.Sprintf("aybabtme #%d", i),
			user)
		expected.Save()

		actual, err := persist.FindAuthorById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying author %d: %v", i, err)
			return
		}

		if actual.Twitter() != expected.Twitter() {
			t.Errorf("Expected <%s> but was <%s>\n",
				expected.Twitter(), actual.Twitter())
			return
		}
	}
}

func TestFindAllAuthor(t *testing.T) {
	findAllAuthor(t, setupSQLitePersist())
	//findAllAuthor(t, setupPGPersist())
}

func findAllAuthor(t *testing.T, pers *Persister) {

	var authorCount = int64(10)

	for i := int64(1); i <= authorCount; i++ {
		var user = pers.NewUser(
			fmt.Sprintf("Author Antoine #%d", i),
			time.Now().UTC(),
			-5,
			fmt.Sprintf("a%d@b.com", i))

		var author = pers.NewAuthor(
			fmt.Sprintf("aybabtme #%d", i),
			user)
		author.Save()
	}

	authors, err := pers.FindAllAuthors()
	if err != nil {
		t.Errorf("Couldn't query authors although just saved %d",
			authorCount)
	}

	if authors == nil {
		t.Errorf("Saved %d authors but query returns none",
			authorCount)
	}

	if int64(len(authors)) != authorCount {
		t.Errorf("Saved and expected <%d> posts, was <%d>",
			authorCount, len(authors))
	}

	for idx, author := range authors {
		if author.Id() != int64(idx)+int64(1) {
			t.Errorf("Expected <%d> but was <%d>", idx+1, author.Id())
			fmt.Printf("UserID(%d) Twitter(%s)\n",
				author.UserId(), author.Twitter())
		}
	}

	defer pers.DeletePersistance()
}

func TestAuthorIdIncrements(t *testing.T) {
	authorIdIncrements(t, setupSQLitePersist())
	// TODO PG doesnt work
	// idIncrements(t, setupPGPersist())
}

func authorIdIncrements(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	for i := int64(1); i < 10; i++ {
		var user = persist.NewUser(
			fmt.Sprintf("Antoine #%d", i),
			time.Now().UTC(),
			-5,
			fmt.Sprintf("a%d@b.com", i))

		var author = persist.NewAuthor(
			fmt.Sprintf("aybabtme #%d", i),
			user)

		if author.Id() != -1 {
			t.Error("Id should be of -1 at this point")
		}

		author.Save()

		if author.Id() != i {
			t.Errorf("Id expected %d but was %d", i, author.Id())
		}
	}
}

func TestDeleteUserCascadesToAuthor(t *testing.T) {
	deleteUserCascadesToAuthor(t, setupSQLitePersist())
	//deleteUserCascadesToAuthor(t, setupPGPersist())
}

func deleteUserCascadesToAuthor(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	var user = persist.NewUser("Antony", time.Now().UTC(), -5, "antoine@g.com")
	var author = persist.NewAuthor("aybabtme", user)
	_ = author.Save()

	var copyUser, _ = persist.FindUserById(author.UserId())
	var copyAuthor, _ = persist.FindAuthorById(author.Id())

	// Tested elsewhere, kind of redundant
	if author.Twitter() != copyAuthor.Twitter() {
		t.Errorf("Expected <%s> but was <%s>", author.Twitter(), copyAuthor.Twitter())
		return
	}

	copyUser.Destroy()

	var _, err = persist.FindAuthorById(author.Id())
	if err == nil {
		t.Error("User was deleted but Author could still be found.")
		return
	}

}

func TestFindAllAuthorPosts(t *testing.T) {
	findAllAuthorPosts(t, setupSQLitePersist())
	//findAllAuthorPosts(t, setupPGPersist())
}

func findAllAuthorPosts(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	var postCount = 10

	var user = persist.NewUser("Antony", time.Now().UTC(), -5, "antoine@g.com")
	var author = persist.NewAuthor("aybabtme", user)
	_ = author.Save()

	var ghostPosts, err = author.Posts()
	if len(ghostPosts) != 0 {
		t.Errorf("Author has no posts yet but query returned posts with len(%d).",
			len(ghostPosts))
		return
	}

	var expected []Post
	for i := 0; i < postCount; i++ {
		var post = persist.NewPost(author.Id(),
			fmt.Sprintf("Great Topic #%d", i),
			fmt.Sprintf("cool content #%d", i),
			fmt.Sprint("awesome@email%d.com", i),
			time.Now().UTC())
		post.Save()
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
