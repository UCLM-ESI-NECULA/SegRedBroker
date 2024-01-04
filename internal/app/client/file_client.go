package client

import (
	"bytes"
	"net/http"
	"os"
	"seg-red-broker/internal/app/dao"
)

// FileClient struct holds the HTTP client and the base URL for the File service
type FileClient struct {
	Client
}

// NewFileClient creates a new instance of FileClient
func NewFileClient() *FileClient {
	return &FileClient{
		Client{
			HttpClient: &http.Client{},
			BaseURL:    os.Getenv("FILE_SERVICE_BASE_URL"),
		},
	}
}

// GetFile requests a specific file from the File service
func (client *FileClient) GetFile(username, docID string) (*dao.FileContent, error) {
	resp, err := client.makeRequest(http.MethodGet, "/"+username+"/"+docID, nil)
	if err != nil {
		return nil, err
	}
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
	err = getBody(resp, &content, err)
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
	err = getBody(resp, &content, err)
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
	err = getBody(resp, &content, err)
	if err != nil {
		return nil, err
	}
	return &content, nil
}
