package main

import (
	"git.soufrabi.com/nextcode/rcee-isolate/internal/jobs"
	"git.soufrabi.com/nextcode/rcee-isolate/internal/web"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

func setupLogHandler(goEnv string) {

	var logHandler *slog.JSONHandler
	if goEnv == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
		gin.SetMode(gin.ReleaseMode)
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		gin.SetMode(gin.DebugMode)
	}
	slog.SetDefault(slog.New(logHandler))
}

func main() {
	var err error
	var goEnv string = os.Getenv("GO_ENV")

	setupLogHandler(goEnv)

	err = jobs.InitializeIsolate()
	if err != nil {
		os.Exit(1)
	}
	web.SetupServer()
}
