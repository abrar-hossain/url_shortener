package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abrar-mashuk/url_shortener/api/database"
	"github.com/abrar-mashuk/url_shortener/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
	}

	database.InitializeClient()

	router := gin.Default()

	setupRouters(router)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(router.Run(":" + port))
}

func setupRouters(router *gin.Engine) {
	router.POST("/api/v1", routes.ShortenURL)
	router.GET("/api/v1/:shortID", routes.GetByShortId)
	router.DELETE("/api/v1/:shortID", routes.DeleteURL)
	router.PUT("/api/v1/:shortID", routes.EditURL)
	router.POST("/api/v1/:addTag", routes.AddTag)
	router.GET("/api/hello", routes.Hello)
}
