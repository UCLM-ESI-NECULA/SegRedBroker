package controller

import (
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
	// Check if the token is valid
	user, err := CheckTokenInput(c, fc.as)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Check username and docID
	username, docID, apiErr := checkParams(c)
	if apiErr != nil {
		common.ForwardError(c, apiErr)
		return
	}

	// Check if the token matches the file owner
	if username != user.Username {
		common.ForwardError(c, common.FileOwnerMismatch())
		return
	}

	// Get the file from the file service
	content, err := fc.fs.GetFile(user.Username, docID)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Return the file content
	c.JSON(http.StatusOK, content)

}
func (fc *FileControllerImpl) CreateFile(c *gin.Context) {
	// Check if the token is valid
	user, err := CheckTokenInput(c, fc.as)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Check username and docID
	username, docID, apiErr := checkParams(c)
	if apiErr != nil {
		common.ForwardError(c, apiErr)
		return
	}

	// Check if the token matches the file owner
	if username != user.Username {
		common.ForwardError(c, common.FileCreationMismatch())
		return
	}

	//ReadBody
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.ForwardError(c, common.BadRequestError("invalid request body"))
		return
	}

	size, err := fc.fs.CreateFile(user.Username, docID, requestBody)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, size)
}

func (fc *FileControllerImpl) UpdateFile(c *gin.Context) {
	// Check if the token is valid
	user, err := CheckTokenInput(c, fc.as)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Check username and docID
	username, docID, apiErr := checkParams(c)
	if apiErr != nil {
		common.ForwardError(c, apiErr)
		return
	}

	// Check if the token matches the file owner
	if username != user.Username {
		common.ForwardError(c, common.FileOwnerMismatch())
		return
	}

	//ReadBody
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.ForwardError(c, common.BadRequestError("invalid request body"))
		return
	}

	// Update the file in the file service
	size, err := fc.fs.UpdateFile(user.Username, docID, requestBody)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Return the file size
	c.JSON(http.StatusOK, size)
}

func (fc *FileControllerImpl) DeleteFile(c *gin.Context) {
	// Check if the token is valid
	user, err := CheckTokenInput(c, fc.as)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Check username and docID
	username, docID, apiErr := checkParams(c)
	if apiErr != nil {
		common.ForwardError(c, apiErr)
		return
	}

	// Check if the token matches the file owner
	if username != user.Username {
		common.ForwardError(c, common.FileOwnerMismatch())
		return
	}

	// Delete the file from the file service
	_, err = fc.fs.DeleteFile(user.Username, docID)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Return OK
	c.JSON(http.StatusOK, gin.H{})
}

func (fc *FileControllerImpl) GetAllUserDocs(c *gin.Context) {
	// Check if the token is valid
	user, err := CheckTokenInput(c, fc.as)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Check username
	username := c.Param("username")
	if username == "" {
		common.ForwardError(c, common.EmptyParamsError("username"))
		return
	}

	// Check if the token matches the file owner
	if username != user.Username {
		common.ForwardError(c, common.FileOwnerMismatch())
		return
	}

	docs, err := fc.fs.GetAllUserDocs(username)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

// checkParams checks if the username and docID are valid
func checkParams(c *gin.Context) (string, string, *common.APIError) {
	username := c.Param("username")
	if username == "" {
		return "", "", common.EmptyParamsError("username")
	}
	docID := c.Param("doc_id")
	if docID == "" {
		return "", "", common.EmptyParamsError("doc_id")
	}
	return username, docID, nil
}
