package ctlr

import (
	"github.com/aybabtme/goblog/auth"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	currentUser, currentAuthor := auth.Login(conn, rw, req)

	data := struct {
		CurrentAuthor *model.Author
		CurrentUser   *model.User
		Posts         []model.Post
	}{
		currentAuthor,
		currentUser,
		posts,
	}

	if err := p.view.Execute(rw, data); nil != err {
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

	currentUser, currentAuthor := auth.Login(conn, rw, req)

	post, err := conn.FindPostById(intId)

	data := struct {
		CurrentAuthor *model.Author
		CurrentUser   *model.User
		Post          *model.Post
	}{
		currentAuthor,
		currentUser,
		post,
	}

	if err := p.view.Execute(rw, data); nil != err {
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

	currentUser, currentAuthor := auth.Login(conn, rw, req)

	data := struct {
		CurrentAuthor *model.Author
		CurrentUser   *model.User
		Labels        []model.Label
	}{
		currentAuthor,
		currentUser,
		labels,
	}

	if err := p.view.Execute(rw, data); nil != err {
		log.Println("PostController for listing 2:", err)
		return
	}
}

func (p *post) forSave(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request) {

	title := strings.Title(req.FormValue("title"))
	imageUrl := req.FormValue("imageUrl")
	markdown := req.FormValue("content")
	labelString := req.FormValue("label_list")
	content := string(blackfriday.MarkdownCommon([]byte(markdown)))

	log.Printf("Title=%s\nContent=%s\nLabels=%s", title, content, labelString)

	currentUser, currentAuthor := auth.Login(conn, rw, req)

	if currentAuthor == nil {
		// Can't save posts when you're not an author
		return
	}

	post := conn.NewPost(currentAuthor, title, content, imageUrl, time.Now().UTC())
	if err := post.Save(); err != nil {
		log.Println("Couldn't save post", err)
		return
	}

	data := struct {
		CurrentAuthor *model.Author
		CurrentUser   *model.User
		Post          *model.Post
	}{
		currentAuthor,
		currentUser,
		post,
	}

	if err := p.view.Execute(rw, data); nil != err {
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

	_, currentAuthor := auth.Login(conn, rw, req)
	if currentAuthor == nil {
		http.Redirect(rw, req, "/", 401)
		return
	}

	post, err := conn.FindPostById(intId)

	if post.Author().Id() != currentAuthor.Id() {
		log.Printf("%s tried to delete a post that isn't their.\n",
			currentAuthor.User().Username())
		return
	}

	if err := post.Destroy(); err != nil {
		log.Println("Couldn't delete post:", err)
	}

	http.Redirect(rw, req, "/", 301)

}
