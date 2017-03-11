package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

const indexPage = `
hello world
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

func main() {
	var router = mux.NewRouter()

	router.HandleFunc("/", indexPageHandler)
	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
}
