package routes

import (
	"net/http" // Import HTTP package to handle requests

	"github.com/gin-gonic/gin" // Import Gin framework for handling API requests
)

func Hello(c *gin.Context) {
	// If shortID is found, return the original URL in a JSON response
	c.JSON(http.StatusOK, gin.H{"data": "hello"})
}
