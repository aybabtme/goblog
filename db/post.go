package db

import (
	"database/sql"
	"fmt"
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
   date DATETIME
)`

var insertOrReplaceRowForId string = `
INSERT OR REPLACE INTO Posts( author, content, date)
VALUES( ?, ?, ?)`

var findRowById string = `
SELECT P.author, P.content, P.date
FROM Posts AS P
WHERE P.id = ?`

var deleteRowById string = `
DELETE FROM Posts
WHERE Posts.id = ?`

var queryForAll string = `
SELECT P.id, P.author, P.content, P.date
FROM Posts AS P`

// Represents a post in the blog
type Post struct {
	id      int64
	author  string
	content string
	date    time.Time
	db      Databaser
}

func (p *Post) Id() int64 {
	return p.id
}

func (p *Post) Author() string {
	return p.author
}

func (p *Post) SetAuthor(author string) {
	p.author = author
}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) SetContent(content string) {
	p.content = content
}

func (p *Post) Date() time.Time {
	return p.date
}

func (p *Post) SetDate(time time.Time) {
	p.date = time
}

//
// SQL stuff
//

//
// Post-specific operations on Databaser
//

// Create the table Post in the database interface
func (d *Databaser) createPostTable() {

	db, err := sql.Open(d.name(), d.driver())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Error creating Posts table:", err)
		return
	}
}

// Creates a new Post attached to the Database (but not saved)
func (d *Databaser) NewPost(author string, content string, date time.Time) {

	return &Post{
		author:  author,
		content: content,
		date:    date,
		db:      d,
	}
}

// Finds all the posts in the database
func (d *Databaser) FindAllPosts() ([]Post, error) {
	var posts []Post
	db, err := sql.Open(d.driver(), d.name())
	if err != nil {
		fmt.Println("FindAllPosts:", err)
		return posts, err
	}
	defer db.Close()

	rows, err := db.Query(queryForAll)
	if err != nil {
		fmt.Println("FindAllPosts:", err)
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var author string
		var content string
		var date time.Time
		rows.Scan(&id, &author, &content, &date)
		p := Post{
			id:      id,
			author:  author,
			content: content,
			date:    date,
			db:      d,
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// Finds a post that matches the given id
func (d *Databaser) FindPostById(id int64) (Post, error) {

	db, err := sql.Open(d.driver(), d.Name())
	if err != nil {
		fmt.Println("FindPostById:", err)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findRowById)
	if err != nil {
		fmt.Println("FindPostById:", err)
		return nil, err
	}
	defer stmt.Close()

	var author string
	var content string
	var date time.Time
	err = stmt.QueryRow(id).Scan(&author, &content, &date)
	if err != nil {
		fmt.Println("FindPostById:", err)
		return nil, err
	}

	return &Post{
		id:     id,
		author: author,
		date:   date,
		db:     d,
	}, nil
}

//
// Operations on Post
//

// Saves the post (or update it if it already exists)
// to the database
func (p *Post) Save() {
	db, err := sql.Open(p.db.driver(), p.db.name())
	if err != nil {
		fmt.Println("Save:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(insertOrReplaceRowForId)
	if err != nil {
		fmt.Println("Save:", err)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.author, p.content, p.date)
	if err != nil {
		fmt.Println("Save:", err)
		return
	}

	p.id, _ = res.LastInsertId()
}

// Deletes the post from the database
func (p *Post) Destroy() {

	db, err := sql.Open(p.db.driver(), p.db.name())
	if err != nil {
		fmt.Println("Destroy:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteRowById)
	if err != nil {
		fmt.Println("Destroy:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.id)
	if err != nil {
		fmt.Println("Destroy:", err)
		return
	}

}
