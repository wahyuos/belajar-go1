package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string, errors interface{}) {
	c.JSON(code, Response{
		Status:  false,
		Message: message,
		Errors:  errors,
	})
}