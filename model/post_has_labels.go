package model

import (
	"database/sql"
	"fmt"
	"time"
)

/*
 * The SQL String Dance Starts Here
 * --------------------------------
 * Relation stuff
 */

var createLabelPostsRelation string = `
CREATE TABLE IF NOT EXISTS LabelPost(
   post_id INTEGER,
   label_id INTEGER,
   PRIMARY KEY (post_id, label_id),
   CONSTRAINT fk_labelpost_post_id
      FOREIGN KEY (post_id) REFERENCES Post(post_id) ON DELETE CASCADE,
   CONSTRAINT fk_labelpost_label_id
      FOREIGN KEY (label_id) REFERENCES Label(label_id) ON DELETE CASCADE
);`

var dropLabelPostsRelation string = `
DROP TABLE LabelPost;`

// used
var insertLabelPostRelation string = `
INSERT INTO LabelPost( post_id, label_id )
VALUES( $1, $2 )`

var findPostsByLabelId string = `
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
	LabelPost AS LP,
	Author AS A,
	BlogUser AS U
WHERE
	LP.label_id = $1
	AND LP.post_id = P.post_id
	AND P.author_id = A.author_id
	AND A.author_id = U.user_id`

// used
var findLabelsByPostId string = `
SELECT L.label_id, L.name
FROM Label AS L, LabelPost AS LP
WHERE LP.post_id = $1 AND LP.label_id = L.label_id`

var deleteAllLabelWithIdFromRelation string = `
DELETE FROM LabelPost
WHERE LabelPost.label_id = $1`

var deleteAllLabelWithIdFromTable string = `
DELETE FROM Label
WHERE Label.label_id = $1;`

// used
var deleteLabelFromPostId string = `
DELETE FROM LabelPost
WHERE LabelPost.post_id = $1`

/*
 * Labels stuff
 */

var insertLabelForId string = `
INSERT INTO Label( name )
VALUES( $1 )`

var queryLabelForName string = `
SELECT L.label_id, L.name
FROM Label AS L
WHERE L.name = $1`

/*
 * No More SQL Strings
 * -------------------
 */

func openDatabase(d *DBVendor) (*sql.DB, error) {
	return sql.Open((*d).Driver(), (*d).Name())
}

func (pers *DBConnection) createLabelPostRelation() {
	var vendor = pers.databaser
	db, err := openDatabase(&vendor)
	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(createLabelPostsRelation)
	if err != nil {
		fmt.Printf("Error creating LabelPost relation, driver \"%s\","+
			" modelname \"%s\", query = \"%s\"\n",
			vendor.Driver(), vendor.Name(), createLabelPostsRelation)
		fmt.Println(err)
		return
	}
}

func (pers *DBConnection) dropLabelPostRelation() {
	var vendor = pers.databaser
	db, err := openDatabase(&vendor)

	if err != nil {
		fmt.Println("Error on open of database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(dropLabelPostsRelation)
	if err != nil {
		fmt.Println("Error droping table:", err)
	}
}

/*
 * Stuff that can be done using a Post
 */

func (p *Post) AddLabel(name string) (Label, error) {
	// Label is a weak entity and thus can't exist outside
	// of a relationship with a Post.  This is enforced by
	// the SQL of Label, which states that it cannot exist
	// without the LabelPosts relation existing itself.
	// Thus we need to create a label manually and put it
	// in the relationship in a single transaction to avoid
	// integrity restriction problems.
	var lbl = Label{
		id:   -1,
		name: name,
		conn: p.conn,
	}

	db, err := openDatabase(&p.conn.databaser)
	if err != nil {
		return lbl, err
	}
	defer db.Close()

	// Start the transaction
	tx, err := db.Begin()
	if err != nil {
		return lbl, err
	}

	// Create the Label
	lblStmt, err := tx.Prepare(insertLabelForId)
	if err != nil {
		fmt.Println("AddLabel 1. Can't create stmt: ", err)
		return tryRollback(lbl, tx, err)
	}
	defer lblStmt.Close()

	_, err = lblStmt.Exec(lbl.Name())
	if err != nil {
		// error may mean it already exists
		// need to restart transaction
		_, _ = tryRollback(lbl, tx, err)
		tx, err = db.Begin()
		if err != nil {
			return lbl, err
		}
	}

	lblFindBack, err := tx.Prepare(queryLabelForName)
	if err != nil {
		fmt.Println("Add Label 3. Can't create stmt: ", err)
		return tryRollback(lbl, tx, err)
	}
	defer lblFindBack.Close()

	err = lblFindBack.QueryRow(lbl.name).Scan(&lbl.id, &lbl.name)
	if err != nil {
		fmt.Println("Add Label 4. Can't query id: ", err)
		return tryRollback(lbl, tx, err)
	}

	// Then establish the relationship
	relStmt, err := tx.Prepare(insertLabelPostRelation)
	if err != nil {
		fmt.Println("Add Label 5. Can't create stmt: ", err)
		return tryRollback(lbl, tx, err)
	}
	defer relStmt.Close()

	_, err = relStmt.Exec(p.Id(), lbl.Id())
	if err != nil {
		fmt.Println("Add Label 6. Can't query ids: ", err)
		return tryRollback(lbl, tx, err)
	}

	// All set, try committing
	err = tx.Commit()
	if err != nil {
		fmt.Println("Add Label 7. Can't commit: ", err)
		return tryRollback(lbl, tx, err)
	}

	return lbl, nil
}

func tryRollback(lbl Label, tx *sql.Tx, err error) (Label, error) {
	rbErr := tx.Rollback()
	if rbErr != nil {
		return lbl, rbErr
	}
	return lbl, err
}

// Removes a label from a post.  If the post if the only post
// referring to that label, it will delete the label altogether.
// Otherwise it will remove the label only for that post, leaving
// other posts unaffected
func (p *Post) RemoveLabel(label *Label) error {
	db, err := openDatabase(&p.conn.databaser)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteLabelFromPostId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Id())
	return err

}

// Returns all the post associated with this post, if any.
func (p *Post) Labels() ([]Label, error) {
	// I prefer returning an empty list than a nil pointer
	var labels []Label

	db, err := openDatabase(&p.conn.databaser)
	if err != nil {
		fmt.Println("PostLabels 1:", err)
		return labels, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findLabelsByPostId)
	if err != nil {
		fmt.Println("PostLabels 2:", err)
		return labels, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(p.Id())
	if err != nil {
		fmt.Println("PostLabels 3:", err)
		return labels, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			fmt.Println("PostLabels 4:", err)
			return labels, err
		}
		var aLabel = Label{
			id:   id,
			name: name,
		}
		labels = append(labels, aLabel)
	}
	return labels, nil
}

/*
 * Stuff that can be done using a Label
 */

// Deletes the label from the database.  If any post is referencing this
// label, they will not do so anymore
func (l *Label) Destroy() error {
	vendor := l.conn.databaser
	db, err := openDatabase(&vendor)
	if err != nil {
		return err
	}
	defer db.Close()

	stmtRelation, err := db.Prepare(deleteAllLabelWithIdFromRelation)
	if err != nil {
		return err
	}
	defer stmtRelation.Close()

	_, err = stmtRelation.Exec(l.Id())
	if err != nil {
		return err
	}

	stmtLabel, err := db.Prepare(deleteAllLabelWithIdFromTable)
	if err != nil {
		return nil
	}
	defer stmtLabel.Close()
	_, err = stmtLabel.Exec(l.Id())

	return err
}

// Returns all the posts making reference to this label.
func (l *Label) Posts() ([]Post, error) {
	var posts []Post

	vendor := l.conn.databaser

	db, err := openDatabase(&vendor)
	if err != nil {
		return posts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findPostsByLabelId)
	if err != nil {
		return posts, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(l.Id())
	if err != nil {
		return posts, err
	}

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
			conn:             l.conn,
		}

		a := &Author{
			id:   authorId,
			user: u,
			conn: l.conn,
		}

		p := Post{
			id:       id,
			author:   a,
			title:    title,
			content:  content,
			imageURL: imageURL,
			date:     date,
			conn:     l.conn,
		}
		posts = append(posts, p)
	}

	return posts, nil
}
