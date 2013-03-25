package ctlr

import (
	"github.com/aybabtme/goblog/model"
	"github.com/gorilla/mux"
	"net/http"
)

func NewAuthorController() Controller {
	return new(author)
}

type author string

func (a *author) Path() string {
	return "/author/{key}"
}

func (a *author) Controller(conn *model.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		_ = vars["key"]

		// dispatch with conn and rw, req

	}
}
