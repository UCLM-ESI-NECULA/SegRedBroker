package client

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"seg-red-broker/internal/app/common"
	"seg-red-broker/internal/app/dao"
)

// FileClient struct holds the HTTP client and the base URL for the File service
type FileClient struct {
	Client *resty.Client
}

// NewFileClient creates a new instance of FileClient
func NewFileClient() *FileClient {
	cl := resty.New()
	cl.
		SetBaseURL(os.Getenv("FILE_SERVICE_BASE_URL")).
		SetHeader("Accept", "application/json").
		SetError(&common.APIError{})
	return &FileClient{
		Client: cl,
	}
}

// GetFile requests a specific file from the File service
func (client *FileClient) GetFile(username, docID string) (*dao.FileContent, error) {
	resp, err := client.Client.R().
		SetResult(&dao.FileContent{}).
		SetPathParams(map[string]string{"username": username, "docID": docID}).
		Get("/{username}/{docID}")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, resp.Error().(*common.APIError)
	}
	return resp.Result().(*dao.FileContent), nil
}

// CreateFile sends a request to create a file in the File service
func (client *FileClient) CreateFile(username, docID string, content []byte) (*dao.FileSize, error) {
	resp, err := client.Client.R().
		SetResult(&dao.FileSize{}).
		SetPathParams(map[string]string{"username": username, "docID": docID}).
		SetBody(content).
		Post("/{username}/{docID}")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, resp.Error().(*common.APIError)
	}
	return resp.Result().(*dao.FileSize), nil
}

// UpdateFile sends a request to update a file in the File service
func (client *FileClient) UpdateFile(username, docID string, content []byte) (*dao.FileSize, error) {
	resp, err := client.Client.R().
		SetResult(&dao.FileSize{}).
		SetPathParams(map[string]string{"username": username, "docID": docID}).
		SetBody(content).
		Put("/{username}/{docID}")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, resp.Error().(*common.APIError)
	}
	return resp.Result().(*dao.FileSize), nil
}

// DeleteFile sends a request to delete a file in the File service
func (client *FileClient) DeleteFile(username, docID string) error {
	resp, err := client.Client.R().
		SetPathParams(map[string]string{"username": username, "docID": docID}).
		Delete("/{username}/{docID}")
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	} else {
		return resp.Error().(*common.APIError)
	}
}

// GetAllUserDocs requests all documents for a specific user from the File service
func (client *FileClient) GetAllUserDocs(username string) (*map[string]string, error) {
	m := make(map[string]string)
	resp, err := client.Client.R().
		SetResult(&m).
		SetPathParams(map[string]string{"username": username}).
		Get("/{username}/_all_docs")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, resp.Error().(*common.APIError)
	}
	return resp.Result().(*map[string]string), nil
}
