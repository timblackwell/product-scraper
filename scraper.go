package product_scraper

import (
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

type Scraper struct {
	httpClient IHttpClient
}

func NewScraper(client IHttpClient) Scraper {
	return Scraper{httpClient: client}
}

func (scraper Scraper) Scrape(urls []string) (Results, error) {
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
				foundProductUrls[url] = true
			}
		// url scraper returned
		case err := <-chScrapeUrl:
			{
				if err != nil {
					//todo log error
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
				foundProducts = append(foundProducts, product)
			}
		// product scraper returned
		case err := <-chScrapeProduct:
			{
				if err != nil {
					//todo log error
				}
				c++
			}
		}
	}

	// create new result object from scraped products and return
	return NewResult(foundProducts), nil
}

func (scraper Scraper) scrapeProduct(url string, chProduct chan Product, chReturn chan error) {
	var err error
	defer func() {
		// ensure we return a value after this function
		chReturn <- err
	}()

	// get http response for url
	resp, err := scraper.httpClient.Get(url)
	if err != nil {
		// error getting url. Cant continue
		return
	}

	// parse the body as html so we can search the DOM
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// error getting parsing html body. Cant continue
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

func (scraper Scraper) scrapeLinks(url string, chUrl chan string, chReturn chan error) {
	var err error
	defer func() {
		// ensure we return a value after this function
		chReturn <- err
	}()

	// get http response for url
	resp, err := scraper.httpClient.Get(url)
	if err != nil {
		// error getting url. Cant continue
		return
	}

	// parse the body as html so we can search the DOM
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// error getting parsing html body. Cant continue
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
