package ctlr

import (
	"fmt"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

type author struct {
	path string
	view *template.Template
}

type data struct {
	AllPosts  []model.Post
	AllLabels []model.Label
}

func NewAuthorController() Controller {
	var a author
	a.path = "/author"
	a.view = view.GetAuthorListingTemplate()
	return a
}

func NewAuthorIdController() Controller {
	var a author
	a.path = "/author/{id:[0-9]+}"
	a.view = view.GetAuthorTemplate()
	return a
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
			a.authorPosts(conn, rw, req, id)
		}
	}
}

func (a *author) authorIndex(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request) {

	authors, err := conn.FindAllAuthors()
	if err != nil {
		fmt.Println("authorController for listing 1:", err)
		return
	}
	if err := a.view.Execute(rw, authors); nil != err {
		fmt.Println("authorController for listing 2:", err)
		return
	}
}

func (a *author) authorPosts(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request,
	id string) {

	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("AuthorController for Posts: ", err)
		return
	}
	author, err := conn.FindAuthorById(intId)
	if err != nil {
		fmt.Println("AuthorController for Posts: ", err)
		return
	}
	posts, err := author.Posts()
	if err != nil {
		fmt.Println("AuthorController for Posts: ", err)
		return
	}
	labels, err := conn.FindAllLabels()
	if err != nil {
		fmt.Println("AuthorController for Posts: ", err)
		return
	}

	d := data{
		AllPosts:  posts,
		AllLabels: labels,
	}

	if err := a.view.Execute(rw, d); nil != err {
		fmt.Println("AuthorController for Posts: ", err)
		return
	}

}
