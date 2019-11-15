// Eliza.go
// Main file - Main file of the Eliza chatbot implementation
// 		Implements a small webserver that serves a html template page and handles requests made to port 8080
//		Has functions that handles requests when the user first visits the index page, when he/she entered a name and when a user question was submitted
// Author - Bastian Graebener

package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"encoding/json"

	"io/ioutil"
	"log"

	"regexp"
	"sort"


)

// struct that holds the data to be injected into the template
type data struct {
	Username string
	Greeting string
	NameSet  bool // determines which part of the hmtl document gets displayed
}

var tmpData data
var t *template.Template

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// load resources from files
	LoadResources()

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
	ran := rand.Intn(len(ElizaGreetings))
	tmpData.Greeting = ElizaGreetings[ran]
	tmpData.NameSet = true

	// execute the html file
	t.Execute(w, &tmpData)
}

// function that gets executed every time a request is made to /question
func questionHandler(w http.ResponseWriter, r *http.Request) {

	// retrieve the header value for the field "user-question"
	question := r.Header.Get("user-question")
	question = strings.ToLower(question)

	// pass the username to be used in the chat window
	w.Header().Set("userName", tmpData.Username)

	// check if user quit the session
	for _, word := range strings.Split(question, " ") {
		if _, ok := UserFarewells[word]; ok {
			// choose a random farewell phrase
			w.Header().Set("quit", "true")
			fmt.Fprintf(w, "%s", ElizaFarewells[rand.Intn(len(ElizaFarewells))])
			return
		}
	}

	// pass eliza response phrase to the responsewriter
	response := GetResponse(question)
	fmt.Fprintf(w, "%s", response)

}



// variable that holds all keywords contained in keywords.json
var elizaData keyWords

// map that holds pronouns and their substitutions
var substitutions map[string]string

// ElizaGreetings holds strings for greeting the user
var ElizaGreetings []string

// ElizaFarewells holds responses for after the user quit the program
var ElizaFarewells []string

// UserFarewells holds string options for the user to quit the program
var UserFarewells map[string]int

//KeyWord holds the keywords, their assoctiated rank and a Decomp struct
type KeyWord struct {
	Keyword string `json:"keyword" `
	Rank    int    `json:"rank"`
	Decomp  []Decomp
}

// Decomp holds the regular expression pattern for decomposition and an array of associated responses
type Decomp struct {
	DisAssRule string   `json:"rule"`
	Responses  []string `json:"responses"`
}

// keyWords redefines a slice of KeyWords
type keyWords []KeyWord

// implementing the sort interface functions for the type keyWords
func (r keyWords) Len() int { return len(r) }
func (r keyWords) Less(r1, r2 int) bool {
	return r[r1].Rank < r[r2].Rank
}
func (r keyWords) Swap(r1, r2 int) {
	r[r1], r[r2] = r[r2], r[r1]
}

// GetResponse returns an appropriate response to the user input
func GetResponse(userInput string) string {

	keyWordList := getKeyWordList(userInput)

	return generateResponse(userInput, keyWordList)
}

// LoadResources loads the resources from the json files
func LoadResources() {

	// load all keyword data into memory
	loadKeywordData()

	// load the substitutions from file
	loadSubstitutions()

	// load the greetings and farewells
	loadGreetings()
}

// loadSubstitutions reads the substitutions file and populates the map of substitutions
func loadSubstitutions() {
	// read the json file
	if raw, err := ioutil.ReadFile("./res/substitutions.json"); err != nil {
		log.Fatal("Couldn't read substitutions.json")
	} else {
		//parse the json file
		if err := json.Unmarshal(raw, &substitutions); err != nil {
			log.Fatal("Couldn't parse substitutions.json")
		}
	}
}

// loadGreetings reads the greetings from the startEnd.json file
func loadGreetings() {
	dataMap := make(map[string][]string)

	// read the json file
	if raw, err := ioutil.ReadFile("./res/startEnd.json"); err != nil {
		log.Fatal("Couldn't read startEnd.json!")
	} else {

		// parse the json data
		if err := json.Unmarshal(raw, &dataMap); err != nil {
			log.Fatal("Couldn't parse startEnd.json")
		}
	}

	// populate the global phrase slices
	ElizaGreetings = dataMap["elizaGreetings"]
	ElizaFarewells = dataMap["elizaFarewells"]
	UserFarewells = SliceToMap(dataMap["userFarewells"])
}

//parses the keyword data from the keyword.json file
func loadKeywordData() {
	// attempt to read the file
	if raw, err := ioutil.ReadFile("./res/keywords.json"); err != nil {
		log.Fatal("Couldn't read keywords.json!")
	} else {
		// parse the json data into the special struct slice
		if err = json.Unmarshal(raw, &elizaData); err != nil {
			log.Fatal("Couldn't parse keywords.json!")
		}
	}
}

// splits the user input string into a string slice
// and finds all keywords contained in the user string
func getKeyWordList(userInput string) (keyWordList keyWords) {

	// replace all non-letter characters with a whitespace
	reg := regexp.MustCompile("[^a-zA-Z]+")
	userInput = reg.ReplaceAllString(userInput, " ")

	// turn phrase string into string slice of individual words
	userWords := strings.Split(userInput, " ")
	// find all keywords contained in user string and store them in a slice of keywords
	for _, word := range userWords {
		for i := range elizaData {
			if elizaData[i].Keyword == word {
				keyWordList = append(keyWordList, elizaData[i])
				break
			}
		}
	}
	// sort the keyword list so highest ranking keyword is first
	sort.Sort(sort.Reverse(keyWordList))
	return keyWordList
}

// searches every found keyword and finds the best matching substring as per the keywords decomposition pattern
// returns a random response from the pool of responses for the specific decomposition pattern of the keyword
func generateResponse(userInput string, keyWordList []KeyWord) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var response string

	// iterate over the all the keywords found in the user input string
	for _, keyWord := range keyWordList {

		// for every keyword iterate over all the decomposition patterns
		for _, decomp := range keyWord.Decomp {

			// compile the decomposition pattern into a regular expression
			reg := regexp.MustCompile(decomp.DisAssRule)

			// check if the decomposition pattern is found in the user question and
			// save the capture group values
			captureGroup := reg.FindStringSubmatch(userInput)
			// no matching substring found for this decomposition pattern
			if len(captureGroup) == 0 {
				continue
			}

			// choose a random response
			response = decomp.Responses[rand.Intn(len(decomp.Responses))]

			// disregard regex capture group if the response doesn't need it
			// or the capture group only contained the whole string
			if !strings.Contains(response, "%s") || len(captureGroup) == 1 {
				return response
			}

			// assemble a response with a randomly chosen response and the captured value
			// from the user input
			if len(captureGroup) > 1 && strings.Contains(userInput, captureGroup[0]) {

				// reflect any pronouns like "my", "yours"
				captureGroupValue := substitute(captureGroup[1])

				// reassamble the response and return it
				return fmt.Sprintf(response, captureGroupValue)
			}
		}
	}

	// if no keyword matches fall back to a random response from the "xnone" keyword
	return elizaData[0].Decomp[0].Responses[rand.Intn(len(elizaData[0].Decomp[0].Responses))]
}

// Substitute takes a string, checks it against the pronouns map and substitutes any found pronouns for their counterpart
func substitute(captureGroupValue string) string {
	// get individual words
	words := strings.Split(captureGroupValue, " ")
	// iterate over every word and if the pronouns map contains it, switch it
	for i, word := range words {
		if _, ok := substitutions[word]; ok {
			words[i] = substitutions[word]
		}
	}
	// reassemble the string and return it
	return strings.Join(words, " ")
}

// SliceToMap converts a string slice into a map
// convience function for fastercheck if user entered a keyword for quiting the program
func SliceToMap(slice []string) map[string]int {

	tmpMap := make(map[string]int)

	for _, i := range slice {
		tmpMap[i]++
	}
	return tmpMap
}