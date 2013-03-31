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

type post struct {
	path string
	view *template.Template
}

func NewPostController() Controller {
	var p post
	p.path = "/post"
	p.view = view.GetPostListingTemplate()
	return p
}

func NewPostIdController() Controller {
	var p post
	p.path = "/post/{id:[0-9]+}"
	p.view = view.GetPostTemplate()
	return p
}

func (p post) Path() string {
	return p.path
}

func (p post) Controller(conn *model.DBConnection) func(http.ResponseWriter, *http.Request) {
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

func (p *post) forListing(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request) {

	posts, err := conn.FindAllPosts()
	if err != nil {
		log.Println("PostController for listing 1:", err)
		return
	}
	if err := p.view.Execute(rw, posts); nil != err {
		log.Println("PostController for listing 2:", err)
		return
	}

}

func (p *post) forId(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request,
	id string) {

	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("PostController for id 1:", err)
		return
	}
	post, err := conn.FindPostById(intId)
	if err := p.view.Execute(rw, post); nil != err {
		log.Println("PostController for id 3:", err)
		return
	}

}
