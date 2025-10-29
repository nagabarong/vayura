package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON response helper functions
func JSONSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func JSONError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func JSONBadRequest(c *gin.Context, err error) {
	JSONError(c, http.StatusBadRequest, err)
}

func JSONUnauthorized(c *gin.Context, err error) {
	JSONError(c, http.StatusUnauthorized, err)
}

func JSONInternalServerError(c *gin.Context, err error) {
	JSONError(c, http.StatusInternalServerError, err)
}

func JSONNotFound(c *gin.Context, err error) {
	JSONError(c, http.StatusNotFound, err)
}
