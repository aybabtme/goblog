package ctlr

import (
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"html/template"
	"log"
	"net/http"
)

type index struct {
	view *template.Template
}

type indexData struct {
	AllLabels []model.Label
	AllPosts  []model.Post
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

		d := indexData{
			AllPosts:  posts,
			AllLabels: labels,
		}

		if err := i.view.Execute(rw, d); nil != err {
			log.Println("IndexController, execute: ", err)
			return
		}
	}
}
