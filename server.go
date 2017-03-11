package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

// index page

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
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
	http.Redirect(response, request, "/profile", 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/", 302)
}

// profile page

const profilePage = `
<h1>Profile</h1>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func profilePageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, profilePage)
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
