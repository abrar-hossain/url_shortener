package routes

import (
	"encoding/json"
	"net/http"

	"github.com/abrar-mashuk/url_shortener/api/database"
	"github.com/gin-gonic/gin"
)

// Define a struct to represent the request body for adding a tag
type Tagrequest struct {
	ShortID string `json:"shortID"`
	Tag     string `json:"tag"`
}

// AddTag function adds a tag to an existing short URL
func AddTag(c *gin.Context) {
	var tagRequest Tagrequest

	if err := c.ShouldBindJSON(&tagRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request Body",
		})
		return
	}

	// Extract shortID and tag from the request
	shortId := tagRequest.ShortID
	tag := tagRequest.Tag

	r := database.CreateClient(0)
	defer r.Close()
	val, err := r.Get(database.Ctx, shortId).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found for the given ShortID",
		})
		return
	}

	var data map[string]interface{}

	// Try to parse the retrieved Redis data as JSON
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		// If data is not in JSON format, create a new JSON structure
		data = make(map[string]interface{})
		data["data"] = val
	}

	var tags []string

	if existingTags, ok := data["tags"].([]interface{}); ok {
		for _, t := range existingTags {
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
			return
		}
	}

	tags = append(tags, tag)

	data["tags"] = tags

	updatedData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to marshal updated data",
		})
		return
	}

	err = r.Set(database.Ctx, shortId, updatedData, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update the database",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
