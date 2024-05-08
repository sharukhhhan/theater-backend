package main

import "theater/internal/app"

var configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
