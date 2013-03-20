package db

import (
	"fmt"
	"testing"
	"time"
)

func TestSQLiteNewPost(t *testing.T) {
	var sqlite = setupSQLitePersist()
	newPost(t, sqlite)

	var pg = setupPGPersist()
	newPost(t, pg)
}

func newPost(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	var post = persist.NewPost("Antoine", "Hello World", time.Now().UTC())
	if post == nil {
		t.Error("Receive a nil post")
	}
}

func TestSavePost(t *testing.T) {
	var pg = setupSQLitePersist()
	savePost(t, pg)

	var sqlite = setupSQLitePersist()
	savePost(t, sqlite)
}

func savePost(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	var post = persist.NewPost("Antoine", "Hello World", time.Now().UTC())
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

func TestIdIncrements(t *testing.T) {
	var sqlite = setupSQLitePersist()
	idIncrements(t, sqlite)

	var pg = setupSQLitePersist()
	idIncrements(t, pg)

}

func idIncrements(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	for i := int64(1); i < 100; i++ {
		var post = persist.NewPost(
			fmt.Sprintf("Author #%d", i),
			fmt.Sprintf("Content #%d", i),
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
