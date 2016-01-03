package product_scraper_test

import (
	"encoding/json"
	"github.com/timblackwell/product-scraper"
	"io/ioutil"
	"testing"
)

func TestResultMarshalJSON(t *testing.T) {
	products := []product_scraper.Product{
		product_scraper.Product{
			Title:       "Sainsbury's Avocado, Ripe & Ready x2",
			Description: "Avocados",
			Size:        44479,
			UnitPrice:   180,
		},
		product_scraper.Product{
			Title:       "Sainsbury's Avocados, Ripe & Ready x4",
			Description: "Avocados",
			Size:        39610,
			UnitPrice:   320,
		},
		product_scraper.Product{
			Title:       "Sainsbury's Avocado Ripe & Ready XL Loose 300g",
			Description: "Avocados",
			Size:        39597,
			UnitPrice:   150,
		},
	}
	result := product_scraper.NewResult(products)
	jbytes, _ := json.MarshalIndent(result, "", "    ")

	expected, err := ioutil.ReadFile("testdata/result.json") // For read access.
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(jbytes) != string(expected) {
		t.Errorf("Did not get expected json")
		t.Logf("Expected:\n%s\n", string(expected))
		t.Logf("Got:\n%s\n", string(jbytes))
	}
}
