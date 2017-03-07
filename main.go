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

	if job == "" || baseURL == "" {
		fmt.Println("You must specify the job type and a URL")
		fmt.Println("use clamber -h for more information")
		os.Exit(1)
	}

	if job == "review" || job == "crawl" {
		if baseURL[len(baseURL)-1:] != "/" {
			baseURL += "/"
		}

		db := boltOn(job)
		defer boltOff(db)

		if job == "crawl" {
			start := time.Now()
			goCrawl(baseURL)
			fmt.Printf(
				"Crawled %d pages from %s in %.2f seconds \n", atomic.LoadUint64(&pageCount), baseURL, time.Since(start).Seconds())

		} else if job == "review" {
			review()
		}
	} else {
		fmt.Println("Job can be crawl or review, not", job)
		os.Exit(1)
	}

}
