package response

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Response struct {
	Message string `json:"message"`
}

func NewResponse(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, Response{message})
}
