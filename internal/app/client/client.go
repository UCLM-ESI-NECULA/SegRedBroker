package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"seg-red-broker/internal/app/common"
)

type Client struct {
	HttpClient *http.Client
	BaseURL    string
}

// makeRequest is a helper function to make an HTTP request with the base URL
func (client Client) makeRequest(method, endpoint string, body *bytes.Buffer) (*http.Response, error) {
	req, _ := http.NewRequest(method, client.BaseURL+endpoint, body)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.HttpClient.Do(req)
	if resp == nil {
		err = fmt.Errorf("no response from the server")
	}
	if resp.StatusCode >= 400 {
		err = &common.APIError{StatusCode: resp.StatusCode, Message: resp.Status, Err: fmt.Errorf("error")}
	}
	return resp, err
}

// getBody reads the response body and unmarshal JSON into the provided target interface
func getBody(resp *http.Response, target interface{}) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	// Unmarshal JSON to the target type
	unmarshallErr := json.Unmarshal(body, &target)
	if unmarshallErr != nil {
		return unmarshallErr
	}

	return nil
}
