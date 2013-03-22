package db

import (
	"database/sql"
	"fmt"
	"time"
)

var createCommentTable string = `
CREATE TABLE IF NOT EXISTS Comments(
   id %s,
   userId INTEGER,
   postId INTEGER,
   content TEXT,
   date %s,
   upVote INTEGER,
   downVote INTEGER,
   FOREIGN KEY (userId) REFERENCES Users(id) ON DELETE CASCADE,
   FOREIGN KEY (postId) REFERENCES Posts(id) ON DELETE CASCADE
)`

var dropCommentTable string = `
DROP TABLE Comments;`

var insertOrReplaceCommentForId string = `
INSERT OR REPLACE INTO Comments( userId, postId, content, date, upVote, downVote )
VALUES( ?, ?, ?, ?, ?, ? )`

var findCommentById string = `
SELECT C.userId, C.postId, C.content, C.date, C.upVote, C.downVote
FROM Comments as C
WHERE C.id = ?`

var deleteCommentById string = `
DELETE FROM Comments
WHERE Comments.id = ?`

var queryForAllComment string = `
SELECT C.id, C.userId, C.postId, C.content, C.date, C.upVote, C.downVote
FROM Comments AS C`

type Comment struct {
	id       int64
	userId   int64
	postId   int64
	content  string
	date     time.Time
	upVote   int64
	downVote int64
	db       Databaser
}

func (c *Comment) Id() int64 {
	return c.id
}

func (c *Comment) Content() string {
	return c.content
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
 * Comment specific operations on Persister
 */

// Create the table Comment in the database interface
func (persist *Persister) createCommentTable() {

	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	var query = fmt.Sprintf(
		createCommentTable,
		dbaser.IncrementPrimaryKey(),
		dbaser.DateField())

	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating Comments table, driver \"%s\", dbname \"%s\", query = \"%s\"\n",
			dbaser.Driver(), dbaser.Name(), query)
		fmt.Println(err)
		return
	}
}

func (persist *Persister) dropCommentTable() {
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(dropCommentTable)
	if err != nil {
		fmt.Println("Error droping table:", err)
		return
	}

}

func (persist *Persister) NewComment(userId int64, postId int64, content string, date time.Time) *Comment {
	return &Comment{
		id:       -1,
		userId:   userId,
		postId:   postId,
		content:  content,
		date:     date,
		upVote:   0,
		downVote: 0,
		db:       persist.databaser,
	}
}

// Finds all the comments in the database
func (persist *Persister) FindAllComments() ([]Comment, error) {

	var comments []Comment
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindAllComments 1:", err)
		return comments, err
	}
	defer db.Close()

	rows, err := db.Query(queryForAllComment)
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
			db:       dbaser,
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// Finds a comment that matches the given id
func (persist *Persister) FindCommentById(id int64) (*Comment, error) {

	var c *Comment
	var dbaser = persist.databaser

	db, err := sql.Open(dbaser.Driver(), dbaser.Name())
	if err != nil {
		fmt.Println("FindCommentById 1:", err)
		return c, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findCommentById)
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
		db:       dbaser,
	}

	return c, nil
}

/*
 * Operations on Comment
 */

// Saves the post (or update it if it already exists)
// to the database
func (c *Comment) Save() error {
	db, err := sql.Open(c.db.Driver(), c.db.Name())
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

	res, err := stmt.Exec(c.userId, c.postId, c.content, c.date, c.upVote, c.downVote)
	if err != nil {
		fmt.Println("Save 3:", err)
		return err
	}

	c.id, _ = res.LastInsertId()
	return nil
}

// Deletes the comment from the database
func (c *Comment) Destroy() error {

	db, err := sql.Open(c.db.Driver(), c.db.Name())
	if err != nil {
		fmt.Println("Destroy:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteCommentById)
	if err != nil {
		fmt.Println("Destroy:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.id)
	if err != nil {
		fmt.Println("Destroy:", err)
		return err
	}

	return nil
}
