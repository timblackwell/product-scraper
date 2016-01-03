package product_scraper_test

import (
	"fmt"
	"github.com/timblackwell/product-scraper"
	"net/http"
	"os"
	"testing"
)

type mockHttp struct {
	urlResource map[string]string
}

func (mock mockHttp) Get(url string) (response *http.Response, err error) {
	resource := mock.urlResource[url]
	if resource == "" {
		err = fmt.Errorf("Url not found")
		return
	}

	resourceHtml, err := os.Open(resource) // For read access.
	if err != nil {
		return
	}
	response = &http.Response{
		Body: resourceHtml,
	}

	return
}

func TestNormal(t *testing.T) {
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
	mock := mockHttp{urlResource: urlResource}

	scraper := product_scraper.NewScraper(mock)
	results, err := scraper.Scrape([]string{index})
	if err != nil {
		t.Error(err)
	}

	if results.Total != 1510 {
		t.Errorf("Expected total %d, got: %d", 151, results.Total)
	}

	if len(results.Results) != 7 {
		t.Errorf("Expected results %d, got: %d", 7, len(results.Results))
	}

}

func TestNoSeed(t *testing.T) {
	urlResource := make(map[string]string)
	mock := mockHttp{urlResource: urlResource}

	scraper := product_scraper.NewScraper(mock)
	results, err := scraper.Scrape([]string{"http://notfound.com"})
	if err != nil {
		t.Error(err)
	}

	if results.Total != 0 {
		t.Errorf("Expected total %d, got: %d", 0, results.Total)
	}

	if len(results.Results) != 0 {
		t.Errorf("Expected results %d, got: %d", 0, len(results.Results))
	}
}

func TestMissingProducts(t *testing.T) {
	index := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
	apricot := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-apricot-ripe---ready-320g.html"
	urlResource := make(map[string]string)
	urlResource[index] = "testdata/Sainsbury's Ripe & ready.html"
	urlResource[apricot] = "testdata/Sainsbury's Apricot Ripe & Ready x5.html"
	mock := mockHttp{urlResource: urlResource}

	scraper := product_scraper.NewScraper(mock)
	results, err := scraper.Scrape([]string{index})
	if err != nil {
		t.Error(err)
	}

	if results.Total != 350 {
		t.Errorf("Expected total %d, got: %d", 350, results.Total)
	}

	if len(results.Results) != 1 {
		t.Errorf("Expected results %d, got: %d", 1, len(results.Results))
	}
}
