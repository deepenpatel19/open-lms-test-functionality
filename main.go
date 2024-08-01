package main

import (
	"net/http"
	"os"

	// JWT
	jwt "github.com/appleboy/gin-jwt/v2"

	// Gin
	"github.com/gin-gonic/gin"

	// Zap logger
	"go.uber.org/zap"

	// Internal packages
	"github.com/open-lms-test-functionality/api"
	"github.com/open-lms-test-functionality/core"
	"github.com/open-lms-test-functionality/logger"
	"github.com/open-lms-test-functionality/middleware"
	"github.com/open-lms-test-functionality/models"
)

func Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ok"})
}

func main() {
	core.ReadEnvFile()        // Configure ENV File
	logger.LoggerInit()       // Configure Logger
	models.RunMigrations()    // Run migrations to sync db schema related changes
	models.CreateConnection() // Create DB connection pool

	authMiddleware, err := middleware.GetAuthMiddleware()
	if err != nil {
		logger.Logger.Error("MAIN :: Error while configuring auth middleware", zap.Error(err))
		os.Exit(1)
	}

	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.POST("/user", api.CreateUser) // Open endpoint

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		logger.Logger.Error("MAIN :: No route found", zap.Any("claims", claims))
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.POST("/token", authMiddleware.LoginHandler)
	r.GET("/logout", authMiddleware.LogoutHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	// Auth Group
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())

	// User APIs
	auth.PUT("/user/:userId", api.UpdateUser)
	auth.DELETE("/user/:userId", api.DeleteUser)

	// Test APIs
	auth.GET("/tests", api.FetchTests)
	auth.GET("/test/:testId", api.FetchTest)
	auth.POST("/test", api.CreateTest)
	auth.PUT("/test/:testId", api.UpdateTest)
	auth.DELETE("/test/:testId", api.DeleteTest)

	// Question APIs
	auth.GET("/questions", api.FetchQuestions)
	auth.GET("/question/:questionId", api.FetchQuestion)
	auth.POST("/question", api.CreateQuestion)
	auth.PUT("/question/:questionId", api.UpdateQuestion)
	auth.DELETE("/question/:questionId", api.DeleteQuestion)

	// Test questionary APIs
	auth.GET("test/:testId/questions", api.FetchTestQuestionary)
	auth.POST("test/:testId/generate_questionary", api.CreateTestQuestionary)
	auth.PUT("test/:testId/question/:questionId/add_question", api.AddTestQuestion)
	auth.DELETE("test/:testId/question/:questionId", api.DeleteTestQuestion)

	// Test question submission APIs
	auth.PUT("test/:testId/question/:questionId", api.SubmitTestQuestionSubmission)
	auth.GET("test/:testId/submissions", api.FetchTestQuestionSubmissions)

	// Starting server
	if err := r.Run(":8000"); err != nil {
		logger.Logger.Fatal("Failed to start the server:", zap.Error(err))
		return
	}
}
