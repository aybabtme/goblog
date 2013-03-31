package ctlr

import (
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func NewAuthorController() Controller {
	var a author
	a.view = view.GetAuthorTemplate()
	return a
}

type author struct {
	view *template.Template
}

type authorData struct {
	Name     string
	AllPosts []model.Post
}

func (a author) Path() string {
	return "/author/{id:[0-9]+}"
}

func (a author) Controller(conn *model.DBConnection) func(http.ResponseWriter,
	*http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("AuthorController, parse id: ", err)
			return
		}

		author, err := conn.FindAuthorById(id)
		if err != nil {
			log.Printf("AuthorController, author db search: ", err)
			return
		}

		posts, err := author.Posts()
		if err != nil {
			log.Printf("AuthorController, posts db search: ", err)
			return
		}

		d := authorData{
			Name:     author.User().Username(),
			AllPosts: posts,
		}

		if err := a.view.Execute(rw, d); nil != err {
			log.Printf("AuthorController for Posts: ", err)
			return
		}
	}
}
