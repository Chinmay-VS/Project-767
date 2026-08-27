# Project-767( Concurrent Synthetic Network Monitor)


A production-ready, highly optimized concurrent network monitoring tool engineered in Go (Golang). The system utilizes a fixed-size worker pool architecture to validate the operational availability, response latency, and security configurations of multiple remote endpoints simultaneously.

## Architectural Deep-Dive & Patterns Used

### 1. Resource Optimization (Worker Pool Strategy)
Instead of spawning an unbounded number of goroutines—which can trigger RAM overhead and destination IP blacklisting—this engine instantiates a bounded **Worker Pool (`numWorkers = 3`)**. A buffered `jobs` channel behaves as a thread-safe task queue, distributing traffic targets evenly across static worker instances.

### 2. Deadlock Avoidance & Memory Safety
* **Non-Blocking Orchestration:** The completion tracker (`wg.Wait()`) and result lifecycle pipeline are decoupled inside an independent supervisor goroutine. This guarantees that the unbuffered `resultsChan` never causes a thread deadlock.
* **Safe Channel Closure:** Closing the `jobs` channel acts as an explicit broadcast to all worker loops to wind down execution and exit cleanly once the task queue drains completely.

### 3. Network Security & SSRF Mitigation
To prevent Server-Side Request Forgery (SSRF) exploits, the monitor parses and enforces explicit network perimeter guards before issuing HTTP client dispatches:
* Enforces structural string filters targeting `localhost` and loopback addresses (`127.0.0.1`).
* Restricts transport layer requests strictly to secure `https://` schemas.
* Implements a tight 5-second connection `http.Client` timeout block to avoid hung connections.

## 🛠️ Tech Stack & Implementation Details
* **Language:** Go (Golang)
* **Standard Library Primitives:** `sync.WaitGroup`, `net/http`, `time`, `strings`
* **Concurrency Mechanics:** Channel streaming (`<-chan`, `chan<-`), Explicit for-range channel drains, Multiplexed control loops.
