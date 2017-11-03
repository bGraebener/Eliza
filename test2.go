package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

//KeyWord holds the keywords, their assoctiated rank and a Decomp struct
type KeyWord struct {
	Keyword string `json:"keyword" `
	Rank    int    `json:"rank"`
	Decomp  []Decomp
}

// Decomp holds the decomposition regular expression pattern and an array of associated responses
type Decomp struct {
	DisAssRule string   `json:"value"`
	Responses  []string `json:"responses"`
}

// keyWords redefines a slice of KeyWords
type keyWords []KeyWord

// defining functions so keyWords implements the sort interface
func (r keyWords) Len() int { return len(r) }
func (r keyWords) Less(r1, r2 int) bool {
	return r[r1].Rank < r[r2].Rank
}
func (r keyWords) Swap(r1, r2 int) {
	r[r1], r[r2] = r[r2], r[r1]
}

// global variable that holds all keywords contained in keywords.json
var elizaData keyWords

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	elizaData = readKeywordData()
	userInput := getUserInput()
	response := getResponse(userInput)
	fmt.Println(response)
}

func getResponse(userInput string) string {

	keyWordList := getKeyWordList(userInput)

	return findDecompPatterns(userInput, keyWordList)

}

//parse the keyword data from the keyword.json file
func readKeywordData() keyWords {
	var list keyWords
	// attempt to read the file
	if raw, err := ioutil.ReadFile("keywords.json"); err != nil {
		log.Fatal(err)
	} else {
		// parse the json data into the special struct slice
		if err = json.Unmarshal(raw, &list); err != nil {
			fmt.Println(err)
		}
	}
	return list
}

func getUserInput() string {
	fmt.Print("Please enter a phrase: ")
	var userInput string
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	userInput = strings.ToLower(reader.Text())

	return userInput
}

// splits the user input string into a string slice
// and finds all keywords contained in the user string
func getKeyWordList(userInput string) (keyWordList keyWords) {

	// replace all non-letter characters with a whitespace
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}
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
	// sort the keyword list by descending rank
	sort.Sort(sort.Reverse(keyWordList))
	return keyWordList
}

// searches every found keyword and finds the best matching substring as per the keywords decomposition pattern
// returns a random response from the pool of responses for the specific decomposition pattern of the keyword
func findDecompPatterns(userInput string, keyWordList []KeyWord) string {
	// iterate over the all the keywords found in the user input string
	for _, keyWord := range keyWordList {
		// for every keyword iterate over all the decomposition patterns
		for _, decomp := range keyWord.Decomp {
			fmt.Println(decomp.DisAssRule)
			// compile them into a regular expression
			reg, err := regexp.Compile(decomp.DisAssRule)
			if err != nil {
				log.Fatal(err)
			}

			// if the decomposition pattern is found in the user input string return a random response from the set of responses for the decomposition pattern
			if reg.MatchString(userInput) {
				return decomp.Responses[rand.Intn(len(decomp.Responses))]
			}

		}
	}

	// if no keyword matches or no matching decomposition pattern is found fall back to a random response from the "xnone" keyword
	// fmt.Println(elizaData[0].Decomp[0].Responses[rand.Intn(len(elizaData[0].Decomp[0].Responses))])
	return elizaData[0].Decomp[0].Responses[rand.Intn(len(elizaData[0].Decomp[0].Responses))]

}
