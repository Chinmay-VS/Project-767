# Project-767( Concurrent Synthetic Network Monitor)



A concurrent network checking tool built in Go. It sends HTTP requests to a list of URLs to track their status, uptime, and response speed in milliseconds. 

## How It Works

1. **Fixed Worker Pool:** Instead of launching a brand-new goroutine for every single website (which overloads memory and gets your IP banned), the script creates a small, fixed pool of 3 workers.
2. **Task Queuing:** The main function sends URLs into a buffered `jobs` channel. The 3 workers read from this queue simultaneously, process the network checks, and pass the data into a `resultsChan`.
3. **Deadlock Prevention:** The `wg.Wait()` and `close(resultsChan)` sequences run inside their own background goroutine. This stops the main thread from blocking or hanging on the unbuffered channel.
4. **Basic Security Checks:** Before making an actual HTTP request, the script filters out `localhost` and `127.0.0.1` addresses to protect against internal network exploits, and strictly requires the `https://` prefix.

## Tech Stack & Go Features Used
* Go (Golang)
* Concurrency tracking using sync.WaitGroup
* Thread-safe data communication via Channels
* Basic network request handling with net/http
