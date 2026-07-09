package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "ok",
		"data":    data,
	})
}

func ResponseError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  400,
		"message": message,
	})
}
