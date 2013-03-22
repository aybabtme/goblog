package db

import (
	"database/sql"
	"fmt"
	"time"
)

//
// SQL queries
//
var createPostTable string = `
CREATE TABLE IF NOT EXISTS Posts(
   id %s,
   authorId INTEGER,
   title VARCHAR(255),
   content TEXT,
   imageURL VARCHAR(255),
   date %s,
   FOREIGN KEY (authorId) REFERENCES Authors(id) ON DELETE SET NULL
)`

var dropPostTable string = `
DROP TABLE Posts;
`

var insertOrReplacePostForId string = `
INSERT OR REPLACE INTO Posts( authorId, title, content, imageURL, date)
VALUES( ?, ?, ?, ?, ?)`

var findPostById string = `
SELECT P.authorId, P.title, P.content, P.imageURL, P.date
FROM Posts AS P
WHERE P.id = ?`

var deletePostById string = `
DELETE FROM Posts
WHERE Posts.id = ?`

var queryForAllPost string = `
SELECT P.id, P.authorId, P.title, P.content, P.imageURL, P.date
FROM Posts AS P`

// Relations
var queryForAllCommentsOfPostId string = `
SELECT C.id, C.userId, C.postId, C.content, C.date, C.upVote, C.downVote
FROM Comments as C
WHERE C.postId = ?`

var queryForAllLabelsOfPostId string = `
SELECT L.id, L.name
FROM Labels AS L, LabelPosts AS LP
WHERE LP.postId = ? AND LP.labelId = P.id`

// Represents a post in the blog
type Post struct {
	id       int64
	authorId int64
	title    string
	content  string
	imageURL string
	date     time.Time
	db       Databaser
}

func (p *Post) Id() int64 {
	return p.id
}

func (p *Post) AuthorId() int64 {
	return p.authorId
}

func (p *Post) SetAuthorId(id int64) {
	p.authorId = id
}

func (p *Post) Title() string {
	return p.title
}

func (p *Post) SetTitle(title string) {
	p.title = title
}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) SetContent(content string) {
	p.content = content
}

func (p *Post) ImageURL() string {
	return p.imageURL
}

func (p *Post) SetImageURL(imageURL string) {
	p.imageURL = imageURL
}

func (p *Post) Date() time.Time {
	return p.date
}

func (p *Post) SetDate(time time.Time) {
	p.date = time
}

func (p *Post) Comments() ([]Comment, error) {
	db, err := sql.Open(p.db.Driver(), p.db.Name())
	if err != nil {
		fmt.Println("Couldn't open DB:", err)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(queryForAllCommentsOfPostId)
	if err != nil {
		fmt.Printf("Couldn't prepare statement: %s", queryForAllCommentsOfPostId)
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()

	var comments []Comment

	rows, err := stmt.Query(p.id)
	if err != nil {
		fmt.Println("Couldn't read rows from statement", err)
		return comments, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var userId int64
		var postId int64
		var content string
		var date time.Time
		var upVote int64
		var downVote int64
		rows.Scan(&id, &userId, &postId, &content, &date, &upVote, &downVote)
		if err != nil {
			fmt.Println("Error while scanning comments", err)
			return comments, err
		}
		c := Comment{
			id:       id,
			userId:   userId,
			postId:   postId,
			content:  content,
			date:     date,
			upVote:   upVote,
			downVote: downVote,
			db:       p.db,
		}
		comments = append(comments, c)
	}

	return comments, nil
}

//
// SQL stuff
//

//
// Post-specific operations on Persister
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
		createPostTable,
		dbaser.IncrementPrimaryKey(),
		dbaser.DateField())

	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating Posts table, driver \"%s\", dbname \"%s\", query = \"%s\"\n",
			dbaser.Driver(), dbaser.Name(), query)
		fmt.Println(err)
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

	_, err = db.Exec(dropPostTable)
	if err != nil {
		fmt.Println("Error droping table:", err)
	}

}

// Creates a new Post attached to the Database (but not saved)
func (persist *Persister) NewPost(authorId int64, title string, content string, imageURL string, date time.Time) *Post {

	return &Post{
		id:       -1,
		authorId: authorId,
		title:    title,
		content:  content,
		imageURL: imageURL,
		date:     date,
		db:       persist.databaser,
	}
}

// Finds all the posts in the database
func (persist *Persister) FindAllPosts() ([]Post, error) {

	var posts []Post
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindAllPosts 1:", err)
		return posts, err
	}
	defer db.Close()

	rows, err := db.Query(queryForAllPost)
	if err != nil {
		fmt.Println("FindAllPosts 2:", err)
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var authorId int64
		var title string
		var content string
		var imageURL string
		var date time.Time
		err := rows.Scan(&id, &authorId, &title, &content, &imageURL, &date)
		if err != nil {
			return posts, err
		}
		p := Post{
			id:       id,
			authorId: authorId,
			title:    title,
			content:  content,
			imageURL: imageURL,
			date:     date,
			db:       dbaser,
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

	stmt, err := db.Prepare(findPostById)
	if err != nil {
		fmt.Println("FindPostById 2:", err)
		return p, err
	}
	defer stmt.Close()

	var authorId int64
	var title string
	var content string
	var imageURL string
	var date time.Time
	err = stmt.QueryRow(id).Scan(&authorId, &title, &content, &imageURL, &date)
	if err != nil {
		// normal if the post doesnt exist
		return p, err
	}

	p = &Post{
		id:       id,
		authorId: authorId,
		title:    title,
		content:  content,
		imageURL: imageURL,
		date:     date,
		db:       dbaser,
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
		fmt.Println("Save 1:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertOrReplacePostForId)
	if err != nil {
		fmt.Println("Save 2:", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.authorId, p.title, p.content, p.imageURL, p.date)
	if err != nil {
		fmt.Println("Save 3:", err)
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

	stmt, err := db.Prepare(deletePostById)
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
