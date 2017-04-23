package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
)

var (
	pageCount uint64
)

func main() {

	job := flag.String("job", "", "What job: crawl or review")
	baseURL := flag.String("url", "", "A URL to start with: eg. https://sbramin.com")
	pretty := flag.Bool("p", false, "Pretty JSON")

	flag.Parse()

	switch {
	case !(*job == "crawl" || *job == "review") || *baseURL == "":
		fmt.Println("You must specify the job type and a URL")
		fmt.Println("eg. clamber -url https://sbramin.com -job crawl")
		fmt.Println("use clamber -h for more information")
		os.Exit(1)
	case *job == "crawl" || *job == "review":
		url := *baseURL
		if url[len(url)-1:] != "/" {
			*baseURL += "/"
		}
	}

	db, err := setupDB(baseURL, job)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch *job {
	case "crawl":
		start := time.Now()
		goCrawl(db, baseURL)
		fmt.Printf(
			"Crawled %d pages from %s in %.2f seconds \n", atomic.LoadUint64(&pageCount), *baseURL, time.Since(start).Seconds())
	case "review":
		review(db, baseURL, pretty)
	}
}
