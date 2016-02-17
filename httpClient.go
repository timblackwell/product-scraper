package product_scraper

import (
	"net/http"
)

// IHttpClient interface used to mock the http client for testing.
type IHttpClient interface {
	// Only need to get web pages
	Get(url string) (*http.Response, error)
}
