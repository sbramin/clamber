package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"golang.org/x/net/html"
)

var semaphore = make(chan struct{}, 32)

// pageType A website URL, its children websites and any static assets on the page.
type pageType struct {
	URL      string      `json:"url"`
	Children []string    `json:"children"`
	Assets   []assetType `json:"assets"`
}

// assetType Sub type of pageType container url and type
type assetType struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// crawler is the worker that runs extract and saves the page to boltDB, as well as passing
// child links back to goCrawl to spawn more crawlers.
func crawler(db *boltDB, baseURL *string, url string) []string {
	semaphore <- struct{}{} // Initialize with empty struct

	page, err := extract(baseURL, url)
	if err != nil {
		log.Print(err)
	}
	db.Write(baseURL, page)
	atomic.AddUint64(&pageCount, 1)
	<-semaphore
	return page.Children
}

// goCrawl is the master that takes links from workers and spawns off more workers, whilst
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

// extract does the main page parsing and applies rules like sticking to the parent domain
// and classifying assets.
func extract(baseURL *string, URL string) (pageType, error) {

	var pT pageType

	resp, err := http.Get(URL)

	if err != nil {
		return pT, err
	}
	if resp.StatusCode != http.StatusOK {
		err := resp.Body.Close()
		if err != nil {
			log.Print(err)
		}
		return pT, fmt.Errorf("getting %s: %s", URL, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Print(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Print(err)
	}

	if err != nil {
		return pT, fmt.Errorf("parsing %s as HTML: %v", URL, err)
	}

	pT.URL = URL

	visitNode := func(n *html.Node) {
		switch {
		case n.Type == html.ElementNode && n.Data == "a":
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)

				switch {
				case err != nil:
					continue
				case link.String() == *baseURL:
					continue
				case link.String() == URL:
					continue
				case strings.Contains(link.String(), "#"):
					continue
				case strings.HasPrefix(link.String(), *baseURL):
					pT.Children = append(pT.Children, link.String())
				}
			}
		case n.Type == html.ElementNode:
			switch n.Data {
			case "img":
				for _, i := range n.Attr {
					if i.Key != "src" {
						continue
					}
					asset, err := resp.Request.URL.Parse(i.Val)
					if err != nil {
						continue
					}
					if strings.HasPrefix(asset.String(), *baseURL) {
						pT.Assets = append(pT.Assets, assetType{asset.String(), "img"})
					}
				}
			case "script":
				for _, s := range n.Attr {
					if s.Key != "src" {
						continue
					}
					asset, err := resp.Request.URL.Parse(s.Val)
					if err != nil {
						continue
					}
					if strings.HasPrefix(asset.String(), *baseURL) {
						pT.Assets = append(pT.Assets, assetType{asset.String(), "script"})
					}
				}
			case "object":
				for _, s := range n.Attr {
					if s.Key != "data" {
						continue
					}
					asset, err := resp.Request.URL.Parse(s.Val)
					if err != nil {
						continue
					}
					if strings.HasPrefix(asset.String(), *baseURL) {
						pT.Assets = append(pT.Assets, assetType{asset.String(), "obj"})
					}
				}
			case "link":
				for _, a := range n.Attr {
					if a.Key == "href" {
						asset, err := resp.Request.URL.Parse(a.Val)
						if err != nil {
							continue
						}
						if strings.HasPrefix(asset.String(), *baseURL) {
							pT.Assets = append(pT.Assets, assetType{asset.String(), "css"})
						}
					}
				}
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return pT, nil
}

// forEachNode traverses the nodes of the HTML doc.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
