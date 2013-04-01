package ctlr

import (
	"github.com/aybabtme/goblog/auth"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"html/template"
	"log"
	"net/http"
)

type index struct {
	view *template.Template
}

func NewIndexController() Controller {
	var i index
	i.view = view.GetIndexTemplate()
	return i
}

func (i index) Path() string {
	return "/"
}

func (i index) Controller(conn *model.DBConnection) func(http.ResponseWriter,
	*http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

		currentUser, currentAuthor := auth.Login(conn, rw, req)

		posts, err := conn.FindAllPosts()
		if err != nil {
			log.Println("IndexController, list posts: ", err)
			return
		}
		labels, err := conn.FindAllLabels()
		if err != nil {
			log.Println("IndexController, list labels: ", err)
			return
		}

		data := struct {
			CurrentUser   *model.User
			CurrentAuthor *model.Author
			AllPosts      []model.Post
			AllLabels     []model.Label
		}{
			currentUser,
			currentAuthor,
			posts,
			labels,
		}

		if err := i.view.Execute(rw, data); nil != err {
			log.Println("IndexController, execute: ", err)
			return
		}
	}
}
