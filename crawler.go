package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

var semaphore = make(chan struct{}, 32)

// getPage - fetchs html body with http.Get and passes back the response.
func getPage(baseURL *string, URL string) (p page, err error) {
	resp, err := http.Get(URL)
	defer resp.Body.Close()
	if err != nil {
		return p, fmt.Errorf("getting %s: %s", URL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return p, fmt.Errorf("status for %s: %s", URL, resp.Status)
	}
	return extract(baseURL, URL, resp)
}

// crawler - is the worker that runs extract and saves the page to boltDB, as well as passing
// child links back to goCrawl to spawn more crawlers.
func crawler(db *boltDB, baseURL *string, url string) []string {
	semaphore <- struct{}{} // Initialize with empty struct

	p, err := getPage(baseURL, url)
	if err != nil {
		log.Print(err)
	}
	db.Writer(baseURL, p)
	atomic.AddUint64(&pageCount, 1)
	<-semaphore
	return p.Children
}

// goCrawl - is the master that takes links from workers and spawns off more workers, whilst
// limiting new crawler functions to the semaphore size.
func goCrawl(db *boltDB, baseURL *string) {
	links := make(chan []string)
	seen := make(map[string]bool)

	startURL := []string{*baseURL}

	go func() { links <- startURL }()

	for n := 1; n > 0; n-- {
		list := <-links
		for _, url := range list {
			if !seen[url] {
				seen[url] = true
				n++
				go func(url string) {
					links <- crawler(db, baseURL, url)
				}(url)
			}
		}
	}
}
