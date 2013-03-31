package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/aybabtme/goblog/ctlr"
	"github.com/aybabtme/goblog/model"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Router string

func (r Router) Start(port string, conn *model.DBConnection) error {

	controllers := []ctlr.Controller{
		ctlr.NewIndexController(),
		ctlr.NewAuthorController(),
		ctlr.NewUserController(),
		ctlr.NewLabelController(),
		ctlr.NewPostController(),
		ctlr.NewPostComposeController(),
		ctlr.NewPostSaveController(),
		ctlr.NewPostDestroyController(),
		ctlr.NewPostIdController()}

	muxer := mux.NewRouter()
	for _, ctlr := range controllers {
		muxer.HandleFunc(ctlr.Path(), ctlr.Controller(conn))
	}
	// serve dynamic resources
	http.Handle("/", muxer)
	// serve static resources
	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir("public/"))))

	// For user authentication
	http.HandleFunc("/authorize", handleAuthorize)
	http.HandleFunc("/oauth2callback", handleOAuth2Callback)

	return http.ListenAndServe(":"+port, nil)
}

// variables used during oauth protocol flow of authentication
var (
	code  = ""
	token = ""
)

//This is the URL that Google has defined so that an authenticated application may obtain the user's info in json format
const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

// Start the authorization process
func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	//Get the Google URL which shows the Authentication page to the user
	url := oauthCfg.AuthCodeURL("")

	//redirect user to that page
	http.Redirect(w, r, url, http.StatusFound)
}

var userInfoTemplate = template.Must(template.New("").Parse(`
<html><body>
This app is now authenticated to access your Google user info.  Your details are:<br />
{{.}}
</body></html>
`))

// Function that handles the callback from the Google server
func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	//Get the code from the response
	code := r.FormValue("code")
	err := r.FormValue("error")

	if "" != err {
		fmt.Fprintf(w, "Access denied.")
		return
	}

	t := &oauth.Transport{oauth.Config: oauthCfg}

	// Exchange the received code for a token
	t.Exchange(code)

	//now get user data based on the Transport which has the token
	resp, _ := t.Client().Get(profileInfoURL)

	buf := make([]byte, 1024)
	resp.Body.Read(buf)
	userInfoTemplate.Execute(w, string(buf))
}
