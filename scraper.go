package product_scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
	"regexp"
	"strconv"
)

// Scraper contains the httpClient and logger needed to scrape urls
type Scraper struct {
	httpClient IHttpClient
	logger     log.Logger
}

// NewScraper is a helper method to construct new Scraper with
// added logging context.
func NewScraper(client IHttpClient, logger log.Logger) Scraper {
	// all logs from this struct will have the key val "struct", "Scraper"
	// included for identification
	contextLogger := log.NewContext(logger).With("struct", "Scraper")
	return Scraper{httpClient: client, logger: contextLogger}
}

// Scrape take slice of urls and returns a Result and error.
// This process scrapes the seed urls for product urls, then scrapes products
// found from these product urls.
func (scraper Scraper) Scrape(urls []string) (Results, error) {
	// all logs from this function will have function name included
	logger := log.NewContext(scraper.logger).With("func", "scrape")
	logger.Log("level", "trace", "message", "started scraping", "urls", len(urls))
	// array and map for storing products scraped and urls visited
	foundProducts := []Product{}
	foundProductUrls := make(map[string]bool)

	// create data channels
	chUrls := make(chan string)
	chProducts := make(chan Product)
	// create control channels
	chScrapeUrl := make(chan error)
	chScrapeProduct := make(chan error)

	// for each of url provided, scrape for product links
	for _, url := range urls {
		go scraper.scrapeLinks(url, chUrls, chScrapeUrl)
	}

	// wait for every url scrape to finish
	for c := 0; c < len(urls); {
		select {
		// we have a product url, add to map. map
		// ensures we dont ever scrape the same url twice
		case url := <-chUrls:
			{
				logger.Log("level", "trace", "message", "Found product URL", "url", url)
				foundProductUrls[url] = true
			}
		// url scraper returned
		case err := <-chScrapeUrl:
			{
				if err != nil {
					logger.Log("level", "info", "message", "Error occured while scraping for product urls.", "error", err.Error())
				}
				c++
			}
		}
	}

	// for each product url we have found, scrape product details
	for url, _ := range foundProductUrls {
		go scraper.scrapeProduct(url, chProducts, chScrapeProduct)
	}

	// wait for every product scraper to finish
	for c := 0; c < len(foundProductUrls); {
		select {
		// we have scraped a product
		case product := <-chProducts:
			{
				// add product to list
				logger.Log("level", "trace", "message", "Scraped product info", "title", product.Title)
				foundProducts = append(foundProducts, product)
			}
		// product scraper returned
		case err := <-chScrapeProduct:
			{
				if err != nil {
					logger.Log("level", "info", "message", "Error occured while scraping for product info.", "error", err.Error())
				}
				c++
			}
		}
	}

	logger.Log("level", "trace", "message", "Finished scraping.", "products", len(foundProducts))

	// create new result object from scraped products and return
	return NewResult(foundProducts), nil
}

// scrapeProduct scrapes the product data from the url, sending the product
// back on the Product channel.
func (scraper Scraper) scrapeProduct(url string, chProduct chan Product, chReturn chan error) {
	var err error
	defer func() {
		// ensure we return a value after this function
		chReturn <- err
	}()

	// all logs from this function will have function name included
	logger := log.NewContext(scraper.logger).With("func", "scrapeProduct", "url", url)

	// get http response for url
	resp, err := scraper.httpClient.Get(url)
	if err != nil {
		// error getting url. Cant continue
		logger.Log("level", "info", "message", "Error during http get", "error", err.Error())
		return
	}

	// parse the body as html so we can search the DOM
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// error getting parsing html body. Cant continue
		logger.Log("level", "info", "message", "Error parsing response body as HTML", "error", err.Error())
		return
	}

	// there should be one div with class productTitleDescriptionContainer
	// contained within should be a h1, the text of the h1 is the title
	title := doc.Find(".productTitleDescriptionContainer").First().Find("h1").Text()
	// the first paragraph of the product text div contains the description
	desc := doc.Find(".productText").First().Find("p").First().Text()

	// the unit price is in the pricePerUnit paragraph
	priceString := doc.Find(".pricePerUnit").Text()
	// use reg expression to remove leading £ and any trailing text
	re := regexp.MustCompile("\\d*.\\d\\d")
	price := re.FindString(priceString)
	// convert this clean unit price string to float64
	floatPrice, _ := strconv.ParseFloat(price, 64)
	// convert from pounds to pence
	floatPence := floatPrice * 100

	// get the size of the product page from the response header
	size := resp.ContentLength

	// store unit price as int64 in pence for ease of addition
	product := Product{
		Title:       title,
		Description: desc,
		UnitPrice:   int64(floatPence),
		Size:        size,
	}

	// push product onto product channel
	chProduct <- product
}

// scrapeLinks scrapes product urls from the seed url, sending found
// product urls bacl on the chUrl channel.
func (scraper Scraper) scrapeLinks(url string, chUrl chan string, chReturn chan error) {
	var err error
	defer func() {
		// ensure we return a value after this function
		chReturn <- err
	}()

	// all logs from this function will have function name included
	logger := log.NewContext(scraper.logger).With("func", "scrapeLinks", "url", url)

	// get http response for url
	resp, err := scraper.httpClient.Get(url)
	if err != nil {
		// error getting url. Cant continue
		logger.Log("level", "info", "message", "Error during http get", "error", err.Error())
		return
	}

	// parse the body as html so we can search the DOM
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// error getting parsing html body. Cant continue
		logger.Log("level", "info", "message", "Error parsing response body as HTML", "error", err.Error())
		return
	}

	// The product urls are in div with class ".productInfo"
	doc.Find(".productInfo").Each(func(i int, s *goquery.Selection) {
		// for each product info div, check its children.
		// we are looking for h3 tag but we will check all children
		s.Children().Each(func(i int, s1 *goquery.Selection) {
			// now check inside for a href to product
			s1.Children().Each(func(i int, s2 *goquery.Selection) {
				url, ok := s2.Attr("href")
				if ok {
					// found a url! Push it onto the channel
					chUrl <- url
				}
			})
		})
	})
}
