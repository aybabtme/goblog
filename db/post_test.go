package db

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	newPost(t, setupSQLitePersist())
	//newPost(t, setupPGPersist())
}

func newPost(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	var post = persist.NewPost("Antoine",
		"Hello World",
		"fake.url/to/image.jpg",
		time.Now().UTC())
	if post == nil {
		t.Error("Receive a nil post")
	}
}

func TestSavePost(t *testing.T) {
	savePost(t, setupSQLitePersist())
	//savePost(t, setupPGPersist())
}

func savePost(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	var post = persist.NewPost("Antoine",
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
	destroyPost(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reason
	// destroyPost(t, setupPGPersist())
}

func destroyPost(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	for i := int64(1); i < 100; i++ {
		var expected = pers.NewPost(
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := pers.FindPostById(id)

		if actual != nil {
			t.Error("Post shouldnt exist in DB after destroy")
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

func TestFindByIdPost(t *testing.T) {
	findByIdPost(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reasons
	//findByIdPost(t, setupPGPersist())
}

func findByIdPost(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	for i := int64(1); i < 100; i++ {
		var expected = persist.NewPost(
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		expected.Save()

		actual, err := persist.FindPostById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying post %d: %v", i, err)
		}

		if actual.Content() != expected.Content() {
			t.Errorf("Expected <%s> but was <%s>\n", expected.Content(), actual.Content())
		}
	}
}

func TestFindAllPost(t *testing.T) {
	findAllPost(t, setupSQLitePersist())
	//findAllPost(t, setupPGPersist())
}

func findAllPost(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	var postCount = int64(10)

	for i := int64(1); i <= postCount; i++ {
		var post = pers.NewPost(
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		post.Save()
	}

	posts, err := pers.FindAllPosts()
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
	idIncrements(t, setupSQLitePersist())
	// TODO PG doesnt work
	// idIncrements(t, setupPGPersist())
}

func idIncrements(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	for i := int64(1); i < 100; i++ {
		var post = persist.NewPost(
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

//
// Helpers
//

func setupSQLitePersist() *Persister {
	var pers, _ = NewPersistance(NewSQLiter("test"))
	return pers
}

func setupPGPersist() *Persister {
	var pers, _ = NewPersistance(NewPostgreser("test", "antoine"))
	return pers
}
