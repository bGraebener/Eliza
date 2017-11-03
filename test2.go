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

type KeyWord struct {
	Keyword string `json:"keyword" `
	Rank    int    `json:"rank"`
	Decomp  []Decomp
}

type Decomp struct {
	DisAssRule string   `json:"value"`
	Responses  []string `json:"responses"`
}

type keyWords []KeyWord

func (r keyWords) Len() int { return len(r) }
func (r keyWords) Less(r1, r2 int) bool {
	return r[r1].Rank < r[r2].Rank
}
func (r keyWords) Swap(r1, r2 int) {
	r[r1], r[r2] = r[r2], r[r1]
}

var elizaData keyWords

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	elizaData = readElizaData()

	userInput := getUserInput()

	keyWordList := getKeyWordList(userInput)

	findDecompPatterns(userInput, keyWordList)
}

func readElizaData() keyWords {
	var list keyWords

	raw, _ := ioutil.ReadFile("test.json")
	err := json.Unmarshal(raw, &list)
	// fmt.Println(string(raw))
	fmt.Println(err)
	// sort.Sort(sort.Reverse(list))

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

func getKeyWordList(userInput string) (keyWordList keyWords) {

	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}

	userInput = reg.ReplaceAllString(userInput, " ")

	userWords := strings.Split(userInput, " ")
	for _, word := range userWords {
		// fmt.Println(word)
		for i := range elizaData {
			if elizaData[i].Keyword == word {
				// keyWordList[list[i].Rank] = word
				keyWordList = append(keyWordList, elizaData[i])
				break
			}
		}
	}
	sort.Sort(sort.Reverse(keyWordList))
	return keyWordList
}

func findDecompPatterns(userInput string, keyWordList []KeyWord) {

	for _, keyWord := range keyWordList {
		for _, decomp := range keyWord.Decomp {
			fmt.Println(decomp.DisAssRule)
			reg, _ := regexp.Compile(decomp.DisAssRule)
			fmt.Println(reg)
			if reg.MatchString(userInput) {
				fmt.Println(decomp.Responses[rand.Intn(len(decomp.Responses))])
				return
			}

		}
	}

	fmt.Println(elizaData[0].Decomp[0].Responses[rand.Intn(len(elizaData[0].Decomp[0].Responses))])

}
