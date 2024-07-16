package web

import (
	"encoding/json"
	"git.soufrabi.com/nextcode/rcee-isolate/internal/api"
	"git.soufrabi.com/nextcode/rcee-isolate/internal/jobs"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

func handleRun(c *gin.Context) {
	var requestBody api.RunRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	slog.Debug("Run Request", "body", requestBody)

	res := jobs.RunCode(requestBody)

	c.JSON(http.StatusOK, res)

}

func SetupServer() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	r.POST("/run", handleRun)

	// listen and serve on :8080
	err := r.Run()
	if err != nil {
		slog.Error("failed to start gin server", "err", err)
	}
}
