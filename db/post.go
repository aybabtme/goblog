package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Post struct {
	id      int
	author  string
	content string
	date    time.Time
}

func init() {
	db, err := sql.Open(DBDriver(), DBName())
	if err != nil {
		fmt.Println("Error on Post init", err)
		return
	}
	defer db.Close()

	createTable := `
   CREATE TABLE Posts(
      id INTEGER NOT NULL PRIMARY KEY,
      author VARCHAR(255),
      content TEXT,
      date TIMESTAMP,
   )`

	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Error creating Posts table", err)
		return
	}

}

func NewPost(author, content string) *Post {

	p := new(Post)
	p.author = author
	p.content = content
	p.date = time.Now().UTC()
	return p
}

func FindPostByID(id int) Post {
	db, err := sql.Open(DBDriver(), DBName())
	if err != nil {
		fmt.Println("FindPostByID", err)
		return
	}
	defer db.Close()

}

func (p *Post) Save() {
	db, err := sql.Open(DBDriver(), DBName())
	if err != nil {
		fmt.Println("Save", err)
		return
	}
	defer db.Close()
}

func (p *Post) Destroy() {

	db, err := sql.Open(DBDriver(), DBName())
	if err != nil {
		fmt.Println("Destroy", err)
		return
	}
	defer db.Close()

	deleteQuery := `
   DELETE FROM Posts
   WHERE Post.id = ?
   `

	stmt, err := db.Prepare(deleteQuery)
	if err != nil {
		fmt.Println("Destroy", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.id)
	if err != nil {
		fmt.Println("Destroy", err)
		return
	}

}
