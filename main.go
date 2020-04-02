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
	_, err := http.Get(link)
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
