package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
