package product_scraper_test

import (
	"encoding/json"
	"github.com/timblackwell/product-scraper"
	"io/ioutil"
	"testing"
)

// TestProductMarshalJSON tests the MarshalJSON method on the product
// class.
func TestProductMarshalJSON(t *testing.T) {
	product := product_scraper.Product{
		Title:       "Sainsbury's Avocado, Ripe & Ready x2",
		Description: "Avocados",
		Size:        44479,
		UnitPrice:   180,
	}

	jbytes, _ := json.MarshalIndent(product, "", "  ")

	expected, err := ioutil.ReadFile("testdata/product.json") // For read access.
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(jbytes) != string(expected) {
		t.Errorf("Did not get expected json")
		t.Logf("Expected:\n%s\n", string(expected))
		t.Logf("Got:\n%s\n", string(jbytes))
	}
}
