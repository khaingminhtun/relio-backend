package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OK sends a 200 OK response with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 Created response with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// Accepted sends a 202 Accepted response with data.
// Use this for async operations where work is queued but not yet complete.
func Accepted(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAccepted, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response.
// Use this for successful DELETE or actions with no return value.
func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, Response{
		Success: true,
		Data:    nil,
		Message: "No Content, soft delete successfully",
	})
}

// BadRequest sends a 400 Bad Request response.
//
// Use this for HTTP/input validation errors that are detected
// directly by the handler.
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	})
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	})
}
