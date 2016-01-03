package main
import (
	"github.com/timblackwell/product-scraper"
	"os"
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/go-kit/kit/log"
)

func main() {
	// get urls to look for products on
	seedUrls := os.Args[1:]

	// set up logger.
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.NewContext(logger).With("timestamp", log.DefaultTimestampUTC)
	// get new scraper that uses default http client
	scraper := product_scraper.NewScraper(http.DefaultClient, logger)
	// scrape urls
	results, err := scraper.Scrape(seedUrls)
	if err != nil {
		fmt.Printf("Error when scraping. %s", err)
		return
	}

	// Got the products, print json
	jbytes, _ := json.MarshalIndent(results, "", "    ")
	fmt.Println(string(jbytes))
}