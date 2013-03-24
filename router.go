package main

import (
	"github.com/aybabtme/goblog-prototype/db"
	"net/http"
)

type Router string

func (r Router) Start(port string, conn *db.DBConnection) {
	user := CurrentUser{false}

	http.HandleFunc("/", indexController(conn))
	http.HandleFunc("/post/", postController(conn))
	http.HandleFunc("/author/", authorController(conn))
	http.HandleFunc("/label/", labelController(conn))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
