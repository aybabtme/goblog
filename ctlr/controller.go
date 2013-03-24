package ctlr

import (
	"github.com/aybabtme/goblog/db"
	"net/http"
)

type Controller interface {
	Path() string
	Controller(*db.DBConnection) func(*http.ResponseWriter, *http.Request)
}
