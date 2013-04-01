package auth

import (
	"github.com/aybabtme/goblog/model"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
)

var store = sessions.NewCookieStore([]byte("auth key avec de la mayo"))

func Login(conn *model.DBConnection, w http.ResponseWriter, r *http.Request) (*model.User, *model.Author) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	user := getUser(conn, store, r)
	author := getAuthor(conn, store, r)

	// Save it.
	sessions.Save(r, w)

	return user, author
}

func getUser(conn *model.DBConnection,
	store sessions.Store,
	r *http.Request) *model.User {

	session, _ := store.Get(r, "user-session")
	idStr, ok := session.Values["userId"].(string)
	if !ok {
		log.Printf("auth.getUser. Couldn't get id string :%v", idStr)
		return nil
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("auth.getUser. Couldn't parse int64 from id string")
		return nil
	}
	user, err := conn.FindUserById(id)
	if err != nil {
		log.Printf("auth.getUser. Couldn't find user with id <%d>", id)
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
