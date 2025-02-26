package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/abrar-mashuk/url_shortener/api/database"
	"github.com/abrar-mashuk/url_shortener/api/models"
	"github.com/abrar-mashuk/url_shortener/api/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// ShortenURL handles shortening a long URL into a short ID
func ShortenURL(c *gin.Context) {

	var body models.Request
	log.Println("hello world")
	if err := c.ShouldBind(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Parse JSON"})
		return
	}

	// Create a Redis client for handling API rate limiting (database 1)
	// r2 := database.CreateClient(1)
	// defer r2.Close()

	// Get the remaining quota for the user's IP address
	//val, err := r2.Get(database.Ctx, c.ClientIP()).Result()
	val, err := database.Client.Get(database.Ctx, c.ClientIP()).Result()
	fmt.Println(val)

	if err == redis.Nil { // If no quota exists, set the default API quota
		_ = database.Client.Set(database.Ctx, c.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()

	} else { // If quota exists, check the remaining requests
		val, _ = database.Client.Get(database.Ctx, c.ClientIP()).Result()
		valInt, _ := strconv.Atoi(val) // Convert quota to an integer

		if valInt <= 0 { // If quota is 0 or negative, reject the request
			limit, _ := database.Client.TTL(database.Ctx, c.ClientIP()).Result() // Get time left before quota resets
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
			return
		}
	}

	// Validate if the input URL is in a valid format
	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	// Ensure that the URL does not belong to the same domain (to prevent abuse)
	if !utils.IsDifferentDomain(body.URL) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "you cant hack this system",
		})
		return
	}

	body.URL = utils.EnsureHTTPPrefix(body.URL)

	var id string

	// If the user does not provide a custom short URL, generate a random ID
	if body.CustomShort == "" {
		id = uuid.New().String()[:6] // Generate a random 6-character UUID
	} else {
		id = body.CustomShort // Use the user-provided custom short URL
	}

	// r := database.CreateClient(0)
	// defer r.Close()

	// Check if the custom short ID already exists in Redis
	val, _ = database.Client.Get(database.Ctx, id).Result()

	if val != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "URL Custom Short Already Exists",
		})
		return
	}

	// Set default expiry time (24 hours) if the user has not provided one
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// Store the short ID and original URL in Redis with an expiry time
	err = database.Client.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to connect to the Redis server",
		})
		return
	}

	// Construct the response object with rate limits and URL information
	resp := models.Response{
		Expiry:          body.Expiry,
		XRateLimitReset: 30, //in minutes
		XRateRemainig:   10,
		URL:             body.URL,
		CustomShort:     "",
	}

	database.Client.Decr(database.Ctx, c.ClientIP())

	val, _ = database.Client.Get(database.Ctx, c.ClientIP()).Result()
	resp.XRateRemainig, _ = strconv.Atoi(val)

	ttl, _ := database.Client.TTL(database.Ctx, c.ClientIP()).Result()

	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	c.JSON(http.StatusOK, resp)
}
