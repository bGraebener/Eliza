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

//KeyWord holds the keywords, their assoctiated rank and a Decomp struct
type KeyWord struct {
	Keyword string `json:"keyword" `
	Rank    int    `json:"rank"`
	Decomp  []Decomp
}

// Decomp holds the decomposition regular expression pattern and an array of associated responses
type Decomp struct {
	DisAssRule string   `json:"rule"`
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

// map of words to be substituted
var substitutions = map[string]string{
	"am":     "are",
	"was":    "were",
	"i":      "you",
	"i'd":    "you would",
	"i've":   "you have",
	"i'll":   "you will",
	"my":     "your",
	"are":    "am",
	"you've": "I have",
	"you'll": "I will",
	"your":   "my",
	"yours":  "mine",
	"you":    "me",
	"me":     "you",
}

// LoadResources loads the resources from the json file
func LoadResources() map[string][]string {

	dataMap := make(map[string][]string)

	// read the json file
	raw, err := ioutil.ReadFile("./res/startEnd.json")
	if err != nil {
		panic("Couldn't read resource file!")
	}

	// parse the json data
	if err := json.Unmarshal(raw, &dataMap); err != nil {
		panic("Couldn't parse json file")
	}

	elizaData = readKeywordData()
	return dataMap

}

// SliceToMap converts a string slice into a map, convience function for faster, easier lookup of keywords and responses
func SliceToMap(slice []string) map[string]int {

	tmpMap := make(map[string]int)

	for _, i := range slice {
		tmpMap[i]++
	}
	return tmpMap
}

// GetResponse returns an appropriate response to the user input
func GetResponse(userInput string) string {

	keyWordList := getKeyWordList(userInput)

	return findDecompPatterns(userInput, keyWordList)

}

//parse the keyword data from the keyword.json file
func readKeywordData() keyWords {
	var list keyWords
	// attempt to read the file
	if raw, err := ioutil.ReadFile("./res/keywords.json"); err != nil {
		log.Fatal(err)
	} else {
		// parse the json data into the special struct slice
		if err = json.Unmarshal(raw, &list); err != nil {
			fmt.Println(err)
		}
	}
	return list
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
	// sort the keyword list so highest ranking keyword is first
	sort.Sort(sort.Reverse(keyWordList))
	return keyWordList
}

// searches every found keyword and finds the best matching substring as per the keywords decomposition pattern
// returns a random response from the pool of responses for the specific decomposition pattern of the keyword
func findDecompPatterns(userInput string, keyWordList []KeyWord) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var response string
	// iterate over the all the keywords found in the user input string
	for _, keyWord := range keyWordList {
		// for every keyword iterate over all the decomposition patterns
		for _, decomp := range keyWord.Decomp {
			// fmt.Println(decomp.DisAssRule)
			// compile the decomposition pattern into a regular expression
			reg, err := regexp.Compile(decomp.DisAssRule)
			if err != nil {
				log.Fatal(err)
			}

			// check if the decomposition pattern is found in the user question and return the capture group values
			fragments := reg.FindStringSubmatch(userInput)

			// fmt.Println(len(fragments))
			fmt.Println(fragments)
			fmt.Println(userInput)
			// fmt.Println(keyWord.Decomp)
			if len(fragments) == 0 {
				continue
			}

			// choose a random response
			response = decomp.Responses[rand.Intn(len(decomp.Responses))]

			// disregard everything caught by the regex capture group if the response doesn't need it or nothing was captured
			if !strings.Contains(response, "%s") || len(fragments) == 1 {
				return response
			}

			// if the decomposition pattern is found in the user input string return a random response from the set of responses for the decomposition pattern
			if len(fragments) > 1 && strings.Contains(userInput, fragments[0]) {
				fmt.Println(reg)

				// reflect words like "my", "yours"
				fragment := substitute(fragments[1])

				// reassamble the response
				response = fmt.Sprintf(response, fragment)
				return response
			}
			// else if len(fragments) == 1 {
			// 	return response
			// }
		}
	}

	// if no keyword matches or no matching decomposition pattern is found fall back to a random response from the "xnone" keyword
	return elizaData[0].Decomp[0].Responses[rand.Intn(len(elizaData[0].Decomp[0].Responses))]

}

func substitute(fragment string) string {
	words := strings.Split(fragment, " ")
	for i, word := range words {
		if _, ok := substitutions[word]; ok {
			words[i] = substitutions[word]
		}
	}
	return strings.Join(words, " ")
}
