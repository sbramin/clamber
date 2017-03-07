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

func crawler(url string) []string {
	semaphore <- struct{}{} // Initialize with empty struct

	page, err := extract(url)
	if err != nil {
		log.Print(err)
	}
	boltDown(page)
	atomic.AddUint64(&pageCount, 1)
	<-semaphore
	return page.Children
}

func goCrawl(url string) {
	links := make(chan []string)
	seen := make(map[string]bool)

	startURL := []string{url}

	go func() { links <- startURL }()

	for n := 1; n > 0; n-- {
		list := <-links
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					links <- crawler(link)
				}(link)
			}
		}
	}
}

// extract - Main link magic
func extract(URL string) (pageType, error) {

	var pT pageType

	resp, err := http.Get(URL)

	if err != nil {
		return pT, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return pT, fmt.Errorf("getting %s: %s", URL, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()

	if err != nil {
		return pT, fmt.Errorf("parsing %s as HTML: %v", URL, err)
	}

	pT.URL = URL

	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				if link.String() == URL || link.String() == baseURL {
					continue
				}
				if strings.Contains(link.String(), "#") {
					continue
				}
				if strings.HasPrefix(link.String(), baseURL) {
					pT.Children = append(pT.Children, link.String())
				}
			}
		} else if n.Type == html.ElementNode {
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
					if strings.HasPrefix(asset.String(), baseURL) {
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
					if strings.HasPrefix(asset.String(), baseURL) {
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
					if strings.HasPrefix(asset.String(), baseURL) {
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
						if strings.HasPrefix(asset.String(), baseURL) {
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
