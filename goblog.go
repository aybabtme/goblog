package main

import (
	"fmt"
	"github.com/aybabtme/goblog/db"
)

type Blog struct {
	title string
}

func NewBlog(title string) *Blog {
	var b = new(Blog)
	b.title = title
	return b
}

func (b *Blog) Start() {
	fmt.Printf("Blog \"%s\" starts\n", b.title)
	db.Start()

	posts, err := db.FindAllPosts()
	if err != nil {
		fmt.Println("Start:", err)
		return
	}

	for idx, val := range posts {
		fmt.Printf("%d : named %s\n", idx, val.Author())
	}
}

func main() {
	b := NewBlog("Patate")
	b.Start()
}
