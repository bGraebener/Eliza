package main

import (
	"fmt"
	"net/http"
	"strings"
)

// function that gets executed every time a request is made to /question
func questionHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the header value for the field "user-question"
	question := r.Header.Get("user-question")

	question = strings.ToUpper(question)

	if len(question) > 0 {
		fmt.Fprintf(w, "%s", question)
	}
}

func main() {

	// handle a request to /question
	http.HandleFunc("/question", questionHandler)

	// handle a request to the root path
	http.Handle("/", http.FileServer(http.Dir("./res")))

	// listen for requests on port 8080
	http.ListenAndServe(":8080", nil)

}
