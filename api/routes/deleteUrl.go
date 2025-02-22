package routes

import (
	"net/http"

	"github.com/abrar-mashuk/url_shortener/api/database"
	"github.com/gin-gonic/gin"
)

func DeleteURL(c *gin.Context) {
	shortID := c.Param("shortID")

	r := database.CreateClient(0)

	defer r.Close()

	_, err := r.Del(database.Ctx, shortID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to Delete shortend Link.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shortend URL Deleted Successfully",
	})
}
