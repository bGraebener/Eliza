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

	// pass the username to be used in the chat window
	w.Header().Set("userName", tmp.Username)

	question = strings.ToUpper(question)

	if len(question) > 0 {
		fmt.Fprintf(w, "%s", question)
	}
}

//function that starts a new "session" after the user put in a username
func sessionHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the user name
	tmp.Username = r.FormValue("userNameInput")

	// redirect to same page if user didn't enter a name
	if len(tmp.Username) < 1 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// parse and execute the main html file
	t, _ := template.ParseFiles("./res/session.html")
	t.Execute(w, &tmp)
}

func main() {

	// handle a request to /question
	http.HandleFunc("/question", questionHandler)

	// start a new "session"
	http.HandleFunc("/session", sessionHandler)

	// handle a request to the root path
	http.Handle("/", http.FileServer(http.Dir("./res")))

	// listen for requests on port 8080
	http.ListenAndServe(":8080", nil)

}
