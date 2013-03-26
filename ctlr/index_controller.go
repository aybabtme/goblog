package ctlr

import (
	"fmt"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"html/template"
	"net/http"
)

type index struct {
	view *template.Template
}

type data struct {
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
			fmt.Println("IndexController: ", err)
			return
		}
		labels, err := conn.FindAllLabels()
		if err != nil {
			fmt.Println("IndexController: ", err)
			return
		}

		d := data{
			AllPosts:  posts,
			AllLabels: labels,
		}

		if err := i.view.Execute(rw, d); nil != err {
			fmt.Println("IndexController: ", err)
			return
		}
	}
}
