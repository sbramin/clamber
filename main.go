package main

import (
	"flag"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var (
	baseURL   string
	pageCount uint64
	pretty    bool
)

func main() {

	j := flag.String("job", "", "What job: crawl or review")
	u := flag.String("url", "", "A URL to start with: eg. https://sbramin.com")
	p := flag.Bool("p", false, "Pretty JSON")

	flag.Parse()
	job := *j
	baseURL = *u
	pretty = *p

	switch {
	case (job != "crawl" && job != "review") || baseURL == "":
		fmt.Println("You must specify the job type and a URL")
		fmt.Println("eg. clamber -url https://sbramin.com -job crawl")
		fmt.Println("use clamber -h for more information")
		os.Exit(1)
	case job == "crawl" || job == "review":
		if baseURL[len(baseURL)-1:] != "/" {
			baseURL += "/"
		}
		db := boltOn(job)
		defer boltOff(db)
	}

	switch job {
	case "crawl":
		start := time.Now()
		goCrawl(baseURL)
		fmt.Printf(
			"Crawled %d pages from %s in %.2f seconds \n", atomic.LoadUint64(&pageCount), baseURL, time.Since(start).Seconds())
	case "review":
		review()
	}
}
