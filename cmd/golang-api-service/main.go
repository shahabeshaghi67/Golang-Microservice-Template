package main

import (
	"github.com/go-kit/log/level"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/app"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/config"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/logger"
)

func main() {
	logger := logger.NewLogger()
	cfg := config.Load()

	app := app.NewApp(cfg, logger, false)
	_, err := app.Wire()
	defer app.Unwire()
	if err != nil {
		_ = level.Error(logger).Log("msg", "app wiring failed", "error:", err)
		panic(err)
	}
	app.Run()
}
