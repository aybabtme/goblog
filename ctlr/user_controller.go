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

func NewUserController() Controller {
	var u user
	u.view = view.GetUserTemplate()
	return u
}

type user struct {
	view *template.Template
}

func (u user) Path() string {
	return "/user/{id:[0-9]+}"
}

func (u user) Controller(conn *model.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Println("UserController, parse id:", err)
			return
		}

		user, err := conn.FindUserById(id)
		if err != nil {
			log.Printf("UserController, for id(%d): \n%v\n", id, err)
			return
		}

		if err := u.view.Execute(rw, user); nil != err {
			log.Println("UserController, execute:", err)
		}

	}
}
