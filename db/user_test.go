package db

import (
	"fmt"
	"testing"
	"time"
)

func generateUser(conn *DBConnection, i int64) *User {
	var user = conn.NewUser(
		fmt.Sprintf("Antoine #%d", i),
		time.Now().UTC(),
		-5,
		"g+",
		fmt.Sprintf("anAuthToken#%d", i),
		fmt.Sprintf("a%d@b.com", i))
	return user
}

func TestNewUser(t *testing.T) {
	newUser(t, setupPGConnection())
}

func newUser(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	user := generateUser(conn, 0)
	if user == nil {
		t.Error("Receive a nil user")
	}
}

func TestSaveUser(t *testing.T) {
	saveUser(t, setupPGConnection())
}

func saveUser(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	user := generateUser(conn, 0)
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
	destroyUser(t, setupPGConnection())
}

func destroyUser(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	for i := int64(1); i < 10; i++ {
		expected := generateUser(conn, i)

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := conn.FindUserById(id)

		if actual != nil {
			t.Error("User shouldnt exist in DB after destroy")
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

func TestFindByIdUser(t *testing.T) {
	findByIdUser(t, setupPGConnection())
}

func findByIdUser(t *testing.T, conn *DBConnection) {

	for i := int64(1); i < 10; i++ {
		expected := generateUser(conn, i)

		expected.Save()

		actual, err := conn.FindUserById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying User %d, expected <%d> but was: %v", i, expected.Id(), err)
			return
		}

		if actual.Username() != expected.Username() {
			t.Errorf("Expected <%s> but was <%s>\n", expected.Username(), actual.Username())
		}
	}

	defer conn.DeleteConnection()
}

func TestFindAllUser(t *testing.T) {
	findAllUser(t, setupPGConnection())
}

func findAllUser(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	var userCount = int64(10)

	for i := int64(1); i <= userCount; i++ {
		user := generateUser(conn, i)
		user.Save()
	}

	users, err := conn.FindAllUsers()
	if err != nil {
		t.Errorf("Couldn't query Users although just saved %d", userCount)
		return
	}

	if users == nil {
		t.Errorf("Saved %d Users but query returns none", userCount)
		return
	}

	if int64(len(users)) != userCount {
		t.Errorf("Saved and expected <%d> Users, was <%d>",
			userCount, len(users))
		return
	}

	for idx, user := range users {
		if user.Id() != int64(idx)+int64(1) {
			t.Errorf("Expected <%d> but was <%d>", idx+1, user.Id())
		}
	}
}

func TestUserIdIncrements(t *testing.T) {
	userIdIncrements(t, setupPGConnection())
}

func userIdIncrements(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()

	for i := int64(1); i < 10; i++ {
		user := generateUser(conn, i)

		if user.Id() != -1 {
			t.Error("Id should be of -1 at this point")
		}

		user.Save()

		if user.Id() != i {
			t.Errorf("Id expected %d but was %d", i, user.Id())
		}
	}
}

func TestFindAllUserComments(t *testing.T) {
	findAllUserComments(t, setupPGConnection())
}

func findAllUserComments(t *testing.T, conn *DBConnection) {
	defer conn.DeleteConnection()
	var commentCount = 10

	var user, post = generateUserAndPost(conn, 0)
	var expected []Comment
	for i := 0; i < commentCount; i++ {
		var comment = conn.NewComment(user.Id(), post.Id(),
			fmt.Sprintf("I agree * %d", i),
			time.Now().UTC())
		comment.Save()
		expected = append(expected, *comment)
	}

	actual, err := user.Comments()
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
