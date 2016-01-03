# product-scraper
Scrapes product details from given url

## To run:

1.  First clone repo or `go get "github.com/timblackwell/product-scraper"`
2.  cd into project folder - `cd $GOPATH/src/github.com/timblackwell/product-scraper`
3.  Run console project `go run console/main.go` with any urls to be scraped as args

example: `go run console/main.go "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html" `

## To test:

1.  First clone repo or `go get "github.com/timblackwell/product-scraper" `
2.  cd into project folder - `cd $GOPATH/src/github.com/timblackwell/product-scraper`
3.  Run `go test` 

example: `go test -coverprofile=coverage.out`