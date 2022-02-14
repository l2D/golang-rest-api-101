package main

import (
	"golang-rest-api-101/internal/config"
)

var server = &config.Server{}

func main() {
	// Load environments
	env := config.LoadServerConfig()

	// Connect to database
	server.InitializeDB(env.DSN)

	// Initialize Gin server
	server.InitializeServer(env)
}
