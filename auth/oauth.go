package auth

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/aybabtme/goblog/model"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

// TODO this shouldn't be hard coded here in the code, remember to reset it
// and put that in config files outside of the repo
var oauthCfg = &oauth.Config{
	ClientId:     "733675763142.apps.googleusercontent.com",
	ClientSecret: "dseJIDxz2ZlpYU6zn-BAMrYK",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
	RedirectURL:  "http://flying-unicorn.aybabt.me:5000/oauth2callback",
	Scope:        "https://www.googleapis.com/auth/userinfo.email profile",
}

var (
	code  = ""
	token = ""
)

func InteractiveOauth(conn *model.DBConnection, port string) {
	fmt.Println(`
Welcome to GoBlog!

In order to get started with your blog, we need to first create a user with
Author access! From this Author account, you will then be able to assign Author
access to other user accounts.

Let's get started! Please open the following URL in your browser:`)

	url := oauthCfg.AuthCodeURL("")

	fmt.Printf("\n%s\n", url)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Can't create listener")
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/oauth2callback",
		interactiveOAuth2Callback(conn, &l))

	fmt.Println(`
Go ahead while I wait here!  I'll carry on once I receive the callback from
Google and create your user.`)

	http.Serve(l, handler)

	authors, err := conn.FindAllAuthors()
	if err != nil {
		log.Println("Couldn't query for authors.", err)
		os.Exit(0)
		return
	}

	if len(authors) == 0 {
		log.Println("Received an empty list of authors.")
		os.Exit(0)
		return
	}

	fmt.Println("Got it!  You can now open the blog in your browser")

	for idx, author := range authors {
		fmt.Printf("\t Author #%d is %s\n", idx, author.User().Username())
	}

}

func interactiveOAuth2Callback(conn *model.DBConnection, lis *net.Listener) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//Get the code from the response
		w.Header().Set("Connection", "close")
		code := r.FormValue("code")
		errResp := r.FormValue("error")

		if "" != errResp {
			fmt.Println("Uhoh, you didn't give me access!\n" +
				"Sorry, I can't run without it!")
			os.Exit(0)
			return
		}

		t := &oauth.Transport{oauth.Config: oauthCfg}

		// Exchange the received code for a token
		t.Exchange(code)

		//now get user data based on the Transport which has the token
		resp, _ := t.Client().Get(profileInfoURL)
		dec := json.NewDecoder(resp.Body)

		var gUser map[string]interface{}
		if err := dec.Decode(&gUser); err != nil {
			log.Println("Couldn't decode json OAuth answer:", err)
			os.Exit(0)
			return
		}

		log.Println("Creating new user")
		user := createAuthUser(conn, gUser, t)
		author := conn.NewAuthor(user)
		if err := author.Save(); err != nil {
			log.Println("Coudln't save author from user!")
			os.Exit(0)
			return
		}
		(*lis).Close()
	}
}

func AuthorizeOauth(w http.ResponseWriter, r *http.Request) {
	//Get the Google URL which shows the Authentication page to the user
	url := oauthCfg.AuthCodeURL("")

	//redirect user to that page
	http.Redirect(w, r, url, http.StatusFound)
}

// Function that handles the callback from the Google server
func GetHandleOAuth2Callback(conn *model.DBConnection) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		//Get the code from the response
		code := r.FormValue("code")
		errResp := r.FormValue("error")

		if "" != errResp {
			http.Error(w, "Access to account was denied", http.StatusExpectationFailed)
			return
		}

		t := &oauth.Transport{oauth.Config: oauthCfg}

		// Exchange the received code for a token
		t.Exchange(code)

		//now get user data based on the Transport which has the token
		resp, _ := t.Client().Get(profileInfoURL)
		dec := json.NewDecoder(resp.Body)

		var gUser map[string]interface{}
		if err := dec.Decode(&gUser); err != nil {
			log.Println("Couldn't decode json OAuth answer:", err)
			http.Error(w, "Error decoding JSON answer.", http.StatusInternalServerError)
			return
		}

		if val, ok := gUser["verified_email"].(bool); !ok || !val {
			http.Error(w, "Need verified email", http.StatusNotAcceptable)
			return
		}

		// If user exists, retrieve it.
		user := recoverAuthUser(conn, t.ClientId)
		// Otherwise save create a new one and save it
		if user == nil {
			log.Println("Creating new user")
			user = createAuthUser(conn, gUser, t)
			if err := user.Save(); err != nil {
				log.Printf("Couldn't save new user <%v>\n", user)
				http.Error(w, "DB couldn't save user", http.StatusInternalServerError)
				return
			}
		}

		session, _ := store.Get(r, "user-session")
		session.Values["userId"] = strconv.FormatInt(user.Id(), 10)

		author, err := conn.FindAuthorById(user.Id())
		if err != nil {
			log.Printf("User <%s> is not an author", user.Username())
		} else if author != nil {
			log.Printf("Author id(%d) %v", author.Id(), author)
			authId := strconv.FormatInt(author.Id(), 10)
			session.Values["authorId"] = authId
			log.Printf("User <%s> is an author", user.Username())
		}

		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)

		userId, _ := session.Values["userId"].(string)
		authorId, _ := session.Values["authorId"].(string)
		log.Printf("userId=%s saved to cookie", userId)
		log.Printf("authorId=%s saved to cookie", authorId)

	}
}

func recoverAuthUser(conn *model.DBConnection,
	oauthId string) *model.User {
	user, err := conn.FindUserByOAuthId(oauthId)
	if err != nil {
		log.Printf("Couldn't find user with id <%v>\n", oauthId)
		log.Println(err)
		return nil
	}
	return user

}

func createAuthUser(conn *model.DBConnection,
	gUser map[string]interface{},
	t *oauth.Transport) *model.User {

	username, ok := gUser["name"].(string)
	if !ok {
		log.Println("HandleOAuth2: username.")
		return nil
	}
	email, ok := gUser["email"].(string)
	if !ok {
		log.Println("HandleOAuth2: email.")
		return nil
	}

	user := conn.NewUser(username,
		time.Now().UTC(),
		-5,
		t.ClientId,
		t.AccessToken,
		t.RefreshToken,
		email)

	return user
}
