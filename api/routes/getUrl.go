package routes

import (
	"net/http" // Import HTTP package to handle requests

	"github.com/abrar-mashuk/url_shortener/api/database" // Import database package to connect to Redis
	"github.com/gin-gonic/gin"                           // Import Gin framework for handling API requests
)

func GetByShortId(c *gin.Context) {
	// Extract the shortID parameter from the URL
	shortID := c.Param("shortID")

	// Create a Redis client (database connection)
	r := database.CreateClient(0)

	// Ensure Redis connection is closed after function execution
	defer r.Close()

	// Retrieve the original URL from Redis using the shortID as a key
	val, err := r.Get(database.Ctx, shortID).Result()

	// If Redis does not find the shortID, return a 404 error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found for given shortID",
		})
		return // Stop execution here
	}

	// If shortID is found, return the original URL in a JSON response
	c.JSON(http.StatusOK, gin.H{"data": val})
}
