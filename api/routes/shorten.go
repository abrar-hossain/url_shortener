package routes

import (
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
	// Declare a variable to store request data
	var body models.Request
	log.Println("hello world")
	// Bind JSON request body to the `body` struct
	if err := c.ShouldBind(&body); err != nil {
		// Return error if the request body is not valid JSON
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Parse JSON"})
		return
	}

	// Create a Redis client for handling API rate limiting (database 1)
	r2 := database.CreateClient(1)
	defer r2.Close() // Ensure Redis connection is closed after execution

	// Get the remaining quota for the user's IP address
	val, err := r2.Get(database.Ctx, c.ClientIP()).Result()

	if err == redis.Nil { // If no quota exists, set the default API quota
		_ = r2.Set(database.Ctx, c.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()

	} else { // If quota exists, check the remaining requests
		val, _ = r2.Get(database.Ctx, c.ClientIP()).Result()
		valInt, _ := strconv.Atoi(val) // Convert quota to an integer

		if valInt <= 0 { // If quota is 0 or negative, reject the request
			limit, _ := r2.TTL(database.Ctx, c.ClientIP()).Result() // Get time left before quota resets
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

	// Ensure the URL has HTTP or HTTPS prefix
	body.URL = utils.EnsureHTTPPrefix(body.URL)

	// Declare a variable to store the shortened ID
	var id string

	// If the user does not provide a custom short URL, generate a random ID
	if body.CustomShort == "" {
		id = uuid.New().String()[:6] // Generate a random 6-character UUID
	} else {
		id = body.CustomShort // Use the user-provided custom short URL
	}

	// Create a Redis client for storing the shortened URL (database 0)
	r := database.CreateClient(0)
	defer r.Close() // Ensure Redis connection is closed after execution

	// Check if the custom short ID already exists in Redis
	val, _ = r.Get(database.Ctx, id).Result()

	if val != "" { // If the short ID already exists, return an error
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
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil { // If storing the short URL fails, return an error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to connect to the Redis server",
		})
		return
	}

	// Construct the response object with rate limits and URL information
	resp := models.Response{
		Expiry:          body.Expiry,
		XRateLimitReset: 30, // Placeholder value (will be updated)
		XRateRemainig:   10, // Placeholder value (will be updated)
		URL:             body.URL,
		CustomShort:     "",
	}

	// Decrease the user's available quota by 1 request
	r2.Decr(database.Ctx, c.ClientIP())

	// Get the updated quota and the time left before reset
	val, _ = r2.Get(database.Ctx, c.ClientIP()).Result()
	resp.XRateRemainig, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.ClientIP()).Result()

	// Update the rate limit reset time in the response
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	// Construct the full short URL using the domain name from environment variables
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	// Return the response with the shortened URL
	c.JSON(http.StatusOK, resp)
}
