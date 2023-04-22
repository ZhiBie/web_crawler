package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type CrawlController struct {
	seen map[string]bool
	mu sync.Mutex
}

func (cc *CrawlController) IsVisited(url string) bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if _, ok := cc.seen[url]; ok {
		// skip if have seen this url
		fmt.Println("skip and return")
		return true
	}
	cc.seen[url] = true
	return false
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func (cc *CrawlController) Crawl(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup) {
	// fmt.Println("Crawl:", url, "|| seen:", cc.seen)
	defer wg.Done()

	if depth <= 0 || cc.IsVisited(url) {
		return
	}
	
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go cc.Crawl(u, depth-1, fetcher, wg)
	}
	return
}

func main() {
	cc := CrawlController{seen: make(map[string]bool)}
	wg := sync.WaitGroup{}
	
	wg.Add(1)
	go cc.Crawl("https://golang.org/", 4, fetcher, &wg)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
