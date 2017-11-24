// Eliza.go
// Main file - Main file of the Eliza chatbot implementation
// 		Implements a small webserver that serves a html template page and handles requests made to port 8080
//		Has functions that take userinput, transform it and puts a response back onto the html page
// Author - Bastian Graebener

package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bGraebener/Eliza/elizaHelper"
)

// struct that holds the data to be injected into the template
type data struct {
	Username string
	Greeting string
	NameSet  bool
}

var tmpData data
var t *template.Template

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// load resources from files
	elizaHelper.LoadResources()

	// parse the index html file
	t, _ = template.ParseFiles("index.html")

	// function that gets executed when the first request is made
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// set name to false so the name input form gets shown
		tmpData.NameSet = false
		t.Execute(w, &tmpData)
	})

	// start a new Eliza "session"
	http.HandleFunc("/session", newSessionHandler)

	// handle a request to /question that is being send when a question was submitted
	http.HandleFunc("/question", questionHandler)

	// serve the resource files
	http.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("res"))))

	// listen for requests on port 8080
	http.ListenAndServe(":8080", nil)

}

//function that starts a new "session" after the user put in a username
func newSessionHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the user name
	tmpData.Username = r.FormValue("userNameInput")

	// redirect to root page if user didn't enter a name
	if len(tmpData.Username) < 1 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// choose a random greeting from the greetings slice
	// elizaGreetings := elizaHelper.GetGreetings()
	ran := rand.Intn(len(elizaHelper.ElizaGreetings))
	tmpData.Greeting = elizaHelper.ElizaGreetings[ran]
	tmpData.NameSet = true

	// execute the html file
	t.Execute(w, &tmpData)
}

// function that gets executed every time a request is made to /question
func questionHandler(w http.ResponseWriter, r *http.Request) {
	// userFarewells := elizaHelper.GetUserFarewells()
	// elizaFarewells := elizaHelper.GetFarewells()

	// retrieve the header value for the field "user-question"
	question := r.Header.Get("user-question")
	question = strings.ToLower(question)

	// pass the username to be used in the chat window
	w.Header().Set("userName", tmpData.Username)

	// check if user quit the session
	for _, word := range strings.Split(question, " ") {
		if _, ok := elizaHelper.UserFarewells[word]; ok {
			// choose a random farewell phrase
			w.Header().Set("quit", "true")
			fmt.Fprintf(w, "%s", elizaHelper.ElizaFarewells[rand.Intn(len(elizaHelper.ElizaFarewells))])
			return
		}
	}

	// pass eliza response phrase to the responsewriter
	response := elizaHelper.GetResponse(question)
	fmt.Fprintf(w, "%s", response)

}
