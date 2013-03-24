package ctlr

import (
	"github.com/aybabtme/goblog/db"
	"github.com/gorilla/mux"
	"net/http"
)

func NewIndexController() Controller {
	return new(index)
}

type index string

func (i *index) Path() string {
	return "/index/{key}"
}

func (i *index) Controller(conn *db.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		_ = vars["key"]

		// dispatch with conn and rw, req

	}
}
