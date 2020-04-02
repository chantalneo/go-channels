package main

import (
	"fmt"
	"net/http"
)

func main() {
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	for _, link := range links {
		checkLink(link)
	}
}

func checkLink(link string) {
	_, err := http.Get(link) // Blocking call. When this runs, Main Go routine can do nothing else!
	if err != nil {
		fmt.Println(link, "might be down!")
		return
	}

	fmt.Println(link, "is up!")
}

// Notes:
// 1. How our code is being executed right now:
//    With our slice of links -> Take first link from slice -> Make request -> GET http://google.com -> Wait for a response, log it
//							  -> Take next link, make request -> GET http://facebook.com -> Wait for a response, log it
//							  -> Repeat for others
//
//    Basically it's in series, every single time we're making a request, We sit around and wait for the response to come back before making the next
//    So our aim to to make requests in parallel
//
// 2. When we launch a Go program, i.e. when we compile it and execute it, we automatically create one Go routine. You can think of a routine as being
//    something that exists inside of our running program or our process. This go routine takes every line of code inside of our program and executes them one by one.
//    Actual compiled form of our code might look a little bit differently than what we have.
//
// 3. Syntax of a Go routine:
//    go checkLink(link)
//    go - creates a new thread go routine
//    checkLink(link) - the function the newly created thread would run
