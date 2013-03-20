package db

import (
	"fmt"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	newUser(t, setupSQLitePersist())
	//newUser(t, setupPGPersist())
}

func newUser(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	var user = persist.NewUser("Antoine",
		"Hello World",
		"fake.url/to/image.jpg",
		time.Now().UTC())
	if user == nil {
		t.Error("Receive a nil user")
	}
}

func TestSaveUser(t *testing.T) {
	saveUser(t, setupSQLitePersist())
	//saveUser(t, setupPGPersist())
}

func saveUser(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	var user = persist.NewUser("Antoine",
		"Hello World",
		"fake.url/to/image.jpg",
		time.Now().UTC())
	if user == nil {
		t.Error("Receive a nil user")
	}

	if user.Id() != -1 {
		t.Error("Id should be of -1 at this point")
	}

	if err := user.Save(); err != nil {
		t.Error("Save failed", err)
	}

	if user.Id() != 1 {
		t.Error("Id should be 1 at this point")
	}
}

func TestDestroyUser(t *testing.T) {
	destroyUser(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reason
	// destroyUser(t, setupPGPersist())
}

func destroyUser(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	for i := int64(1); i < 100; i++ {
		var expected = pers.NewUser(
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := pers.FindUserById(id)

		if actual != nil {
			t.Error("User shouldnt exist in DB after destroy")
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

func TestFindByIdUser(t *testing.T) {
	findByIdUser(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reasons
	//findByIdUser(t, setupPGPersist())
}

func findByIdUser(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	for i := int64(1); i < 100; i++ {
		var expected = persist.NewUser(
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		expected.Save()

		actual, err := persist.FindUserById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying User %d: %v", i, err)
		}

		if actual.Content() != expected.Content() {
			t.Errorf("Expected <%s> but was <%s>\n", expected.Content(), actual.Content())
		}
	}
}

func TestFindAllUser(t *testing.T) {
	findAllUser(t, setupSQLitePersist())
	//findAllUser(t, setupPGPersist())
}

func findAllUser(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	var userCount = int64(10)

	for i := int64(1); i <= userCount; i++ {
		var user = pers.NewUser(
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		user.Save()
	}

	users, err := pers.FindAllUsers()
	if err != nil {
		t.Errorf("Couldn't query Users although just saved %d", userCount)
	}

	if users == nil {
		t.Errorf("Saved %d Users but query returns none", userCount)
	}

	if int64(len(users)) != userCount {
		t.Errorf("Saved and expected <%d> Users, was <%d>",
			userCount, len(users))
	}

	for idx, user := range users {
		if user.Id() != int64(idx)+int64(1) {
			t.Errorf("Expected <%d> but was <%d>", idx+1, user.Id())
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
		var user = persist.NewUser(
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())

		if user.Id() != -1 {
			t.Error("Id should be of -1 at this point")
		}

		user.Save()

		if user.Id() != i {
			t.Errorf("Id expected %d but was %d", i, user.Id())
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
	var pers, _ = NewPersistance(NewUsergreser("test", "antoine"))
	return pers
}
