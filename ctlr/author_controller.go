package ctlr

import (
	"github.com/aybabtme/goblog/auth"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func NewAuthorController() Controller {
	var a author
	a.path = "/author/{id:[0-9]+}"
	a.view = view.GetAuthorTemplate()
	return a
}

func NewAuthorListController() Controller {
	var a author
	a.path = "/author"
	a.view = view.GetAuthorListTemplate()
	return a
}

type author struct {
	path string
	view *template.Template
}

func (a author) Path() string {
	return a.path
}

func (a author) Controller(conn *model.DBConnection) func(http.ResponseWriter,
	*http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]

		if id == "" {
			a.authorIndex(conn, rw, req)
		} else {
			a.authorId(conn, rw, req, id)
		}
	}
}

func (a author) authorIndex(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request) {

	authors, err := conn.FindAllAuthors()
	if err != nil {
		log.Printf("AuthorController, find all authors: ", err)
	}

	currentUser, currentAuthor := auth.Login(conn, rw, req)

	data := struct {
		CurrentUser   *model.User
		CurrentAuthor *model.Author
		Authors       []model.Author
	}{
		currentUser,
		currentAuthor,
		authors,
	}

	if err := a.view.Execute(rw, data); nil != err {
		log.Printf("AuthorController for Listing", err)
	}
}

func (a author) authorId(conn *model.DBConnection,
	rw http.ResponseWriter, req *http.Request, id string) {

	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("AuthorController, parse id: ", err)
		return
	}

	author, err := conn.FindAuthorById(intId)
	if err != nil {
		log.Printf("AuthorController, author db search: ", err)
		return
	}

	posts, err := author.Posts()
	if err != nil {
		log.Printf("AuthorController, posts db search: ", err)
		return
	}

	currentUser, currentAuthor := auth.Login(conn, rw, req)

	data := struct {
		CurrentUser   *model.User
		CurrentAuthor *model.Author
		Author        *model.Author
		Posts         []model.Post
	}{
		currentUser,
		currentAuthor,
		author,
		posts,
	}

	if err := a.view.Execute(rw, data); nil != err {
		log.Printf("AuthorController for Posts: ", err)
		return
	}
}
