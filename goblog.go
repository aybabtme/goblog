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
}

func main() {
	b := NewBlog("Patate")
	b.Start()
}
