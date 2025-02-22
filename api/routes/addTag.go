package routes

import (
	"encoding/json" // Import JSON package for data conversion
	"net/http"      // Import HTTP package for handling API responses

	"github.com/abrar-mashuk/url_shortener/api/database" // Import database package to connect with Redis
	"github.com/gin-gonic/gin"                           // Import Gin framework to handle API requests
)

// Define a struct to represent the request body for adding a tag
type Tagrequest struct {
	ShortID string `json:"shortID"` // Shortened URL ID
	Tag     string `json:"tag"`     // Tag to be added
}

// AddTag function adds a tag to an existing short URL
func AddTag(c *gin.Context) {
	var tagRequest Tagrequest

	// Bind the incoming JSON request to the tagRequest struct
	if err := c.ShouldBindJSON(&tagRequest); err != nil {
		// Return an error response if the JSON request is invalid
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request Body",
		})
		return // Stop execution
	}

	// Extract shortID and tag from the request
	shortId := tagRequest.ShortID
	tag := tagRequest.Tag

	// Create a Redis client to interact with the database
	r := database.CreateClient(0)
	defer r.Close() // Ensure Redis connection is closed after function execution

	// Retrieve existing data from Redis using the shortID
	val, err := r.Get(database.Ctx, shortId).Result()
	if err != nil {
		// If shortID does not exist in Redis, return an error response
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found for the given ShortID",
		})
		return // Stop execution
	}

	// Declare a map to store the retrieved data
	var data map[string]interface{}

	// Try to parse the retrieved Redis data as JSON
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		// If data is not in JSON format, create a new JSON structure
		data = make(map[string]interface{})
		data["data"] = val
	}

	// Initialize a slice to store existing tags
	var tags []string

	// Check if "tags" already exists in the data and is a slice of interface{}
	if existingTags, ok := data["tags"].([]interface{}); ok {
		// Loop through each existing tag
		for _, t := range existingTags {
			// Convert the tag to a string and add it to the tags slice
			if strTag, ok := t.(string); ok {
				tags = append(tags, strTag)
			}
		}
	}

	// Check for duplicate tags before adding the new tag
	for _, existingTag := range tags {
		if existingTag == tag {
			// If the tag already exists, return an error response
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tag Already Exists",
			})
			return // Stop execution
		}
	}

	// Add the new tag to the tags slice
	tags = append(tags, tag)

	// Update the "tags" field in the data map
	data["tags"] = tags

	// Convert the updated data map back to JSON format
	updatedData, err := json.Marshal(data)
	if err != nil {
		// If JSON conversion fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to marshal updated data",
		})
		return // Stop execution
	}

	// Store the updated data in Redis with the same shortID
	err = r.Set(database.Ctx, shortId, updatedData, 0).Err()
	if err != nil {
		// If Redis update fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update the database",
		})
		return // Stop execution
	}

	// Respond with the updated data containing the new tag
	c.JSON(http.StatusOK, data)
}
