# clamber 

Go go WebCrawler with a bolt on.

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

    clamber  -job review -url https://sbramin.com
Pretty:

    clamber -job review -url https://sbramin.com -p
