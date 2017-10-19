package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "index.html")
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	question := r.Header.Get("user-question")

	question = strings.ToUpper(question)

	if len(question) > 0 {
		fmt.Fprintf(w, "%v", question)
	}
}

func main() {

	http.HandleFunc("/question", questionHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}
