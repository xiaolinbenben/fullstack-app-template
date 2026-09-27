package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Success: true, Message: "", Data: data})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Envelope{Success: false, Message: message, Data: map[string]any{}})
}

func Write(c *gin.Context, status int, body Envelope) {
	c.JSON(status, body)
}
