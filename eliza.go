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
	Quit     bool
}

// arrays that hold strings for greeting the user and responses after the user quit the program
var elizaGreetings []string
var elizaFarewells []string

// options for the user to quit the program
var userFarewells map[string]int

var tmpData data
var t *template.Template

// function that gets executed every time a request is made to /question
func questionHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the header value for the field "user-question"
	question := r.Header.Get("user-question")
	question = strings.ToLower(question)

	// check if user quit the session
	for _, word := range strings.Split(question, " ") {
		if _, ok := userFarewells[word]; ok {
			// choose a random farewell phrase
			w.Header().Set("quit", "true")
			fmt.Fprintf(w, "%s", elizaFarewells[rand.Intn(len(elizaFarewells))])
			return
		}
	}

	// pass the username to be used in the chat window
	w.Header().Set("userName", tmpData.Username)

	// pass eliza response phrase to the responsewriter
	response := elizaHelper.GetResponse(question)
	fmt.Fprintf(w, "%s", response)

}

//function that starts a new "session" after the user put in a username
func newSessionHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the user name
	tmpData.Username = r.FormValue("userNameInput")

	// redirect to same page if user didn't enter a name
	if len(tmpData.Username) < 1 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	// choose a random greeting from the greetings slice
	ran := rand.Intn(len(elizaGreetings))
	tmpData.Greeting = elizaGreetings[ran]

	// execute the html file
	t.Execute(w, &tmpData)
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	// load resources from resources file
	dataMap := elizaHelper.LoadResources()

	// split individual string data in the correct containers
	elizaGreetings = dataMap["elizaGreetings"]
	elizaFarewells = dataMap["elizaFarewells"]
	userFarewells = elizaHelper.SliceToMap(dataMap["userFarewells"])
	// fmt.Println(farewells)

	// parse the session html file
	t, _ = template.ParseFiles("./res/session.html")

	// handle a request to /question
	http.HandleFunc("/question", questionHandler)

	// start a new "session"
	http.HandleFunc("/session", newSessionHandler)

	// handle a request to the root path
	http.Handle("/", http.FileServer(http.Dir("./res")))

	// listen for requests on port 8080
	http.ListenAndServe(":8080", nil)

}
