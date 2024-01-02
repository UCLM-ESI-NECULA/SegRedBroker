package client

import (
	"bytes"
	"io"
	"net/http"
	"seg-red-broker/internal/app/dao"
)

const (
	fileServiceBaseURL = "http://file-service" // Replace with actual File Service URL
)

// FileClient struct holds the HTTP client and the base URL for the File service
type FileClient struct {
	HttpClient *http.Client
	BaseURL    string
}

// NewFileClient creates a new instance of FileClient
func NewFileClient() *FileClient {
	return &FileClient{
		HttpClient: &http.Client{},
		BaseURL:    fileServiceBaseURL,
	}
}

// makeRequest is a helper function to make an HTTP request with the base URL
func (client *FileClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, client.BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return client.HttpClient.Do(req)
}

// GetFile requests a specific file from the File service
func (client *FileClient) GetFile(username, docID string) (*dao.FileContent, error) {
	resp, err := client.makeRequest(http.MethodGet, "/"+username+"/"+docID, nil)
	return getFileContent(resp, err)
}

// CreateFile sends a request to create a file in the File service
func (client *FileClient) CreateFile(username, docID string, content []byte) (*dao.FileSize, error) {
	resp, err := client.makeRequest(http.MethodPost, "/"+username+"/"+docID, bytes.NewBuffer(content))
	return getFileSize(resp, err)
}

// UpdateFile sends a request to update a file in the File service
func (client *FileClient) UpdateFile(username, docID string, content []byte) (*dao.FileSize, error) {
	resp, err := client.makeRequest(http.MethodPut, "/"+username+"/"+docID, bytes.NewBuffer(content))
	return getFileSize(resp, err)
}

// DeleteFile sends a request to delete a file in the File service
func (client *FileClient) DeleteFile(username, docID string) (bool, error) {
	resp, err := client.makeRequest(http.MethodDelete, "/"+username+"/"+docID, nil)
	return err == nil && resp.StatusCode == http.StatusOK, err
}

// GetAllUserDocs requests all documents for a specific user from the File service
func (client *FileClient) GetAllUserDocs(username string) (*map[string]string, error) {
	resp, err := client.makeRequest(http.MethodGet, "/"+username+"/_all_docs", nil)
	return getFiles(resp, err)
}

func getFileContent(resp *http.Response, err error) (*dao.FileContent, error) {
	if err != nil {
		return nil, err
	}

	var content dao.FileContent
	err = getBody(resp, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func getFileSize(resp *http.Response, err error) (*dao.FileSize, error) {
	if err != nil {
		return nil, err
	}

	var content dao.FileSize
	err = getBody(resp, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func getFiles(resp *http.Response, err error) (*map[string]string, error) {
	if err != nil {
		return nil, err
	}

	var content map[string]string
	err = getBody(resp, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}
