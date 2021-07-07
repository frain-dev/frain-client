package frain

import (
	"encoding/json"
	"fmt"
	"frain-dev/frain-client-go/types"
	"io/ioutil"
	"log"
	"net/http"
)

var BanksEndpoint = "https://api.frain.dev/api/v1/banks"

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (f *Frain) GetBanksFromApi() ([]types.Component, error) {
	httpClient := f.options.HTTPClient
	token := f.options.APIToken

	req, err := http.NewRequest("GET", BanksEndpoint, nil)
	if err != nil {
		errorMsg := fmt.Sprint("Error creating new request | ", err)
		return nil, &types.ServerException{Message: &errorMsg}
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		errorMsg := fmt.Sprint("Error processing request | ", err)
		return nil, &types.ServerException{Message: &errorMsg}
	}

	var components []types.Component

	err = parseAPIResponse(resp, &components)
	if err != nil {
		return nil, err
	}

	return components, nil
}

func parseAPIResponse(resp *http.Response, resultPtr interface{}) error {
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMsg := fmt.Sprint("Error while reading the response bytes | ", err)
		return &types.ServerException{Message: &errorMsg}
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("Error closing response body | ", err)
		}
	}()

	var response types.APIResponse

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		errorMsg := fmt.Sprint("Error while unmarshalling the response bytes | ", err)
		return &types.ServerException{Message: &errorMsg}
	}

	if !response.Status {
		return &types.ClientException{Message: &response.Message}
	}

	err = json.Unmarshal(*response.Data, resultPtr)
	if err != nil {
		errorMsg := fmt.Sprint("Error while unmarshalling the response data bytes | ", err)
		return &types.ServerException{Message: &errorMsg}
	}
	return nil
}
