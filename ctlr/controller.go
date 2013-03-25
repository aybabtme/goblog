package ctlr

import (
	"github.com/aybabtme/goblog/model"
	"net/http"
)

type Controller interface {
	Path() string
	Controller(*model.DBConnection) func(http.ResponseWriter, *http.Request)
}
