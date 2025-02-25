package routes

import (
	"net/http"

	"github.com/abrar-mashuk/url_shortener/api/database"
	"github.com/gin-gonic/gin"
)

func GetByShortId(c *gin.Context) {

	shortID := c.Param("shortID")

	r := database.CreateClient(0)

	defer r.Close()

	// Retrieve the original URL from Redis using the shortID as a key
	val, err := r.Get(database.Ctx, shortID).Result()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found for given shortID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": val})
}
