package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	urls := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	c := make(chan string)
	// start to make the requests
	for _, url := range urls {
		go checkURL(url, c)
	}
	// wait for the responses and keep making requests
	// similar to: for { go checkURL(<-c, c) }
	for url := range c {
		go func(url string) {
			time.Sleep(time.Second * 2)
			checkURL(url, c)
		}(url)
	}
}

func checkURL(url string, channel chan string) {
	_, err := http.Get(url)
	if err != nil {
		fmt.Println(url, "might be down!")
		channel <- url
		return
	}
	fmt.Println(url, "is up!")
	channel <- url
}
