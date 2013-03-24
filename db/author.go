package db

import (
	"database/sql"
	"fmt"
	"time"
)

/*
 * SQL stuff
 */
var createAuthorTable string = `
CREATE TABLE IF NOT EXISTS Author(
   author_id SERIAL PRIMARY KEY,
   user_id INTEGER UNIQUE NOT NULL,
   CONSTRAINT fk_author_user_id
   	FOREIGN KEY(user_id) REFERENCES BlogUser(user_id) ON DELETE CASCADE
)`

var dropAuthorTable string = `
DROP TABLE Author;
`

var insertAuthorForId string = `
INSERT INTO Author(user_id) VALUES ( $1 )`

var findAuthorById string = `
SELECT
   A.user_id,
   U.username,
   U.registration_date,
   U.timezone,
   U.email
FROM Author AS A, BlogUser AS U
WHERE A.author_id = $1 AND A.user_id = U.user_id`

var deleteAuthorById string = `
DELETE FROM Author
WHERE Author.author_id = $1`

var queryForAllAuthor string = `
SELECT
   A.author_id,
   A.user_id,
   U.username,
   U.registration_date,
   U.timezone,
   U.email
FROM Author AS A, BlogUser AS U
WHERE A.user_id = U.user_id`

var queryAuthorForUserId string = `
SELECT A.author_id
FROM Author AS A
WHERE A.user_id = $1`

// Relations
var queryForAllPostsOfAuthorId string = `
SELECT P.post_id, P.author_id, P.title, P.content, P.image_url, P.date
FROM Post AS P
WHERE P.author_id = $1
`

// Represents an author of the blog
type Author struct {
	id     int64
	userId int64
	user   *User
	db     DBVendor
}

func (a *Author) Id() int64 {
	return a.id
}

func (a *Author) UserId() int64 {
	return a.userId
}

func (a *Author) User() *User {
	return a.user
}

func (a *Author) Posts() ([]Post, error) {
	db, err := sql.Open(a.db.Driver(), a.db.Name())
	if err != nil {
		fmt.Println("Couldn't open DB:", err)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(queryForAllPostsOfAuthorId)
	if err != nil {
		fmt.Printf("Couldn't prepare statement: %s", queryForAllPostsOfAuthorId)
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()

	var posts []Post

	rows, err := stmt.Query(a.id)
	if err != nil {
		fmt.Println("Couldn't read rows from statement", err)
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
		rows.Scan(&id, &authorId, &title, &content, &imageURL, &date)
		if err != nil {
			fmt.Println("Error while scanning comments", err)
			return posts, err
		}
		p := Post{
			id:       id,
			authorId: authorId,
			title:    title,
			content:  content,
			imageURL: imageURL,
			date:     date,
			db:       a.db,
		}
		posts = append(posts, p)
	}
	return posts, nil
}

/*
 *  SQL Stuff
 */

func (p *DBConnection) createAuthorTable() {
	var dbvendor = p.databaser

	db, err := sql.Open(dbvendor.Driver(), dbvendor.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}

	_, err = db.Exec(createAuthorTable)
	if err != nil {
		fmt.Printf("Error creating Author table, driver \"%s\", dbname \"%s\", query = %s\n",
			dbvendor.Driver(), dbvendor.Name(), createAuthorTable)
		fmt.Println(err)
		return
	}
}

func (conn *DBConnection) dropAuthorTable() {
	var dbvendor = conn.databaser

	db, err := sql.Open(dbvendor.Driver(), dbvendor.Name())
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(dropAuthorTable)
	if err != nil {
		fmt.Println("Error droping table:", err)
	}
}

// Creates an author.  The Author is NOT saved.  To save it, you must call
// the save method on the returned Author.
func (conn *DBConnection) NewAuthor(user *User) *Author {
	return &Author{
		id:     -1,
		userId: user.Id(),
		user:   user,
		db:     conn.databaser,
	}
}

// Finds all the Authors known to this blog.  Returns an empty slice and
// an error stating no rows matched the request if no authors are known
// to this blog.
func (conn *DBConnection) FindAllAuthors() ([]Author, error) {

	var authors []Author
	var dbvendor = conn.databaser

	db, err := sql.Open(dbvendor.Driver(), dbvendor.Name())
	if err != nil {
		fmt.Println("FindAllAuthors 1:", err)
		return authors, err
	}
	defer db.Close()

	rows, err := db.Query(queryForAllAuthor)
	if err != nil {
		fmt.Println("FindAllAuthors 2:", err)
		return authors, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var userId int64
		var username string
		var date time.Time
		var timezone int
		var email string
		rows.Scan(&id, &userId, &username, &date, &timezone, &email)
		u := &User{
			id:               userId,
			username:         username,
			registrationDate: date,
			timezone:         timezone,
			email:            email,
			db:               dbvendor,
		}

		a := Author{
			id:     id,
			userId: userId,
			user:   u,
			db:     dbvendor,
		}
		authors = append(authors, a)
	}

	return authors, nil
}

// Returns an author given its id.  If the id is not known to the blog
// a nil value is returned with an error.
func (conn *DBConnection) FindAuthorById(id int64) (*Author, error) {

	var a *Author
	var dbvendor = conn.databaser

	db, err := sql.Open(dbvendor.Driver(), dbvendor.Name())
	if err != nil {
		fmt.Println("FindAuthorById 1:", err)
		return a, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findAuthorById)
	if err != nil {
		fmt.Println("FindAuthorById 2:", err)
		return a, err
	}
	defer stmt.Close()

	var userId int64
	var username string
	var date time.Time
	var timezone int
	var email string

	err = stmt.QueryRow(id).Scan(&userId,
		&username,
		&date,
		&timezone,
		&email)

	if err != nil {
		// normal if the author doesnt exist
		return a, err
	}

	u := &User{
		id:               userId,
		username:         username,
		registrationDate: date,
		timezone:         timezone,
		email:            email,
		db:               dbvendor,
	}

	a = &Author{
		id:     id,
		userId: userId,
		user:   u,
		db:     dbvendor,
	}

	return a, nil
}

/*
*  Operations on Author
 */

// Save an author to the connence.  If the provided
// user didn't exist, it will create it first.
func (a *Author) Save() error {
	db, err := sql.Open(a.db.Driver(), a.db.Name())
	if err != nil {
		fmt.Println("Save 1:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertAuthorForId)
	if err != nil {
		fmt.Printf("Save 2 query=%s\n", insertAuthorForId)
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	err = a.user.Save()
	if err != nil {
		fmt.Println("Save 3:", err)
		return err
	}
	a.userId = a.user.Id()

	_, err = stmt.Exec(a.userId)
	if err != nil {
		fmt.Println("Save 4:", err)
		return err
	}

	// query the ID we inserted
	idStmt, err := db.Prepare(queryAuthorForUserId)
	if err != nil {
		fmt.Println("Save 5:", err)
		return err
	}
	defer idStmt.Close()

	row := idStmt.QueryRow(a.UserId())

	return row.Scan(&a.id)

}

// Removes the user from the author table.  The user attached to the author
// is not destroyed.
func (a *Author) Destroy() error {
	db, err := sql.Open(a.db.Driver(), a.db.Name())
	if err != nil {
		fmt.Println("Destroy 1:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteAuthorById)
	if err != nil {
		fmt.Println("Destroy 2:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.id)
	if err != nil {
		fmt.Println("Destroy 3:", err)
		return err
	}

	return nil
}
