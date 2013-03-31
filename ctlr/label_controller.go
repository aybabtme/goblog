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

func NewLabelController() Controller {
	var l label
	l.view = view.GetLabelTemplate()
	return l
}

type label struct {
	view *template.Template
}

type labelData struct {
	Name     string
	AllPosts []model.Post
}

func (l label) Path() string {
	return "/label/{id:[0-9]+}"
}

func (l label) Controller(conn *model.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Println("LabelController, parse id:", err)
			return
		}

		label, err := conn.FindLabelById(id)
		if err != nil {
			log.Printf("LabelController, for id(%d): \n%v\n", id, err)
			return
		}

		posts, err := label.Posts()
		if err != nil {
			log.Println("LabelController, listing posts.", err)
			return
		}

		d := labelData{
			Name:     label.Name(),
			AllPosts: posts,
		}

		if err := l.view.Execute(rw, d); nil != err {
			log.Println("LabelController, execute:", err)
		}

	}
}
