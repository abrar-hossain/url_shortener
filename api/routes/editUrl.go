package routes

import (
	"net/http"
	"time"

	"github.com/abrar-mashuk/url_shortener/api/database"
	"github.com/abrar-mashuk/url_shortener/api/models"
	"github.com/gin-gonic/gin"
)

// EditURL updates the URL and expiry time for a given short ID
func EditURL(c *gin.Context) {

	shortID := c.Param("shortID")

	var body models.Request

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Can not parse JSON",
		})
		return
	}

	r := database.CreateClient(0)

	defer r.Close()

	val, err := r.Get(database.Ctx, shortID).Result()

	if err != nil || val == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ShortID doesn't exist here",
		})
		return
	}

	// Update the URL and expiry time for the given shortID in Redis
	err = r.Set(database.Ctx, shortID, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update the shortened content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The content has been updated.",
	})
}
