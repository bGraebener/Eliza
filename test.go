package main

import (
	// "regexp"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)


type response struct {
	Keyword   string   `json:"keyword"`
	Rank      int      `json:"rank"`
	Responses []string `json:"responses"`
}

func (r responseList) Len() int { return len(r) }
func (r responseList) Less(r1, r2 int) bool {
	return r[r1].Rank < r[r2].Rank
}
func (r responseList) Swap(r1, r2 int) {
	r[r1], r[r2] = r[r2], r[r1]
}

type responseList []response

func main() {

	var list2 = responseList{}

	raw, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(raw, &list2)
	// fmt.Printf("Keyword: %v Rank: %v Keyword: %v Rank: %v", list2[0].Keyword, list2[0].Rank, list2[1].Keyword, list2[1].Rank)

	sort.Sort(list2)
	// fmt.Printf("\nKeyword: %v Rank: %v Keyword: %v Rank: %v", list2[0].Keyword, list2[0].Rank, list2[1].Keyword, list2[1].Rank)
	sort.Sort(sort.Reverse(list2))
	// fmt.Printf("\nKeyword: %v Rank: %v Keyword: %v Rank: %v", list2[0].Keyword, list2[0].Rank, list2[1].Keyword, list2[1].Rank)

	for _, v := range list2 {
		fmt.Printf("Keyword: %s\nRank: %d\n", v.Keyword, v.Rank)
		fmt.Println("Responses:")
		for _, j := range v.Responses {
			fmt.Println(j)
		}
	}
	// fmt.Printf(list2[0].Responses[3], "your family")

	fmt.Print("Please enter a phrase: ")
	var userInput string
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	userInput = reader.Text()
	fmt.Println(userInput)

	// r,_ regexp.Compile("")

}

func decompose(original string) string{
	var response string

	// r,_ regexp.Compile("[myMy]")




}
