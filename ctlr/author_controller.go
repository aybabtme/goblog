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

type authroData struct {
	Name     string
	AllPosts []model.Post
}

func (l label) Path() string {
	return "/label?{id:[0-9]+}"
}

func (a author) Path() string {
	return a.path
}

func (a author) Controller(conn *model.DBConnection) func(http.ResponseWriter,
	*http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			fmt.Println("AuthorController, parse id: ", err)
			return
		}

		author, err := conn.FindAuthorById(id)
		if err != nil {
			fmt.Println("AuthorController, author db search: ", err)
			return
		}

		posts, err := author.Posts()
		if err != nil {
			fmt.Println("AuthorController, posts db search: ", err)
			return
		}

		d := authroData{
			Name:     author.User().Username(),
			AllPosts: posts,
		}

		if err := a.view.Execute(rw, d); nil != err {
			fmt.Println("AuthorController for Posts: ", err)
			return
		}
	}
}
