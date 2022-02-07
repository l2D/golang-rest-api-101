package main

import (
	"fmt"
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"golang-rest-api-101/internal/routes"
)

func main() {

	// Load environment variables
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	goEnv := os.Getenv("GO_ENV")

	// Connect to database
	dsn := os.Getenv("DB_CONNECTION")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	// Initialize Gin server
	/*
		r := gin.New()

		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	*/
	r := gin.Default() // If use this, you don't need to inject gin.Logger() and gin.Recovery()

	if goEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		r.TrustedPlatform = gin.PlatformCloudflare
		fmt.Println("Production mode")
	} else {
		gin.SetMode(gin.DebugMode)
		err := r.SetTrustedProxies([]string{"127.0.0.1"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Development mode")
	}

	// Initialize routes
	r = routes.InitRoutes(r, db)

	// Run the server on port from .env file
	appPort := os.Getenv("PORT")
	log.Fatal(r.Run(":" + appPort))
}
