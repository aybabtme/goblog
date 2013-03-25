package ctlr

import (
	"fmt"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/goblog/view"
	"html/template"
	"net/http"
)

type index struct {
	view *template.Template
}

func NewIndexController() Controller {
	var i index
	i.view = view.GetIndexTemplate()
	return i
}

func (i index) Path() string {
	return "/"
}

func (i index) Controller(conn *model.DBConnection) func(http.ResponseWriter,
	*http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if err := i.view.Execute(rw, nil); nil != err {
			fmt.Println("IndexController: ", err)
		}
	}
}
