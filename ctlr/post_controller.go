package ctlr

import (
	"fmt"
	"github.com/aybabtme/goblog/db"
	"github.com/gorilla/mux"
	"net/http"
)

type post struct {
	path string
}

func NewPostController() Controller {
	return &post{path: "/post"}
}

func NewPostIdController() Controller {
	return &post{path: "/post/{id:[0-9]+}"}
}

func (p *post) Path() string {
	return p.path
}

func (p *post) Controller(conn *db.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]

		if id == "" {
			p.forListing(conn, rw, req)
		} else {
			p.forId(conn, rw, req, id)
		}
	}
}

func (p *post) forListing(conn *db.DBConnection,
	rw http.ResponseWriter,
	req *http.Request) {

	// when templates will be there, we'll need to pass them data
	fmt.Fprintf(rw, "received listing request")
}

func (p *post) forId(conn *db.DBConnection,
	rw http.ResponseWriter,
	req *http.Request,
	id string) {

	// when templates will be there, we'll need to pass them data
	fmt.Fprintf(rw, "received id %s", id)
}
