package main

import (
	"frain-dev/frain-client-go"
	"frain-dev/frain-client-go/cache"
	"log"
	"net/http"
	"time"
)

func main() {

	// with default options - retrieves API Token from ENV
	runDefault()

	// with options
	runWithOptions()
}

func runDefault() {
	mockCache := cache.NewMockCache()

	f := frain.New(mockCache)

	banks := f.GetBanks()

	log.Printf("\nGet Banks (default): %+v\n", banks)
}

func runWithOptions() {
	mockCache := cache.NewMockCache()

	f := frain.NewWithOptions(mockCache, frain.Options{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		APIToken:   "api_token",
		ExpiryTime: 10 * time.Minute,
	})

	banks := f.GetBanks()

	log.Printf("\nGet Banks (with options): %+v\n", banks)
}
