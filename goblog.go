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

	postgres := db.NewPostgreser(dburl)
	conn, err := db.NewConnection(postgres)
	if err != nil {
		fmt.Println("Couldn't connect to database.")
		panic(err)
	}

	err = generateData(conn)
	if err != nil {
		fmt.Println("Couldn't generate data")
	}

	http.HandleFunc("/", indexHandler(conn))
	fmt.Println("listening...")
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func indexHandler(conn *db.DBConnection) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "<h1>hello, world</h1>")
	}
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
