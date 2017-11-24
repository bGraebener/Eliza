// Helper.go
// Utility file - Contains utility functions and custom data types for the eliza chatbot
// 				Keyword struct: represents a keyword. It contains the keyword itself, its rank and a slice of decomposition structs.
// 				Decomp struct: represents a decomposition rule. It contains a dissassembly rule and a slice of response strings.
// 				keywords type: a redefinition of a keyword slice. Implements the Sort interface to sort the keywords by rank.
// Author - Bastian Graebener

package elizaHelper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

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

// SliceToMap converts a string slice into a map,
// convience function for faster, easier lookup of keywords and responses
func SliceToMap(slice []string) map[string]int {

	tmpMap := make(map[string]int)

	for _, i := range slice {
		tmpMap[i]++
	}
	return tmpMap
}
