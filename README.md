# clamber 

Web Clamberer with a bolt on.

## Install

    go get github.com/sbramin/clamber
	go get github.com/boltdb/bolt
	cd $GOPATH/src/github.com/sbramin/clamber
	go install

## Usage

### Help

    clamber -h

### Crawl
Normal:   

	clamber -job crawl -url https://sbramin.com

### Review
Normal:

    clamber -job review -url https://sbramin.com
Pretty:

    clamber -job review -url https://sbramin.com -p

## TODO

* Write some more tests
