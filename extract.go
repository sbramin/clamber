package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// extract does the main page parsing and applies rules like sticking to the parent domain
// and classifying assets.
func extract(baseURL *string, URL string) (page, error) {

	var p page

	resp, err := http.Get(URL)

	if err != nil {
		return p, err
	}
	if resp.StatusCode != http.StatusOK {
		err := resp.Body.Close()
		if err != nil {
			log.Print(err)
		}
		return p, fmt.Errorf("getting %s: %s", URL, resp.Status)
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
		return p, fmt.Errorf("parsing %s as HTML: %v", URL, err)
	}

	p.URL = URL

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
					p.Children = append(p.Children, link.String())
				}
			}
		case n.Type == html.ElementNode:
			switch n.Data {
			case "img":
				for _, i := range n.Attr {
					if i.Key != "src" {
						continue
					}
					a, err := resp.Request.URL.Parse(i.Val)
					if err != nil {
						continue
					}
					if strings.HasPrefix(a.String(), *baseURL) {
						p.Assets = append(p.Assets, asset{a.String(), "img"})
					}
				}
			case "script":
				for _, s := range n.Attr {
					if s.Key != "src" {
						continue
					}
					a, err := resp.Request.URL.Parse(s.Val)
					if err != nil {
						continue
					}
					if strings.HasPrefix(a.String(), *baseURL) {
						p.Assets = append(p.Assets, asset{a.String(), "script"})
					}
				}
			case "object":
				for _, s := range n.Attr {
					if s.Key != "data" {
						continue
					}
					a, err := resp.Request.URL.Parse(s.Val)
					if err != nil {
						continue
					}
					if strings.HasPrefix(a.String(), *baseURL) {
						p.Assets = append(p.Assets, asset{a.String(), "obj"})
					}
				}
			case "link":
				for _, a := range n.Attr {
					if a.Key == "href" {
						a, err := resp.Request.URL.Parse(a.Val)
						if err != nil {
							continue
						}
						if strings.HasPrefix(a.String(), *baseURL) {
							p.Assets = append(p.Assets, asset{a.String(), "css"})
						}
					}
				}
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return p, nil
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
