package ctlr

import (
	"github.com/aybabtme/goblog/db"
	"github.com/gorilla/mux"
	"net/http"
)

func NewPostController() Controller {
	return new(post)
}

type post string

func (p *post) Path() string {
	return "/post/{key}"
}

func (p *post) Controller(conn *db.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		_ = vars["key"]

		// dispatch with conn and rw, req

	}
}
