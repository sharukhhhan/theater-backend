package main

import "theater/internal/app"

var configPath = "config/config.yaml"

// @title Theater API
// @version 1.0
// @description This is a sample server for a theater app.
// @host localhost:8080
// @BasePath /api/v1
// @Security ApiKeyAuth
func main() {
	app.Run(configPath)
}
