package main

import (
	"github.com/aybabtme/goblog/ctlr"
	"github.com/aybabtme/goblog/model"
	"github.com/gorilla/mux"
	"net/http"
)

type Router string

func (r Router) Start(port string, conn *model.DBConnection) error {

	controllers := []ctlr.Controller{
		ctlr.NewIndexController(),
		ctlr.NewAuthorController(),
		ctlr.NewLabelController(),
		ctlr.NewPostController(),
		ctlr.NewPostIdController()}

	muxer := mux.NewRouter()
	for _, ctlr := range controllers {
		muxer.HandleFunc(ctlr.Path(), ctlr.Controller(conn))
	}
	http.Handle("/", muxer)
	return http.ListenAndServe(":"+port, nil)
}
