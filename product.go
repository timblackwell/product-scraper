package product_scraper

import (
	"encoding/json"
	"github.com/dustin/go-humanize"
)

// Product stores the product details. Unit price is in pence.
// Size is in bytes
type Product struct {
	Title       string
	Description string
	UnitPrice   int64
	Size        int64
}

// MarshalJSON converts the product to JSON byte array
// and returns this with an error.
func (p Product) MarshalJSON() ([]byte, error) {
	// convert unit price to float and divide so
	// can be shown as pounds.
	price := float64(p.UnitPrice)
	price = price / 100

	if p.Size < 0 {
		p.Size = 0
	}
	size := uint64(p.Size)

	// return field value map of result with float
	// storing unit price in pounds and string
	// of human readable size
	return json.Marshal(map[string]interface{}{
		"title":       p.Title,
		"size":        humanize.Bytes(size),
		"unit_price":  price,
		"description": p.Description,
	})
}
