package main

import (
	"github.com/4udiwe/download-task-service/internal/app"
)

func main() {
	app := app.New("config/config.yaml")
	app.Start()
}
