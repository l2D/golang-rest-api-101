package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"golang-rest-api-101/internal/routes"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

type ServerEnvironments struct {
	GoEnv   string
	AppPort string
	DSN     string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLmode  string
}

func LoadServerConfig() ServerEnvironments {
	// Load environment variables from .env file if not use Doppler.
	// err := godotenv.Load()

	// if err != nil {
	// 	panic("Error loading .env file")
	// }

	env := ServerEnvironments{
		GoEnv:   os.Getenv("GO_ENV"),
		AppPort: os.Getenv("PORT"),
		DSN:     LoadDBConfig(),
	}

	// Check struct fields for empty values
	if env.GoEnv == "" {
		env.GoEnv = "development"
	}

	if env.AppPort == "" {
		env.AppPort = "3000"
	}

	if env.DSN == "" {
		panic("DSN cannot be empty")
	}

	return env
}

func LoadDBConfig() string {
	config := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLmode:  os.Getenv("DB_SSL_MODE"),
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLmode)
}

func (server *Server) InitializeDB(dsn string) {
	var err error

	// dsn := LoadDBConfig()

	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Cannot connect to database\n")
		log.Fatal("Error:", err)
	} else {
		fmt.Printf("Connected to the database\n")
	}

	// server.DB.Debug().AutoMigrate()
}

func (server *Server) InitializeServer(config ServerEnvironments) {
	// Initialize Gin server
	/*
		r := gin.New()

		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	*/
	server.Router = gin.Default() // If use this, you don't need to inject gin.Logger() and gin.Recovery()

	if config.GoEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		server.Router.TrustedPlatform = gin.PlatformCloudflare
		fmt.Println("Production mode")
	} else {
		gin.SetMode(gin.DebugMode)
		err := server.Router.SetTrustedProxies([]string{"127.0.0.1"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Development mode")
	}

	// Initialize routes
	server.Router = routes.InitializeRoutes(server.Router, server.DB)

	// server.Router.Run(fmt.Sprintf(":%s", config.AppPort))

	// Graceful shutdown section. If you wantn't to use this, comment this section and uncomment above line.
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.AppPort),
		Handler: server.Router,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		fmt.Println("")
		log.Println("Shutting down server...")
		log.Println("STEP 1 | Shutdown Server <START>")

		// Received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Println("STEP 1 <ERROR> | Error HTTP server Shutdown: ", err)
		}

		log.Println("STEP 1 | Server Shutdown <SUCCESS>")

		log.Println("STEP 2 | Closing DB connection ... <START>")

		if db, err := server.DB.DB(); err == nil {
			if err := db.Close(); err != nil {
				log.Println("STEP 2 <ERROR> | Error closing DB connection: ", err)
			} else {
				log.Println("STEP 2 | DB connection closed <SUCCESS>")
			}
		} else {
			log.Println("STEP 2 <ERROR> | Error starting close DB connection: ", err)
		}

		close(idleConnsClosed)
	}()

	log.Println("Server should running on port: ", config.AppPort)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("Error starting or closing listener.")
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed

	log.Println("Server gracefully shutdown.")

	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf(":%s", config.AppPort),
	// 	Handler: server.Router,
	// }

	// go func() {
	// 	// service connections
	// 	if err := srv.ListenAndServe(); err != nil {
	// 		log.Printf("%s\n", err)
	// 	}
	// }()

	// fmt.Println("Server is running on port " + config.AppPort)

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)
	// <-quit

	// log.Println("STEP 1 | Shutdown Server ... <START>")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := srv.Shutdown(ctx); err == nil {
	// 	log.Println("STEP 1 | Server Shutdown <SUCCESS>")
	// } else {
	// 	log.Println("STEP 1 <ERROR> | Error closing server: ", err)
	// }

	// log.Println("STEP 2 | Closing DB connection ... <START>")

	// if db, err := server.DB.DB(); err == nil {
	// 	if err := db.Close(); err != nil {
	// 		log.Println("STEP 2 <ERROR> | Error closing DB connection: ", err)
	// 	} else {
	// 		log.Println("STEP 2 | DB connection closed <SUCCESS>")
	// 	}
	// } else {
	// 	log.Println("STEP 2 <ERROR> | Error starting close DB connection: ", err)
	// }

	// log.Println("STEP 3 | Server exiting <COMPLETE>")
}
