# clamber 

Web Clamberer with a bolt on.

## Install

    go get github.com/sbramin/clamber
	cd $GOPATH/src/github.com/sbramin/clamber
	go install

## Usage

### Help

    clamber -h

### Crawl
Normal:   

	clamber -job crawl -url https://sbramin.com
Skip SSL verification

    clamber -job crawl -url https://sbramin.com -nossl

### Review
Normal:

    clamber -job review -url https://sbramin.com
Pretty:

    clamber -job review -url https://sbramin.com -p

## TODO

* Write some more tests
