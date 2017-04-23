package main

// page - a website URL, its children websites and any static assets on the page.
type page struct {
	URL      string   `json:"url"`
	Children []string `json:"children"`
	Assets   []asset  `json:"assets"`
}

// asset Sub type of page container url and type
type asset struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}
