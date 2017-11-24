# Eliza - Chatbot

## Author: Bastian Graebener

>Student-ID: G00340600
>
>Module: Data Representation and Querying
>
>Year: 3
>
>Program: Bsc. in Software Development
>
>Lecturer: Ian McLoughlin
>
>Institute: Galway-Mayo Institute of Technology

## Introduction

This is an implementation of the classic Eliza Chatbot as described [here](https://en.wikipedia.org/wiki/ELIZA) and [here](https://www.masswerk.at/elizabot/) in Golang. The application is a webapplication that can be accessed through a browser by navigating to the address "localhost:8080".

It was created as a project for the module "Data Representation and Querying" in 3rd year of the Bsc. Software Development Course at Galway-Mayo Institute of Technology.

You can find the [Project Specifications](https://github.com/bGraebener/Eliza.git) here.

The Eliza-Chatbot is a chat bot implementation that analyses user input on a very simple level and responses with a phrase, mimiking a psychotherapy session.

## How the application works

Upon the first request to port 8080 the user is asked to enter a name. After entering a name the name is stored in a HTTP-Header and the user is rerouted to the "session"-page.  The user is greeted with a randomly chosen greeting from an array of greetings. The user can then enter a phrase in the input textbox and either press enter or click on the arrow button to send the question.
Submitting the question triggers an AJAX-Request to be sent passing the user input string.

The input string is split into a slice of individual words. The slice is iterated over and a list is created of all keywords that are contained in the users question. The list is then sorted by rank in descending order, so the highest ranked keyword is first in the list.
For all keywords in the list, the keywords decomposition rules are compiled into regular expressions and checked against the original question. If a pattern is found, a random response is chosen. The captured text from the regular expression is searched for pronouns. Every found pronoun is then reflected by an appropriate counterpart.
Afterwards a response is assembled from the randomly chosen response and the reflected capure group.

If no keyword or matching pattern is found, a randomly chosen, generic response is returned.

The session can be ended by submitting "stop", "bye", "goodbye", "quit" or "exit" as a question.

## Usage

With Go and Git installed and the GOPATH set, follow these steps to download and install the Eliza-Chatbot.

Download the project with

```go
go get github.com/bGraebener/Eliza
```

or

```git
git clone https://github.com/bGraebener/Eliza.git
```

Navigate to

    %GOPATH%/src/github.com/bGraebener/Eliza

Build with:

```go
go build ./eliza.go
```

And execute with

    ./eliza

In any browser navigate to

    localhost:8080

and follow the instructions.

## References and Sources used to complete the project

1. Basic ideas and data for Eliza responses: http://www.masswerk.at/elizabot/
2. Instructions on how to implement Eliza: http://www.chayden.net/eliza/instructions.txt
3. List of Reflections for pronouns: https://www.smallsurething.com/implementing-the-famous-eliza-chatbot-in-python/
4. Reading json files in golang: https://www.chazzuka.com/2015/03/load-parse-json-file-golang/
5. Using regular expressions in golang:
https://github.com/StefanSchroeder/Golang-Regex-Tutorial
6. Expression language for golangs html templates:
http://2016.8-p.info/post/06-18-go-html-template/
7. Blurred Background: https://stackoverflow.com/questions/38366571/how-to-blur-the-background-image-only-in-css