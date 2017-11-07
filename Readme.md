# Eliza - Chatbot

## Introduction

This is an implementation of the classic Eliza Chatbot as described [here](https://en.wikipedia.org/wiki/ELIZA) and [here](https://www.masswerk.at/elizabot/).

It was created as a project for the module "Data Representation and Querying" in 3rd year of the Bsc. Software Development Course at Galway-Mayo Institute of Technology.

The Eliza-Chatbot is a chat bot implementation that analyses user input on a very simple level and responses with a phrase, mimiking a psychotherapy session.

## How the application works

After entering a user name, the user is brought to the "session"-page. The user name is stored in a HTTP-Header for later reference. The user is greeted with a randomly chosen greeting from an array of greetings. The user can then enter a phrase in the input textbox and either press enter or click on the arrow button to send the question.
Submitting the question triggers an AJAX-Request to be sent passing the user input.

A ranked list is created of all keywords that are contained in the users question.
For all keywords in the list, the keywords decomposition rules are compiled into regular expressions and checked against the original question. If a pattern is found, a random response is chosen and reassembled with any captured text from the regular expression.

If no keyword or matching pattern is found, a randomly chosen, generic response is returned.

The session can be ended by submitting "stop", "bye", "goodbye", "quit" or "exit" as a question.

## Usage

With Go installed and the GOPATH set, follow these steps to download and install the Eliza-Chatbot.

Download the project with

```go
go get github.com/bgraebener/Eliza
```

Navigate to

    %GOPATH%/src/github.com/bgraebener/Eliza

Build with:

```go
go build ./eliza.go
```

And execute with

    ./eliza

In any browser navigate to

    localhost:8080

and follow the instructions

## References

1. Basic ideas and data for Eliza responses: http://www.masswerk.at/elizabot/
2. Blurred Background: https://stackoverflow.com/questions/38366571/how-to-blur-the-background-image-only-in-css
3. Instructions on how to implement Eliza: http://www.chayden.net/eliza/instructions.txt

