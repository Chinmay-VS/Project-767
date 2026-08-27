package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
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
		result.Error = "Internal URLs restricted"
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

// 1. THE WORKER FUNCTION
func worker(id int, jobs <-chan string, resultsChan chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// Worker keeps pulling jobs until the 'jobs' channel is closed and drained
	for url := range jobs {
		res := checkWebsite(url)
		resultsChan <- res
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

	numWorkers := 3

	// 2. CHANNELS & WAITGROUP CREATION
	jobs := make(chan string, len(urls))        // Buffered job queue
	resultsChan := make(chan CheckResult)       // Unbuffered results stream
	var wg sync.WaitGroup

	// 3. SPAWN FIXED WORKER POOL
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, resultsChan, &wg)
	}

	// 4. FEED JOBS & CLOSE JOBS QUEUE
	for _, url := range urls {
		jobs <- url
	}
	close(jobs) // Signals to workers: "No more jobs are coming!"

	// 5. INLINE SUPERVISOR 
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 6. EXPLICIT RECEIVE LOOP 
	for {
		res, ok := <-resultsChan
		if !ok {
			fmt.Println("\nAll jobs completed. Worker pool shut down cleanly.")
			break
		}

		fmt.Printf("URL :- %s -> Status :-%s\n", res.URL, res.Status)
	}
}
