package frain

import (
	"encoding/json"
	"frain-dev/frain-client-go/cache"
	"frain-dev/frain-client-go/types"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	BanksEndpointCacheKey = "frain_banks_uptime_status"
	ApiTokenKey           = "FRAIN_API_TOKEN"
)

type Frain struct {
	Cache cache.Cache

	options Options
}

type Options struct {
	HTTPClient HTTPClient

	APIToken string

	ExpiryTime time.Duration
}

func New(cache cache.Cache) *Frain {
	options := Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIToken:   retrieveTokenFromEnv(),
		ExpiryTime: 1 * time.Minute,
	}
	return NewWithOptions(cache, options)
}

func NewWithOptions(cache cache.Cache, options Options) *Frain {
	ensureCache(&cache)

	ensureHTTPClient(&options)
	ensureAPIToken(&options)
	ensureExpiryTime(&options)

	frain := &Frain{
		Cache:   cache,
		options: options,
	}

	return frain
}

func ensureCache(c *cache.Cache) {
	if c != nil {
		return
	}
	log.Fatalf("Bad Cache configured: %+v", c)
}

func ensureHTTPClient(o *Options) {
	if o != nil && o.HTTPClient != nil {
		return
	}
	log.Fatalf("Bad HTTPClient in options: %+v", o)
}

func ensureAPIToken(o *Options) {
	if o != nil && o.APIToken != "" {
		return
	}
	log.Fatalf("Bad APIToken in options: %+v", o)
}

func ensureExpiryTime(o *Options) {
	if o != nil && o.ExpiryTime != 0 {
		return
	}
	log.Fatalf("Bad ExpiryTime in options: %+v", o)
}

func retrieveTokenFromEnv() string {
	token := os.Getenv(ApiTokenKey)
	if token == "" || len(token) == 0 {
		log.Println("Unable to retrieve Frain SDK API token")
	}
	return token
}

func (f *Frain) GetBanks() []types.Component {
	expiryTime := f.options.ExpiryTime

	var response []types.Component

	banksFromCache, err := f.GetBanksFromCache()
	if err == nil && len(banksFromCache) > 0 {
		log.Println("WARNING: Fetching from cache")
		response = banksFromCache
	} else {
		log.Println("WARNING: Fetching from API due to", err)
		banksFromApi, err := f.GetBanksFromApi()
		if err == nil {
			f.SaveBanksToCache(banksFromApi, expiryTime)
			response = banksFromApi
		} else {
			log.Println("ERROR: Failed to fetch banks from API due to", err)
		}
	}

	return response
}

func (f *Frain) SaveBanksToCache(components []types.Component, expiryTime time.Duration) {
	if len(components) > 0 {
		bytes, err := json.Marshal(components)
		if err != nil {
			log.Println("ERROR: Failed to Marshal JSON: ", err)
			return
		}
		f.Cache.Set(BanksEndpointCacheKey, string(bytes), expiryTime)
	}
}

func (f *Frain) GetBanksFromCache() ([]types.Component, error) {
	var components []types.Component
	dataString := f.Cache.Get(BanksEndpointCacheKey)
	if dataString == "" || len(dataString) == 0 {
		return nil, &types.FrainException{Message: types.ErrorBanksNotFoundInCache}
	}

	err := json.Unmarshal([]byte(dataString), &components)
	if err != nil {
		log.Println("Error while unmarshalling the components bytes: ", err)
		return nil, err
	}

	return components, nil
}
