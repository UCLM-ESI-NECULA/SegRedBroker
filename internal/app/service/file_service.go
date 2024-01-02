package service

import (
	"seg-red-broker/internal/app/client"
	"seg-red-broker/internal/app/dao"
)

type FileServiceImpl struct {
	fc client.FileClient
}

func NewFileService(fc client.FileClient) *FileServiceImpl {
	return &FileServiceImpl{fc}
}

type FileService interface {
	GetFile(username, docID string) (*dao.FileContent, error)
	CreateFile(username, docID string, content []byte) (*dao.FileSize, error)
	UpdateFile(username, docID string, content []byte) (*dao.FileSize, error)
	DeleteFile(username, docID string) (bool, error)
	GetAllUserDocs(username string) (*map[string]string, error)
}

func (fs *FileServiceImpl) GetFile(username, docID string) (*dao.FileContent, error) {
	return fs.fc.GetFile(username, docID)
}

func (fs *FileServiceImpl) CreateFile(username, docID string, content []byte) (*dao.FileSize, error) {
	return fs.fc.CreateFile(username, docID, content)
}

func (fs *FileServiceImpl) UpdateFile(username, docID string, content []byte) (*dao.FileSize, error) {
	return fs.fc.UpdateFile(username, docID, content)
}

func (fs *FileServiceImpl) DeleteFile(username, docID string) (bool, error) {
	return fs.fc.DeleteFile(username, docID)
}

func (fs *FileServiceImpl) GetAllUserDocs(username string) (*map[string]string, error) {
	return fs.fc.GetAllUserDocs(username)
}
