package model

import (
	"database/sql"
	"fmt"
	"github.com/russross/blackfriday"
	"time"
)

/*
 * SQL strings
 */

var createCommentTable string = `
CREATE TABLE IF NOT EXISTS Comment(
   comment_id	SERIAL PRIMARY KEY,
   user_id		INTEGER NOT NULL,
   post_id		INTEGER NOT NULL,
   content		TEXT NOT NULL,
   date			TIMESTAMP NOT NULL,
   up_vote		INTEGER NOT NULL,
   down_vote	INTEGER NOT NULL,
   CONSTRAINT fk_comment_user_id
   	FOREIGN KEY (user_id) REFERENCES BlogUser(user_id) ON DELETE CASCADE,
   CONSTRAINT fk_comment_post_id
   	FOREIGN KEY (post_id) REFERENCES Post(post_id) ON DELETE CASCADE
)`

var dropCommentTable string = `
DROP TABLE Comment;`

var insertOrReplaceCommentForId string = `
INSERT INTO Comment(
	user_id,
	post_id,
	content,
	date,
	up_vote,
	down_vote )
VALUES( $1, $2, $3, $4, $5, $6 )`

var findCommentById string = `
SELECT
	C.user_id,
	C.post_id,
	C.content,
	C.date,
	C.up_vote,
	C.down_vote
FROM
	Comment as C
WHERE
	C.comment_id = $1`

var deleteCommentById string = `
DELETE FROM Comment
WHERE Comment.comment_id = $1`

var queryForAllComment string = `
SELECT
	C.comment_id,
	C.user_id,
	C.post_id,
	C.content,
	C.date,
	C.up_vote,
	C.down_vote
FROM Comment AS C`

var queryCommentIdFromDetails string = `
SELECT
	C.comment_id
FROM
	Comment AS C
WHERE
	C.date = $1
`

// Represents a comment on a post.  Comments are made by Users.
type Comment struct {
	id       int64
	userId   int64
	postId   int64
	content  string
	date     time.Time
	upVote   int64
	downVote int64
	conn     *DBConnection
}

func (c *Comment) Id() int64 {
	return c.id
}

func (c *Comment) User() (*User, error) {
	return c.conn.FindUserById(c.userId)
}

func (c *Comment) Post() (*Post, error) {
	return c.conn.FindPostById(c.postId)
}

func (c *Comment) Content() string {
	return c.content
}

func (c *Comment) ContentMarkdown() string {
	return string(blackfriday.MarkdownCommon([]byte(c.content)))
}

func (c *Comment) SetContent(content string) {
	c.content = content
}

func (c *Comment) Date() time.Time {
	return c.date
}

func (c *Comment) SetDate(date time.Time) {
	c.date = date
}

func (c *Comment) UpVote() int64 {
	return c.upVote
}

func (c *Comment) SetUpVote(count int64) {
	c.upVote = count
}

func (c *Comment) DownVote() int64 {
	return c.downVote
}

func (c *Comment) SetDownVote(count int64) {
	c.downVote = count
}

/*
 * SQL stuff
 */

/*
 * Comment specific operations on DBConnection
 */

// Create the table Comment in the database interface
func (persist *DBConnection) createCommentTable() {

	var modelaser = persist.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer model.Close()

	_, err = model.Exec(createCommentTable)
	if err != nil {
		fmt.Printf("Error creating Comments table, driver \"%s\", modelname \"%s\", query = \"%s\"\n",
			modelaser.Driver(), modelaser.Name(), createCommentTable)
		fmt.Println(err)
		return
	}
}

func (persist *DBConnection) dropCommentTable() {
	var modelaser = persist.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer model.Close()

	_, err = model.Exec(dropCommentTable)
	if err != nil {
		fmt.Println("Error droping table:", err)
		return
	}

}

// Creates a new Comment attached to the database.  It is NOT saved in the
// database, you must call "Save" on this comment to have it persisted
func (conn *DBConnection) NewComment(userId int64, postId int64, content string, date time.Time) *Comment {
	return &Comment{
		id:       -1,
		userId:   userId,
		postId:   postId,
		content:  content,
		date:     date,
		upVote:   0,
		downVote: 0,
		conn:     conn,
	}
}

// Finds all the comments in the database.  Returns an empty slice with
// an error if not comments were found.
func (conn *DBConnection) FindAllComments() ([]Comment, error) {

	var comments []Comment
	var modelaser = conn.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("FindAllComments 1:", err)
		return comments, err
	}
	defer model.Close()

	rows, err := model.Query(queryForAllComment)
	if err != nil {
		fmt.Println("FindAllComments 2:", err)
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
		c := Comment{
			id:       id,
			userId:   userId,
			postId:   postId,
			content:  content,
			date:     date,
			upVote:   upVote,
			downVote: downVote,
			conn:     conn,
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// Finds a comment that matches the given id.  Returns nil and an error
// if the id didn't match any comment.
func (conn *DBConnection) FindCommentById(id int64) (*Comment, error) {

	var c *Comment
	var modelaser = conn.databaser

	model, err := sql.Open(modelaser.Driver(), modelaser.Name())
	if err != nil {
		fmt.Println("FindCommentById 1:", err)
		return c, err
	}
	defer model.Close()

	stmt, err := model.Prepare(findCommentById)
	if err != nil {
		fmt.Println("FindCommentById 2:", err)
		return c, err
	}
	defer stmt.Close()

	var userId int64
	var postId int64
	var content string
	var date time.Time
	var upVote int64
	var downVote int64
	err = stmt.QueryRow(id).Scan(&userId, &postId, &content, &date, &upVote, &downVote)
	if err != nil {
		// normal if the comment doesnt exist
		return c, err
	}

	c = &Comment{
		id:       id,
		userId:   userId,
		postId:   postId,
		content:  content,
		date:     date,
		upVote:   upVote,
		downVote: downVote,
		conn:     conn,
	}

	return c, nil
}

/*
 * Operations on Comment
 */

// Saves the post (or update it if it already exists)
// to the database.  Returns an error if something went wrong.
func (c *Comment) Save() error {
	model := c.conn.databaser
	db, err := sql.Open(model.Driver(), model.Name())
	if err != nil {
		fmt.Println("Save 1:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertOrReplaceCommentForId)
	if err != nil {
		fmt.Println("Save 2:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.userId, c.postId, c.content, c.date, c.upVote, c.downVote)
	if err != nil {
		fmt.Println("Save 3:", err)
		return err
	}

	// query the ID we inserted
	idStmt, err := db.Prepare(queryCommentIdFromDetails)
	if err != nil {
		fmt.Println("Save 5:", err)
		return err
	}
	defer idStmt.Close()

	row := idStmt.QueryRow(c.date)

	return row.Scan(&c.id)
}

// Deletes the comment from the database.  Returns an error if something
// went wrong.
func (c *Comment) Destroy() error {

	model := c.conn.databaser
	db, err := sql.Open(model.Driver(), model.Name())
	if err != nil {
		fmt.Println("Comment Destroy 1:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteCommentById)
	if err != nil {
		fmt.Println("Comment Destroy 2:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.id)
	if err != nil {
		fmt.Println("Comment Destroy 3:", err)
		return err
	}

	return nil
}
