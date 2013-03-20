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

	var author = persist.NewAuthor(
	// TODO
	)
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
	var author = persist.NewAuthor(
	// TODO
	)
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

	for i := int64(1); i < 100; i++ {
		var expected = pers.NewAuthor(
		// TODO
		)

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
	for i := int64(1); i < 100; i++ {
		var expected = persist.NewAuthor(
		// TODO
		)
		expected.Save()

		actual, err := persist.FindAuthorById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying author %d: %v", i, err)
		}

		if actual.Content() != expected.Content() {
			t.Errorf("Expected <%s> but was <%s>\n", expected.Content(), actual.Content())
		}
	}
}

func TestFindAllAuthor(t *testing.T) {
	findAllAuthor(t, setupSQLitePersist())
	//findAllAuthor(t, setupPGPersist())
}

func findAllAuthor(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	var authorCount = int64(10)

	for i := int64(1); i <= authorCount; i++ {
		var author = pers.NewAuthor(
		// TODO
		)
		author.Save()
	}

	authors, err := pers.FindAllAuthors()
	if err != nil {
		t.Errorf("Couldn't query authors although just saved %d", authorCount)
	}

	if authors == nil {
		t.Errorf("Saved %d authors but query returns none", authorCount)
	}

	if int64(len(authors)) != authorCount {
		t.Errorf("Saved and expected <%d> posts, was <%d>",
			authorCount, len(authors))
	}

	for idx, author := range authors {
		if author.Id() != int64(idx)+int64(1) {
			t.Errorf("Expected <%d> but was <%d>", idx+1, author.Id())
		}
	}
}

func TestAuthorIdIncrements(t *testing.T) {
	authorIdIncrements(t, setupSQLitePersist())
	// TODO PG doesnt work
	// idIncrements(t, setupPGPersist())
}

func authorIdIncrements(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	for i := int64(1); i < 100; i++ {
		var author = persist.NewAuthor(
		// TODO
		)

		if author.Id() != -1 {
			t.Error("Id should be of -1 at this point")
		}

		author.Save()

		if author.Id() != i {
			t.Errorf("Id expected %d but was %d", i, author.Id())
		}
	}
}
