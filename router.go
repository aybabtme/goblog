package main

import (
	"github.com/aybabtme/goblog/ctlr"
	"github.com/aybabtme/goblog/db"
	"net/http"
)

type Router string

func (r Router) Start(port string, conn *db.DBConnection) error {

	controllers := []Controller{
		ctlr.NewIndexController(),
		ctlr.NewAuthorController(),
		ctlr.NewLabelController(),
		ctlr.NewPostController()}

	for _, ctlr := range controllers {
		http.HandleFunc(ctlr.Path(), ctlr.Controller(conn))
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
}
