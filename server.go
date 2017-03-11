package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64), // hash key at least 32 bytes
	securecookie.GenerateRandomKey(32)) // block key at least 16 bytes

func getUsername(request *http.Request) string {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			return cookieValue["username"]
		}
	}
	return ""
}

func createSession(username string, response http.ResponseWriter) {
	value := map[string]string{
		"username": username,
	}

	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func destroySession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(response, cookie)
}

// index page

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="username">User name</label>
    <input type="text" id="username" name="username">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

// index handler

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	password := request.FormValue("password")
	redirectURL := "/"

	if username != "" && password != "" {
		// TODO check username and password
		createSession(username, response)
		redirectURL = "/profile"
	}

	http.Redirect(response, request, redirectURL, 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	destroySession(response)
	http.Redirect(response, request, "/", 302)
}

// profile page

const profilePage = `
<h1>Profile</h1>
<small>Username: %s</small>

<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func profilePageHandler(response http.ResponseWriter, request *http.Request) {
	username := getUsername(request)
	if username != "" {
		fmt.Fprintf(response, profilePage, username)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func main() {
	var router = mux.NewRouter()

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/profile", profilePageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
}
