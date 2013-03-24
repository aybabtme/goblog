package ctlr

import (
	"fmt"
	"github.com/aybabtme/goblog/db"
	"github.com/gorilla/mux"
	"net/http"
)

func NewLabelController() Controller {
	return new(label)
}

type label string

func (l *label) Path() string {
	return "/label/{key}"
}

func (l *label) Controller(conn *db.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		key := vars["key"]

		// dispatch with conn and rw, req

		fmt.Fprintf(rw, "<h1>I love %s</h1>", key)
	}
}
