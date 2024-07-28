package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-lms-test-functionality/core"
	"github.com/open-lms-test-functionality/logger"
	"github.com/open-lms-test-functionality/models"
	"github.com/open-lms-test-functionality/utils"
	"go.uber.org/zap"
)

func main() {
	core.ReadEnvFile()        // Configure ENV File
	logger.LoggerInit()       // Configure Logger
	models.RunMigrations()    // Run migrations to sync db schema related changes
	models.CreateConnection() // Create DB connection pool

	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		uuidString := utils.GetUUID()
		c.Header("X-REQUEST-ID", uuidString)
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	// Starting server
	if err := r.Run(":8000"); err != nil {
		logger.Logger.Fatal("Failed to start the server:", zap.Error(err))
		return
	}
}
