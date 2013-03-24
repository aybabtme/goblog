package main

import (
	"fmt"
	"github.com/aybabtme/goblog/db"
	"net/http"
	"os"
	"time"
)

func main() {

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		fmt.Println("Need a database to connect to!\n" +
			"export DATABASE_URL=<your db url here>")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("No port specified.\n" +
			"export PORT=<port number here>")
		return
	}

	conn, err := db.NewConnection(postgres)
	if err != nil {
		fmt.Println("Couldn't connect to database.")
		panic(err)
	}

	err = generateData(conn)
	if err != nil {
		fmt.Println("Couldn't generate data")
	}

	var r Router
	if err := r.Start(port, conn); err != nil {
		panic(err)
	}
}

func setupDatabase(dburl string) (*db.DBConnection, error) {
	postgres := db.NewPostgreser(dburl)
	return db.NewConnection(postgres)
}

func generateData(conn *db.DBConnection) error {
	user := conn.NewUser(
		"antoine",
		time.Now().UTC(),
		-5,
		"google+",
		"heheveb7673tygvh23",
		"antoinegrondin@gmail.com")
	if err := user.Save(); nil != err {
		return err
	}
	author := conn.NewAuthor(user)
	if err := author.Save(); nil != err {
		return err
	}
	post := conn.NewPost(author.Id(),
		"Il était une fois",
		"Lorem ipsum shit chien vache",
		"/path/to/image.png",
		time.Now().UTC())
	if err := post.Save(); nil != err {
		return err
	}
	return nil
}
