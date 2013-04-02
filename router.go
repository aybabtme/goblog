package main

import (
	"github.com/aybabtme/goblog/auth"
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
		ctlr.NewAuthorListController(),
		ctlr.NewUserController(),
		ctlr.NewLabelController(),
		ctlr.NewPostController(),
		ctlr.NewPostComposeController(),
		ctlr.NewPostSaveController(),
		ctlr.NewPostDestroyController(),
		ctlr.NewPostCommentController(),
		ctlr.NewPostIdController()}

	muxer := mux.NewRouter()
	for _, ctlr := range controllers {
		muxer.HandleFunc(ctlr.Path(), ctlr.Controller(conn))
	}
	// serve dynamic resources
	http.Handle("/", muxer)
	// serve static resources
	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir("public/"))))

	// For user authentication
	http.HandleFunc("/authorize", auth.AuthorizeOauth)
	http.HandleFunc("/oauth2callback", auth.GetHandleOAuth2Callback(conn))
	http.HandleFunc("/logout", auth.Logout(conn))

	return http.ListenAndServe(":"+port, nil)
}
