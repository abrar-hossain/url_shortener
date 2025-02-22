package routes

import (
	"net/http" // Import HTTP package for handling responses
	"time"     // Import time package for handling expiry time

	"github.com/abrar-mashuk/url_shortener/api/database" // Import database connection (Redis)
	"github.com/abrar-mashuk/url_shortener/api/models"   // Import models for request structure
	"github.com/gin-gonic/gin"                           // Import Gin framework for handling API requests
)

// EditURL updates the URL and expiry time for a given short ID
func EditURL(c *gin.Context) {
	// Extract the shortID parameter from the URL (e.g., /api/v1/edit/abc123 → shortID = "abc123")
	shortID := c.Param("shortID")

	// Declare a variable to hold the request data
	var body models.Request

	// Bind JSON request body to `body` struct
	if err := c.ShouldBind(&body); err != nil {
		// If the request body is invalid, return a 400 Bad Request error
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Can not parse JSON",
		})
		return // Stop execution if there's an error
	}

	// Create a Redis client to interact with the database
	r := database.CreateClient(0)

	// Ensure Redis connection is closed after function execution
	defer r.Close()

	// Check if the shortID exists in the database
	val, err := r.Get(database.Ctx, shortID).Result()

	// If the shortID does not exist, return an error message
	if err != nil || val == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ShortID doesn't exist here",
		})
		return // Stop execution
	}

	// Update the URL and expiry time for the given shortID in Redis
	err = r.Set(database.Ctx, shortID, body.URL, body.Expiry*3600*time.Second).Err()

	// If Redis update fails, return an internal server error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update the shortened content",
		})
		return // Stop execution
	}

	// Return a success message if everything is updated successfully
	c.JSON(http.StatusOK, gin.H{
		"message": "The content has been updated.",
	})
}
