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

	if user != nil {
		if author != nil {
			log.Printf("LOGIN: Author id(%d)<%v>",
				author.User().Id(),
				author.User().Username())
		} else {
			log.Printf("LOGIN: User id(%d)<%v>",
				user.Id(),
				user.Username())
		}
	}

	// Save it.
	sessions.Save(r, w)

	return user, author
}

func Logout(conn *model.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-session")
		session.Values["userId"] = ""
		session.Values["authorId"] = ""
		sessions.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func getUser(conn *model.DBConnection,
	store sessions.Store,
	r *http.Request) *model.User {

	session, _ := store.Get(r, "user-session")
	idStr, ok := session.Values["userId"].(string)
	if !ok || idStr == "" {
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

	session, err := store.Get(r, "user-session")
	idStr, ok := session.Values["authorId"].(string)
	if !ok || idStr == "" {
		return nil
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		log.Printf("auth.getAuthor. Couldn't parse int64 from id string: %s\n",
			idStr)
		return nil
	}
	author, err := conn.FindAuthorById(id)

	if err != nil {
		log.Printf("auth.getAuthor. Couldn't find user with id <%d>", id)
		return nil
	}

	return author

}
