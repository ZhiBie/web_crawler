# Exercise from A Tour of Go

Question: https://go.dev/tour/concurrency/10

Notes:
- To run: `go run crawl.go`
- Used `go mod init github.com/ZhiBie/web_crawler/go` to add `go.mod` file, which is needed to import local test file.
- Go concurrency features used: 
  - `sync.WaitGroup`, without this the program directly returns, and goroutines are not run.
  - `sync.Mutex`
- Successful test result: 4 lines of "found..." are printed
