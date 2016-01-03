package product_scraper

import (
	"encoding/json"
)

// store the collection of products and total of
// unit price in pence
type Results struct {
	Results []Product
	Total   int64
}

// helper method to calculate total for products
func NewResult(products []Product) Results {
	// go though every product, keeping running total
	// of unit price
	var runningTotal int64
	for _, product := range products {
		runningTotal += product.UnitPrice
	}

	// create result with running total
	results := Results{
		Results: products,
		Total:   runningTotal,
	}

	return results
}

// custom JSON marshaller to display total price in pounds
// rather than pence
func (r Results) MarshalJSON() ([]byte, error) {
	// convert int64 to float 64 and devide to get pence
	total := float64(r.Total)
	total = total / 100

	// return field value map of result with float for total
	return json.Marshal(map[string]interface{}{
		"results": r.Results,
		"total":   total,
	})
}
