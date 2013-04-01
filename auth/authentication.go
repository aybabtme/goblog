package auth

import (
	"github.com/aybabtme/goblog/model"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
)

var store = sessions.NewCookieStore(
	[]byte("auth key avec de la mayo"),
	[]byte("crypto key avec de la mayo"),
	[]byte("auth key avec du ketchup"),
	[]byte("crypto key avec du ketchup"),
)

func Login(conn *model.DBConnection, w http.ResponseWriter, r *http.Request) (*model.User, *model.Author) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	user := getUser(conn, store, r)
	author := getAuthor(conn, store, r)

	// Save it.
	sessions.Save(r, w)

	log.Printf("Session:User=%v Author=%v\n", user, author)

	return user, author
}

func getUser(conn *model.DBConnection,
	store sessions.Store,
	r *http.Request) *model.User {

	session, err := store.Get(r, "user-session")
	if err == nil {
		return nil
	}
	idStr, ok := session.Values["userId"].(string)
	if !ok {
		return nil
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil
	}
	user, err := conn.FindUserById(id)
	if err != nil {
		return nil
	}
	return user
}

func getAuthor(conn *model.DBConnection,
	store sessions.Store,
	r *http.Request) *model.Author {

	session, err := store.Get(r, "author-session")
	if err == nil {
		return nil
	}
	idStr, ok := session.Values["authorId"].(string)
	if !ok {
		return nil
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil
	}
	author, err := conn.FindAuthorById(id)

	if err != nil {
		return nil
	}

	return author

}
