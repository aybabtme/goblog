package ctlr

import (
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func NewPostComposeController() Controller {
	var p post
	p.path = "/post/compose"
	p.view = view.GetPostComposeTemplate()
	return p
}

func NewPostSaveController() Controller {
	var p post
	p.path = "/post/save"
	p.view = view.GetPostTemplate()
	return p
}

func NewPostIdController() Controller {
	var p post
	p.path = "/post/{id:[0-9]+}"
	p.view = view.GetPostTemplate()
	return p
}

func NewPostDestroyController() Controller {
	var p post
	p.path = "/post/destroy/{destroyId:[0-9]+}"
	p.view = view.GetPostDestroyTemplate()
	return p
}

func (p post) Path() string {
	return p.path
}

func (p post) Controller(conn *model.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		destroyId := vars["destroyId"]

		if p.path == "/post/compose" {
			p.forCompose(conn, rw, req)
		} else if p.path == "/post/save" {
			p.forSave(conn, rw, req)
		} else if destroyId != "" {
			p.forDestroy(conn, rw, req, destroyId)
		} else if id == "" {
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

func (p *post) forCompose(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request) {

	labels, err := conn.FindAllLabels()
	if err != nil {
		log.Println("Couldn't find previous labels for autosuggestion")
	}

	if err := p.view.Execute(rw, labels); nil != err {
		log.Println("PostController for listing 2:", err)
		return
	}
}

func (p *post) forSave(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request) {

	title := strings.Title(req.FormValue("title"))
	content := req.FormValue("content")
	labelString := req.FormValue("label_list")

	log.Printf("Title=%s\nContent=%s\nLabels=%s", title, content, labelString)

	if err := p.view.Execute(rw, nil); nil != err {
		log.Println("PostController for listing 2:", err)
		return
	}
}

func (p *post) forDestroy(conn *model.DBConnection,
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
