package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CheckResult struct {
	URL        string
	Error      string
	Status     string
	StatusCode int
	ResponseMS int64
	CheckedAt  time.Time
}

func checkWebsite(url string) CheckResult {

	var result CheckResult

	result.URL = url

	if strings.Contains(url, "localhost") || strings.Contains(url, "127.0.0.1") {

		result.Status = "BLOCKED"
		result.Error = "Internal URL'S restricted"
		return result
	}

	if !strings.HasPrefix(url, "https://") {

		result.Status = "REJECTED"
		result.Error = "HTTPS Required"
		return result
	}

	client := &http.Client{

		Timeout: 5 * time.Second,
	}

	result.CheckedAt = time.Now()

	resp, err := client.Get(url)

	duration := time.Since(result.CheckedAt)

	if err != nil {

		result.Error = err.Error()
		result.Status = "ERROR"
		return result
	}

	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	if resp.StatusCode == http.StatusOK {

		result.Status = "UP"
		result.ResponseMS = duration.Milliseconds()

	} else {

		result.Status = "DOWN"
		result.Error = "Website Content not accessible"

	}

	return result

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

		res := checkWebsite(url)
		fmt.Printf("%+v\n", res)

	}

}
