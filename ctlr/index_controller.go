package ctlr

import (
	"fmt"
	"github.com/aybabtme/goblog/db"
	"net/http"
)

func NewIndexController() Controller {
	return new(index)
}

type index string

func (i *index) Path() string {
	return "/"
}

func (i *index) Controller(conn *db.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		// dispatch with conn and rw, req
		posts, err := conn.FindAllPosts()
		if err != nil {
			fmt.Fprintf(rw, "<h1>Something went wrong!</h1>")
			return
		}
		if len(posts) == 0 {
			fmt.Fprintf(rw, "<h1>No post to show!</h1><h2>Fill your data!</h2>")
			return
		}

		var template string = `
<h1>Post #%d</h1>
<h2>Titled: %s</h2>
<img src="%s" height="100" width="100"></img>
<h3>By %s posted on %s</h3>
<p>%s<p>
</hr>
      `
		for i, p := range posts {
			author, err := conn.FindAuthorById(p.AuthorId())
			if err != nil {
				fmt.Fprintf(rw, "Oups, something went wrong")
				return
			}
			fmt.Fprintf(rw, template, i,
				p.Title(),
				p.ImageURL(),
				author.User().Username(),
				p.Date(),
				p.Content())
		}

	}
}
