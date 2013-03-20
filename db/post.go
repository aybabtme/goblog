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
   id %s,
   author VARCHAR(255),
   content TEXT,
   date %s
)`

var dropTable string = `
DROP TABLE Posts;
`

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
func (persist *Persister) createPostTable() {

	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	var query = fmt.Sprintf(
		createTable,
		dbaser.IncrementPrimaryKey(),
		dbaser.DateField())

	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating Posts table, query = \"%s\":", query, err)
		return
	}
}

func (persist *Persister) dropPostTable() {
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(dropTable)
	if err != nil {
		fmt.Println("Error droping table", err)
	}

}

// Creates a new Post attached to the Database (but not saved)
func (persist *Persister) NewPost(author string, content string, date time.Time) *Post {

	return &Post{
		id:      -1,
		author:  author,
		content: content,
		date:    date,
		db:      persist.databaser,
	}
}

// Finds all the posts in the database
func (persist *Persister) FindAllPosts() ([]Post, error) {

	var posts []Post
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
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
			db:      dbaser,
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// Finds a post that matches the given id
func (persist *Persister) FindPostById(id int64) (*Post, error) {

	var p *Post
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindPostById 1:", err)
		return p, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findRowById)
	if err != nil {
		fmt.Println("FindPostById 2:", err)
		return p, err
	}
	defer stmt.Close()

	var author string
	var content string
	var date time.Time
	err = stmt.QueryRow(id).Scan(&author, &content, &date)
	if err != nil {
		// normal if the post doesnt exist
		return p, err
	}

	p = &Post{
		id:      id,
		author:  author,
		content: content,
		date:    date,
		db:      dbaser,
	}

	return p, nil
}

//
// Operations on Post
//

// Saves the post (or update it if it already exists)
// to the database
func (p *Post) Save() error {
	db, err := sql.Open(p.db.Driver(), p.db.Name())
	if err != nil {
		fmt.Println("Save:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertOrReplaceRowForId)
	if err != nil {
		fmt.Println("Save:", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.author, p.content, p.date)
	if err != nil {
		fmt.Println("Save:", err)
		return err
	}

	p.id, _ = res.LastInsertId()
	return nil
}

// Deletes the post from the database
func (p *Post) Destroy() error {

	db, err := sql.Open(p.db.Driver(), p.db.Name())
	if err != nil {
		fmt.Println("Destroy:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteRowById)
	if err != nil {
		fmt.Println("Destroy:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.id)
	if err != nil {
		fmt.Println("Destroy:", err)
		return err
	}

	return nil
}
