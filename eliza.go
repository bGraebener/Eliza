package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type data struct {
	Username string
}

var tmp data

// function that gets executed every time a request is made to /question
func questionHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the header value for the field "user-question"
	question := r.Header.Get("user-question")

	w.Header().Set("userName", tmp.Username)

	question = strings.ToUpper(question)

	if len(question) > 0 {
		fmt.Fprintf(w, "%s", question)
	}
}

func newSessionHandler(w http.ResponseWriter, r *http.Request) {
	tmp.Username = r.FormValue("userNameInput")
	t, _ := template.ParseFiles("./res/session.html")
	t.Execute(w, &tmp)
}

func main() {

	// handle a request to /question
	http.HandleFunc("/question", questionHandler)

	http.HandleFunc("/newSession", newSessionHandler)

	// handle a request to the root path
	http.Handle("/", http.FileServer(http.Dir("./res")))

	// listen for requests on port 8080
	http.ListenAndServe(":8080", nil)

}
