package ctlr

import (
	"github.com/aybabtme/goblog/auth"
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

		currentUser, currentAuthor := auth.Login(conn, rw, req)

		data := struct {
			CurrentAuthor *model.Author
			CurrentUser   *model.User
			Name          string
			AllPosts      []model.Post
		}{
			currentAuthor,
			currentUser,
			label.Name(),
			posts,
		}

		if err := l.view.Execute(rw, data); nil != err {
			log.Println("LabelController, execute:", err)
		}

	}
}
