package ctlr

import (
	"github.com/aybabtme/goblog/auth"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
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

func NewPostUpdateController() Controller {
	var p post
	p.path = "/post/save/{saveId:[0-9]+}"
	p.view = view.GetPostTemplate()
	return p
}

func NewPostEditController() Controller {
	var p post
	p.path = "/post/edit/{editId:[0-9]+}"
	p.view = view.GetPostComposeTemplate()
	return p
}

func NewPostCommentController() Controller {
	var p post
	p.path = "/post/comment/{commentId:[0-9]+}"
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
		commentId := vars["commentId"]
		editId := vars["editId"]
		saveId := vars["saveId"]

		if p.path == "/post/compose" {
			p.forCompose(conn, rw, req)
		} else if p.path == "/post/save" {
			p.forSave(conn, rw, req)
		} else if commentId != "" {
			p.forComment(conn, rw, req, commentId)
		} else if destroyId != "" {
			p.forDestroy(conn, rw, req, destroyId)
		} else if editId != "" {
			p.forEdit(conn, rw, req, editId)
		} else if saveId != "" {
			p.forUpdate(conn, rw, req, saveId)
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
		Post          *model.Post
	}{
		currentAuthor,
		currentUser,
		labels,
		nil,
	}

	if err := p.view.Execute(rw, data); nil != err {
		log.Println("PostController for listing 2:", err)
		return
	}
}

func (p *post) forEdit(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request,
	id string) {

	currentUser, currentAuthor := auth.Login(conn, rw, req)

	if currentAuthor == nil {
		http.Error(rw, "/post/"+id, http.StatusForbidden)
		return
	}

	postId, _ := strconv.ParseInt(id, 10, 64)
	post, err := conn.FindPostById(postId)
	if err != nil || post == nil {
		log.Println("Can't edit, post doesn't exist")
		http.Error(rw, "/post/"+id, http.StatusBadRequest)
		return
	}

	labels, err := conn.FindAllLabels()
	if err != nil {
		log.Println("Couldn't find previous labels for autosuggestion")
	}

	data := struct {
		CurrentAuthor *model.Author
		CurrentUser   *model.User
		Labels        []model.Label
		Post          *model.Post
	}{
		currentAuthor,
		currentUser,
		labels,
		post,
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
	content := req.FormValue("content")
	labelString := req.FormValue("label_list")

	currentUser, currentAuthor := auth.Login(conn, rw, req)

	if currentUser == nil {
		http.Redirect(rw, req, "/", http.StatusForbidden)
		return
	}

	if currentAuthor == nil {
		http.Redirect(rw, req, "/", http.StatusForbidden)
		return
	}

	post := conn.NewPost(currentAuthor, title, content, imageUrl, time.Now().UTC())
	if err := post.Save(); err != nil {
		log.Println("Couldn't save post", err)
		return
	}

	for _, label := range strings.Split(labelString, ",") {
		if _, err := post.AddLabel(label); err != nil {
			log.Printf("Couldn't add label <%s> to Post id<%d>\n", label, post.Id())
		}
	}

	id := strconv.FormatInt(post.Id(), 10)
	http.Redirect(rw, req, "/post/"+id, http.StatusFound)
}

func (p *post) forUpdate(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request,
	postId string) {

	_, currentAuthor := auth.Login(conn, rw, req)

	if currentAuthor == nil {
		http.Redirect(rw, req, "/post/"+postId, http.StatusForbidden)
		return
	}

	id, _ := strconv.ParseInt(postId, 10, 64)

	post, err := conn.FindPostById(id)
	if err != nil {
		http.Redirect(rw, req, "/post/"+postId, http.StatusForbidden)
		return
	}

	title := strings.Title(req.FormValue("title"))
	imageUrl := req.FormValue("imageUrl")
	content := req.FormValue("content")
	labelString := req.FormValue("label_list")

	post.SetTitle(title)
	post.SetImageURL(imageUrl)
	post.SetDate(time.Now().UTC())
	post.SetContent(content)
	if err := post.Update(); err != nil {
		log.Println("Couldn't update post", err)
		return
	}

	for _, label := range strings.Split(labelString, ",") {
		if _, err := post.AddLabel(label); err != nil {
			log.Printf("Couldn't add label <%s> to Post id<%d>\n", label, post.Id())
			log.Println(err)
		}
	}

	http.Redirect(rw, req, "/post/"+postId, http.StatusFound)
}

func (p *post) forComment(conn *model.DBConnection,
	rw http.ResponseWriter,
	req *http.Request,
	id string) {

	content := req.FormValue("content")

	currentUser, _ := auth.Login(conn, rw, req)
	if currentUser == nil {
		http.Redirect(rw, req, "/", http.StatusForbidden)
		return
	}

	postId, _ := strconv.ParseInt(id, 10, 64)
	post, err := conn.FindPostById(postId)
	if err != nil || post == nil {
		log.Printf("Post id<%d> doesn't exist", postId)
		log.Println(err)
		http.Redirect(rw, req, "/", http.StatusForbidden)
		return
	}

	comment := conn.NewComment(currentUser.Id(),
		postId,
		content,
		time.Now().UTC())

	if err := comment.Save(); err != nil {
		log.Printf("Error saving comment on post id<%d>\n", postId)
		log.Println(err)
		http.Error(rw, "/", http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, req, "/post/"+id, http.StatusFound)

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

	currentUser, currentAuthor := auth.Login(conn, rw, req)
	if currentUser == nil {
		http.Redirect(rw, req, "/post/"+id, http.StatusForbidden)
		return
	}

	if currentAuthor == nil {
		http.Redirect(rw, req, "/post/"+id, http.StatusForbidden)
		return
	}

	post, err := conn.FindPostById(intId)

	if err := post.Destroy(); err != nil {
		log.Println("Couldn't delete post:", err)
	}

	http.Redirect(rw, req, "/", http.StatusFound)

}
