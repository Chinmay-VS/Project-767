package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CheckList struct {

}

func checkWebsite(url string) {

	if strings.Contains(url, "localhost") || strings.Contains(url, "127.0.0.1") {

		fmt.Printf("BLOCKED %s - INTERNAL URL RESTRICTED\n", url)
		return
	}

	if !strings.HasPrefix(url, "https://") {

		fmt.Printf("BLOCKED %s -  NOT SECURED by https\n", url)
		return
	}

	client := &http.Client{

		Timeout: 5 * time.Second,
	}

	start := time.Now()

	resp, err := client.Get(url)

	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Website %s not found\n", url)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		fmt.Printf("Status of URL %s : %d\n", url, resp.StatusCode)

	} else {

		fmt.Printf("WARN - %s - STATUS - %d (%dms)\n", url, resp.StatusCode, duration.Milliseconds())
		return
	}

}

func main() {

	urls := []string{
		"https://google.com",
		"https://facebook.com",
		"https://golang.org",
		"https://twitter.com",
		"https://this-is-a-fake-site-123456.com",
		"https://google.com/this-page-does-not-exist",
	}

	for _, url := range urls {

		checkWebsite(url)
	}

}
