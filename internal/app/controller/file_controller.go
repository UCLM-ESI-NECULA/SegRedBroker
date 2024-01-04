package controller

import (
	"errors"
	"io"
	"net/http"
	"seg-red-broker/internal/app/client"
	"seg-red-broker/internal/app/common"
	"seg-red-broker/internal/app/service"

	"github.com/gin-gonic/gin"
)

type FileControllerImpl struct {
	fs service.FileService
	as service.AuthService
}

func NewFileController(r *gin.RouterGroup) *FileControllerImpl {
	c := &FileControllerImpl{
		fs: service.NewFileService(*client.NewFileClient()),
		as: service.NewAuthService(*client.NewAuthClient()),
	}
	c.RegisterRoutes(r)
	return c
}

type FileController interface {
	GetFile(c *gin.Context)
	CreateFile(c *gin.Context)
	UpdateFile(c *gin.Context)
	DeleteFile(c *gin.Context)
	GetAllUserDocs(c *gin.Context)
}

// RegisterRoutes registers the authentication routes
func (fc *FileControllerImpl) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/:username/:doc_id", fc.GetFile)
	router.POST("/:username/:doc_id", fc.CreateFile)
	router.PUT("/:username/:doc_id", fc.UpdateFile)
	router.DELETE("/:username/:doc_id", fc.DeleteFile)
	router.GET("/:username/_all_docs", fc.GetAllUserDocs)
}

func (fc *FileControllerImpl) GetFile(c *gin.Context) {
	username, docID := checkParams(c)
	_, err := fc.as.ValidateToken(checkTokenInput(c))
	if err != nil {
		var apiErr *common.APIError
		if errors.As(err, &apiErr) {
			common.NewAPIError(c, apiErr.StatusCode, apiErr.Err, apiErr.Message)
			return
		}
		common.NewAPIError(c, http.StatusInternalServerError, err, err.Error())
		return
	}
	content, err := fc.fs.GetFile(username, docID)
	if err != nil {
		var apiErr *common.APIError
		if errors.As(err, &apiErr) {
			common.NewAPIError(c, apiErr.StatusCode, apiErr.Err, apiErr.Message)
			return
		}
		common.NewAPIError(c, http.StatusInternalServerError, err, err.Error())
		return
	}
	c.JSON(http.StatusOK, content)
}

func (fc *FileControllerImpl) CreateFile(c *gin.Context) {
	username, docID := checkParams(c)
	_, err := fc.as.ValidateToken(checkTokenInput(c))
	if err != nil {
		var apiErr *common.APIError
		if errors.As(err, &apiErr) {
			common.NewAPIError(c, apiErr.StatusCode, apiErr.Err, apiErr.Message)
			return
		}
		common.NewAPIError(c, http.StatusInternalServerError, err, err.Error())
		return
	}
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, "invalid request body")
		return
	}

	size, err := fc.fs.CreateFile(username, docID, requestBody)
	if err != nil {
		var apiErr *common.APIError
		if errors.As(err, &apiErr) {
			common.NewAPIError(c, apiErr.StatusCode, apiErr.Err, apiErr.Message)
			return
		}
		common.NewAPIError(c, http.StatusInternalServerError, err, err.Error())
		return
	}
	c.JSON(http.StatusOK, size)
}

func (fc *FileControllerImpl) UpdateFile(c *gin.Context) {
	username, docID := checkParams(c)
	_, err := fc.as.ValidateToken(checkTokenInput(c))
	if err != nil {
		var apiErr *common.APIError
		if errors.As(err, &apiErr) {
			common.NewAPIError(c, apiErr.StatusCode, apiErr.Err, apiErr.Message)
			return
		}
		common.NewAPIError(c, http.StatusInternalServerError, err, err.Error())
		return
	}
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, "invalid request body")
		return
	}

	size, _ := fc.fs.UpdateFile(username, docID, requestBody)
	c.JSON(http.StatusOK, size)
}

func (fc *FileControllerImpl) DeleteFile(c *gin.Context) {
	username, docID := checkParams(c)
	_, err := fc.as.ValidateToken(checkTokenInput(c))
	if err != nil {
		var apiErr *common.APIError
		if errors.As(err, &apiErr) {
			common.NewAPIError(c, apiErr.StatusCode, apiErr.Err, apiErr.Message)
			return
		}
		common.NewAPIError(c, http.StatusInternalServerError, err, err.Error())
		return
	}
	_, err = fc.fs.DeleteFile(username, docID)
	if err != nil {
		common.NewAPIError(c, http.StatusNotFound, err, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (fc *FileControllerImpl) GetAllUserDocs(c *gin.Context) {
	username := c.Param("username")
	_, err := fc.as.ValidateToken(checkTokenInput(c))
	if err != nil {
		var apiErr *common.APIError
		if errors.As(err, &apiErr) {
			common.NewAPIError(c, apiErr.StatusCode, apiErr.Err, apiErr.Message)
			return
		}
		common.NewAPIError(c, http.StatusInternalServerError, err, err.Error())
		return
	}
	if username == "" {
		common.NewAPIError(c, http.StatusBadRequest, nil, "username cannot be empty")
		return
	}
	docs, _ := fc.fs.GetAllUserDocs(username)
	c.JSON(http.StatusOK, docs)
}

// checkParams checks if the username and docID are valid
func checkParams(c *gin.Context) (string, string) {
	username := c.Param("username")
	docID := c.Param("doc_id")
	if username == "" || docID == "" {
		common.NewAPIError(c, http.StatusBadRequest, nil, "invalid input parameters")
	}
	return username, docID
}
