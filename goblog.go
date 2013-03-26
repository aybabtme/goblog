package main

import (
	"fmt"
	"github.com/aybabtme/goblog/model"
	"github.com/aybabtme/gypsum"
	"os"
	"time"
)

func main() {

	modelurl := os.Getenv("DATABASE_URL")
	if modelurl == "" {
		fmt.Println("Need a database to connect to!\n" +
			"export DATABASE_URL=<your model url here>")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("No port specified.\n" +
			"export PORT=<port number here>")
		return
	}

	conn, err := setupDatabase(modelurl)
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

func setupDatabase(modelurl string) (*model.DBConnection, error) {
	postgres := model.NewPostgreser(modelurl)
	conn, err := model.NewConnection(postgres)
	if err != nil {
		return nil, err
	}
	conn.DeleteConnection()
	return model.NewConnection(postgres)
}

func generateData(conn *model.DBConnection) error {

	user1 := conn.NewUser(
		"Antoine Grondin",
		time.Now().UTC(),
		-5,
		"google+",
		"heheveb7673tygvh23",
		"antoinegrondin@gmail.com")
	author := conn.NewAuthor(user1)
	if err := author.Save(); nil != err {
		return err
	}
	post1 := conn.NewPost(author,
		"Il était une fois",
		gypsum.Lorem(),
		"http://media.zoom-cinema.fr/photos/news/2380/il-etait-une-fois-2007-4.jpg",
		time.Now().UTC())
	if err := post1.Save(); nil != err {
		return err
	}
	if _, err := post1.AddLabel("prince"); err != nil {
		return err
	}
	if _, err := post1.AddLabel("princess"); err != nil {
		return err
	}
	if _, err := post1.AddLabel("drama"); err != nil {
		return err
	}

	post2 := conn.NewPost(author,
		"Parenthood",
		gypsum.Lorem(),
		"http://www.blessedquietness.com/STUPID01.jpg",
		time.Now().UTC())
	if err := post2.Save(); nil != err {
		return err
	}

	if _, err := post2.AddLabel("childhood"); err != nil {
		return err
	}
	if _, err := post2.AddLabel("funny"); err != nil {
		return err
	}

	user2 := conn.NewUser(
		"John Smith",
		time.Now().UTC(),
		-5,
		"google+",
		"fdvh23",
		"ajfdsin@gmail.com")
	user2.Save()
	user3 := conn.NewUser(
		"Chris Poirier",
		time.Now().UTC(),
		-5,
		"google+",
		"fdvh23",
		"ajf21vbnlkdsin@gmail.com")
	user3.Save()

	conn.NewComment(user1.Id(), post1.Id(), gypsum.Lorem(), time.Now()).Save()
	conn.NewComment(user2.Id(), post1.Id(), gypsum.Lorem(), time.Now()).Save()
	conn.NewComment(user3.Id(), post1.Id(), gypsum.Lorem(), time.Now()).Save()
	conn.NewComment(user1.Id(), post2.Id(), gypsum.Lorem(), time.Now()).Save()
	conn.NewComment(user2.Id(), post2.Id(), gypsum.Lorem(), time.Now()).Save()
	conn.NewComment(user3.Id(), post2.Id(), gypsum.Lorem(), time.Now()).Save()
	conn.NewComment(user1.Id(), post2.Id(), gypsum.Lorem(), time.Now()).Save()
	conn.NewComment(user2.Id(), post2.Id(), gypsum.Lorem(), time.Now()).Save()
	conn.NewComment(user3.Id(), post2.Id(), gypsum.Lorem(), time.Now()).Save()

	return nil
}
