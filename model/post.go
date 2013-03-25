package model

import (
	"database/sql"
	"fmt"
	"time"
)

//
// SQL queries
//
var createPostTable string = `
CREATE TABLE IF NOT EXISTS Post(
   post_id		SERIAL PRIMARY KEY,
   author_id	INTEGER NOT NULL,
   title			VARCHAR(255) NOT NULL,
   content		TEXT NOT NULL,
   image_url	VARCHAR(255) NOT NULL,
   date			TIMESTAMP NOT NULL,
   CONSTRAINT fk_post_authorid
   	FOREIGN KEY (author_id) REFERENCES Author(author_id) ON DELETE SET NULL
)`

var dropPostTable string = `
DROP TABLE Post;
`

var insertOrReplacePostForId string = `
INSERT INTO Post(
	author_id,
	title,
	content,
	image_url,
	date)
VALUES( $1, $2, $3, $4, $5)`

var findPostById string = `
SELECT
	P.author_id,
	P.title,
	P.content,
	P.image_URL,
	P.date,
   A.user_id,
   U.username,
   U.registration_date,
   U.timezone,
   U.email
FROM
	Post AS P,
	Author AS A,
	BlogUser AS U
WHERE
	P.post_id = $1
	AND P.author_id = A.author_id
	AND A.user_id = U.user_id`

var deletePostById string = `
DELETE FROM
	Post
WHERE
	Post.post_id = $1`

var queryForAllPost string = `
SELECT
	P.post_id,
	P.author_id,
	P.title,
	P.content,
	P.image_url,
	P.date,
	A.user_id,
   U.username,
   U.registration_date,
   U.timezone,
   U.email
FROM
	Post AS P,
	Author AS A,
	BlogUser AS U
WHERE
	P.author_id = A.author_id
	AND A.user_id = U.user_id`

// Relations
var queryForAllCommentsOfPostId string = `
SELECT
	C.comment_id,
	C.user_id,
	C.post_id,
	C.content,
	C.date,
	C.up_vote,
	C.down_vote
FROM
	Comment as C
WHERE
	C.post_id = $1`

var queryForAllLabelsOfPostId string = `
SELECT
	L.label_id,
	L.name
FROM
	Label AS L,
	LabelPost AS LP
WHERE
	LP.post_id = $1
	AND LP.label_id = P.id`

var queryPostIdFromDate string = `
SELECT
	P.post_id
FROM
	Post AS P
WHERE
	P.date = $1
`

// Represents a post in the blog
type Post struct {
	id       int64
	author   *Author
	title    string
	content  string
	imageURL string
	date     time.Time
	model    DBVendor
}

func (p *Post) Id() int64 {
	return p.id
}

func (p *Post) Author() *Author {
	return p.author
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
	model, err := sql.Open(p.model.Driver(), p.model.Name())
	if err != nil {
		fmt.Println("Couldn't open DB:", err)
		return nil, err
	}
	defer model.Close()

	stmt, err := model.Prepare(queryForAllCommentsOfPostId)
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
			model:    p.model,
		}
		comments = append(comments, c)
	}

	return comments, nil
}

//
// SQL stuff
//

//
// Post-specific operations on DBConnection
//

// Create the table Post in the database interface
func (conn *DBConnection) createPostTable() {

	var model = conn.databaser

	db, err := sql.Open(model.Driver(), model.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(createPostTable)
	if err != nil {
		fmt.Printf("Error creating Posts table, driver \"%s\", modelname \"%s\", query = \"%s\"\n",
			model.Driver(), model.Name(), createPostTable)
		fmt.Println(err)
		return
	}
}

func (conn *DBConnection) dropPostTable() {
	var model = conn.databaser

	db, err := sql.Open(model.Driver(), model.Name())
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
func (conn *DBConnection) NewPost(author *Author, title string, content string, imageURL string, date time.Time) *Post {

	return &Post{
		id:       -1,
		author:   author,
		title:    title,
		content:  content,
		imageURL: imageURL,
		date:     date,
		model:    conn.databaser,
	}
}

// Finds all the posts in the database
func (conn *DBConnection) FindAllPosts() ([]Post, error) {

	var posts []Post
	var model = conn.databaser

	db, err := sql.Open(model.Driver(), model.Name())
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
		var userId int64
		var username string
		var registDate time.Time
		var timezone int
		var email string
		err := rows.Scan(&id,
			&authorId,
			&title,
			&content,
			&imageURL,
			&date,
			&userId,
			&username,
			&registDate,
			&timezone,
			&email)
		if err != nil {
			return posts, err
		}
		u := &User{
			id:               userId,
			username:         username,
			registrationDate: registDate,
			timezone:         timezone,
			email:            email,
			model:            model,
		}

		a := &Author{
			id:    authorId,
			user:  u,
			model: model,
		}

		p := Post{
			id:       id,
			author:   a,
			title:    title,
			content:  content,
			imageURL: imageURL,
			date:     date,
			model:    model,
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// Finds a post that matches the given id
func (conn *DBConnection) FindPostById(id int64) (*Post, error) {

	var p *Post
	var model = conn.databaser

	db, err := sql.Open(model.Driver(), model.Name())
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
	var userId int64
	var username string
	var registDate time.Time
	var timezone int
	var email string
	err = stmt.QueryRow(id).Scan(
		&authorId,
		&title,
		&content,
		&imageURL,
		&date,
		&userId,
		&username,
		&registDate,
		&timezone,
		&email)
	if err != nil {
		// normal if the post doesnt exist
		return p, err
	}

	u := &User{
		id:               userId,
		username:         username,
		registrationDate: registDate,
		timezone:         timezone,
		email:            email,
		model:            model,
	}

	a := &Author{
		id:    authorId,
		user:  u,
		model: model,
	}

	p = &Post{
		id:       id,
		author:   a,
		title:    title,
		content:  content,
		imageURL: imageURL,
		date:     date,
		model:    model,
	}

	return p, nil
}

//
// Operations on Post
//

// Saves the post (or update it if it already exists)
// to the database
func (p *Post) Save() error {
	model, err := sql.Open(p.model.Driver(), p.model.Name())
	if err != nil {
		fmt.Println("Save 1:", err)
		return err
	}
	defer model.Close()

	stmt, err := model.Prepare(insertOrReplacePostForId)
	if err != nil {
		fmt.Println("Save 2:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.author.Id(), p.title, p.content, p.imageURL, p.date)
	if err != nil {
		fmt.Println("Save 3:", err)
		return err
	}

	// query the ID we inserted
	idStmt, err := model.Prepare(queryPostIdFromDate)
	if err != nil {
		fmt.Println("Save 5:", err)
		return err
	}
	defer idStmt.Close()

	row := idStmt.QueryRow(p.Date())

	return row.Scan(&p.id)
}

// Deletes the post from the database
func (p *Post) Destroy() error {

	model, err := sql.Open(p.model.Driver(), p.model.Name())
	if err != nil {
		fmt.Println("Destroy:", err)
		return err
	}
	defer model.Close()

	stmt, err := model.Prepare(deletePostById)
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
