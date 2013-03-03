package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

//
// SQL queries
//
var createTable string = `
CREATE TABLE IF NOT EXISTS Posts(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   author VARCHAR(255),
   content TEXT,
   date TIMESTAMP
)`

var insertOrReplaceRowForId string = `
INSERT OR REPLACE INTO Posts( author, content, date)
VALUES( ?, ?, ?)
`
var findRowById string = `
SELECT P.author, P.content, P.date
FROM Posts AS P
WHERE P.id = ?`

var deleteRowById string = `
DELETE FROM Posts
WHERE Posts.id = ?`

// Represents a post in the blog
type Post struct {
	id      int64
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

	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Error creating Posts table:", err)
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

// Finds a post that match the given id
func FindPostByID(id int64) (Post, error) {
	var p Post

	db, err := sql.Open(DBDriver(), DBName())
	if err != nil {
		fmt.Println("FindPostByID", err)
		return p, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findRowById)
	if err != nil {
		fmt.Println("FindPostByID", err)
		return p, err
	}
	defer stmt.Close()

	var author string
	var content string
	var date time.Time
	err = stmt.QueryRow(id).Scan(&author, &content, &date)
	if err != nil {
		fmt.Println("FindPostByID", err)
		return p, err
	}

	p.id = id
	p.author = author
	p.content = content

	return p, nil
}

// Saves the post (or update it if it already exists)
// to the database
func (p *Post) Save() {
	db, err := sql.Open(DBDriver(), DBName())
	if err != nil {
		fmt.Println("Save", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(insertOrReplaceRowForId)
	if err != nil {
		fmt.Println("Save", err)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.author, p.content, p.date)
	if err != nil {
		fmt.Println("Save", err)
		return
	}

	p.id, _ = res.LastInsertId()
}

// Deletes the post from the database
func (p *Post) Destroy() {

	db, err := sql.Open(DBDriver(), DBName())
	if err != nil {
		fmt.Println("Destroy", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteRowById)
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
