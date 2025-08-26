package helper

import (
	"term-service/logger"

	"github.com/gin-gonic/gin"
)

const (
	ErrInvalidOperation = "ERR_INVALID_OPERATION"
	ErrInvalidRequest   = "ERR_INVALID_REQUEST"
	ErrNotFount         = "ERR_NOT_FOUND"
	ErrInternal         = "ERR_INTERNAL"
)

type APIResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	ErrorCode  string      `json:"error_code,omitempty"`
}

func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func SendError(c *gin.Context, statusCode int, err error, errorCode string) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	} else {
		errMsg = errorCode
	}

	// Ghi log lỗi
	logger.WriteLogEx("error", errMsg, map[string]interface{}{
		"status_code": statusCode,
		"error_code":  errorCode,
		"path":        c.Request.URL.Path,
		"method":      c.Request.Method,
	})

	c.JSON(statusCode, APIResponse{
		StatusCode: statusCode,
		Error:      errMsg,
		ErrorCode:  errorCode,
	})
}
