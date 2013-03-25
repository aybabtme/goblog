package ctlr

import (
	"fmt"
	"github.com/aybabtme/goblog/model"
	"net/http"
)

func NewIndexController() Controller {
	return new(index)
}

type index string

func (i *index) Path() string {
	return "/"
}

func (i *index) Controller(conn *model.DBConnection) func(http.ResponseWriter,
	*http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

		fmt.Fprintf(rw, "this is the index page!")

	}
}
