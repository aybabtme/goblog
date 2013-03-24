package main

import (
	"fmt"
	"github.com/aybabtme/goblog/db"
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

	conn, err := setupDatabase(dburl)
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
	conn, err := db.NewConnection(postgres)
	if err != nil {
		return nil, err
	}
	conn.DeleteConnection()
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
	author := conn.NewAuthor(user)
	if err := author.Save(); nil != err {
		return err
	}
	post1 := conn.NewPost(author.Id(),
		"Il était une fois",
		"Lorem ipsum shit chien vache",
		"http://media.zoom-cinema.fr/photos/news/2380/il-etait-une-fois-2007-4.jpg",
		time.Now().UTC())
	if err := post1.Save(); nil != err {
		return err
	}

	post2 := conn.NewPost(author.Id(),
		"Grosse Truie avec un Gros Cul",
		"XXX gratis, donne nous juste ton carte de crédit pis on te promet de pas l'utiliser",
		"http://www.blacktowhite.net/wp-content/uploads/2011/05/cock-sucking-bitches-05-590x398.jpg",
		time.Now().UTC())
	if err := post2.Save(); nil != err {
		return err
	}
	return nil
}
