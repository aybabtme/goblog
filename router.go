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
		ctlr.NewUserController(),
		ctlr.NewLabelController(),
		ctlr.NewPostController(),
		ctlr.NewPostComposeController(),
		ctlr.NewPostSaveController(),
		ctlr.NewPostDestroyController(),
		ctlr.NewPostIdController()}

	muxer := mux.NewRouter()
	for _, ctlr := range controllers {
		muxer.HandleFunc(ctlr.Path(), ctlr.Controller(conn))
	}
	// serve dynamic resources
	http.Handle("/", muxer)
	// serve static resources
	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir("public/"))))
	return http.ListenAndServe(":"+port, nil)
}
