package product_scraper_test

import (
	"fmt"
	"github.com/fortx/api/vendor/github.com/go-kit/kit/log"
	"github.com/timblackwell/product-scraper"
	"net/http"
	"os"
	"testing"
)

// mock http client used to control the tests
type mockHttp struct {
	urlResource map[string]string
}

// implementing IHttpClient interface
func (mock mockHttp) Get(url string) (response *http.Response, err error) {
	// check to see if there is a local resource for the url
	resource := mock.urlResource[url]
	if resource == "" {
		// dont know how to handel this url
		err = fmt.Errorf("Url not found")
		return
	}

	// open the local resource for url
	resourceHtml, err := os.Open(resource) // For read access.
	if err != nil {
		// something went wrong while trying to open the file
		return
	}

	// get file info for content length
	info, err := resourceHtml.Stat()
	if err != nil {
		// something went wrong while trying to open the file
		return
	}

	// create response, giving file (implements io.ReadCloser)
	response = &http.Response{
		Body: resourceHtml,
		ContentLength: info.Size(),
	}

	// return
	return
}

// test with working example
func TestNormal(t *testing.T) {
	// create mapping between the urls and local resources we want to serve
	index := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
	apricot := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-apricot-ripe---ready-320g.html"
	avocadoXL := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocado-xl-pinkerton-loose-300g.html"
	avocadoX2 := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocado--ripe---ready-x2.html"
	avocadoX4 := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocados--ripe---ready-x4.html"
	pears := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-conference-pears--ripe---ready-x4-%28minimum%29.html"
	kiwiGolden := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-golden-kiwi--taste-the-difference-x4-685641-p-44.html"
	kiwi := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-kiwi-fruit--ripe---ready-x4.html"
	urlResource := make(map[string]string)
	urlResource[index] = "testdata/Sainsbury's Ripe & ready.html"
	urlResource[apricot] = "testdata/Sainsbury's Apricot Ripe & Ready x5.html"
	urlResource[avocadoXL] = "testdata/Sainsbury's Avocado Ripe & Ready XL Loose 300g.html"
	urlResource[avocadoX2] = "testdata/Sainsbury's Avocado, Ripe & Ready x2.html"
	urlResource[avocadoX4] = "testdata/Sainsbury's Avocados, Ripe & Ready x4.html"
	urlResource[pears] = "testdata/Sainsbury's Conference Pears, Ripe & Ready x4 (minimum).html"
	urlResource[kiwiGolden] = "testdata/Sainsbury's Golden Kiwi x4.html"
	urlResource[kiwi] = "testdata/Sainsbury's Kiwi Fruit, Ripe & Ready x4.html"
	// using this mapping, create mock http client
	mock := mockHttp{urlResource: urlResource}

	// create a new scraper using the mock http client
	// and a dummy logger (we just want to test, dont want logs)
	scraper := product_scraper.NewScraper(mock, log.NewNopLogger())
	// scrape using the index url
	results, err := scraper.Scrape([]string{index})
	if err != nil {
		// this is an unexpected error, test failed
		t.Error(err)
	}

	// the total of all unit prices should be 1510 pence
	if results.Total != 1510 {
		// test failed!
		t.Errorf("Expected total %d, got: %d", 151, results.Total)
	}

	// we should have 7 products returned
	if len(results.Results) != 7 {
		// test failed!
		t.Errorf("Expected results %d, got: %d", 7, len(results.Results))
	}

}

func TestNoSeed(t *testing.T) {
	// create mock with no mapping. we want the error to hit the scraper
	urlResource := make(map[string]string)
	mock := mockHttp{urlResource: urlResource}

	// create a new scraper using the mock http client
	// and a dummy logger (we just want to test, dont want logs)
	scraper := product_scraper.NewScraper(mock, log.NewNopLogger())
	// scrape random url
	results, err := scraper.Scrape([]string{"http://notfound.com"})
	if err != nil {
		// this is an unexpected error, test failed
		t.Error(err)
	}

	// the total of all unit prices should be 0 pence
	if results.Total != 0 {
		// test failed!
		t.Errorf("Expected total %d, got: %d", 0, results.Total)
	}

	// we should have no results
	if len(results.Results) != 0 {
		// test failed!
		t.Errorf("Expected results %d, got: %d", 0, len(results.Results))
	}
}

// this test makes sure that if an error occurs scraping one
// url, it doesn't effect another
func TestMissingProducts(t *testing.T) {
	// create a mapping with only a few urls
	index := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
	apricot := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocado--ripe---ready-x2.html"
	urlResource := make(map[string]string)
	urlResource[index] = "testdata/Sainsbury's Ripe & ready.html"
	urlResource[apricot] = "testdata/Sainsbury's Avocado, Ripe & Ready x2.html"
	mock := mockHttp{urlResource: urlResource}

	// create a new scraper using the mock http client
	// and a dummy logger (we just want to test, dont want logs)
	scraper := product_scraper.NewScraper(mock, log.NewNopLogger())
	// scrape index
	results, err := scraper.Scrape([]string{index})
	if err != nil {
		// this is an unexpected error, test failed
		t.Error(err)
	}

	// the total of all unit prices should be 180 pence
	if results.Total != 180 {
		// test failed
		t.Errorf("Expected total %d, got: %d", 180, results.Total)
	}

	// we should have one result
	if len(results.Results) != 1 {
		// test failed
		t.Errorf("Expected results %d, got: %d", 1, len(results.Results))
	}

	// get the only product returned
	product := results.Results[0]

	// create the product we expected
	expected := product_scraper.Product{
		Title:       "Sainsbury's Avocado, Ripe & Ready x2",
		Description: "Avocados",
		Size:        44479,
		UnitPrice:   180,	}

	if product.Title != expected.Title {
		// test failed
		t.Errorf("Expected Title %d, got: %d", expected.Title, product.Title)
	}
	if product.Description != expected.Description {
		// test failed
		t.Errorf("Expected Description %d, got: %d", expected.Description, product.Description)
	}
	if product.Size != expected.Size {
		// test failed
		t.Errorf("Expected Size %d, got: %d", expected.Size, product.Size)
	}
	if product.UnitPrice != expected.UnitPrice {
		// test failed
		t.Errorf("Expected UnitPrice %d, got: %d", expected.UnitPrice, product.UnitPrice)
	}
}
