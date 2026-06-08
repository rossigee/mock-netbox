// Mock NetBox service for E2E testing.
package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rossigee/mock-netbox/internal/handler"
	"github.com/rossigee/mock-netbox/internal/middleware"
)

func main() {
	// Configure logging
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	slog.SetDefault(slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
	))

	// Configure Gin
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release"
	}
	gin.SetMode(ginMode)

	// Create router
	router := gin.New()

	// Apply middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.StructuredLogging())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoints
	health := handler.NewHealthHandler()
	router.GET("/health", health.Health)
	router.GET("/ready", health.Ready)

	// API endpoints
	api := router.Group("/api/dcim")
	{
		devices := handler.NewDeviceHandler()
		api.GET("/devices", devices.List)
		api.POST("/devices", devices.Create)
		api.GET("/devices/:id", devices.Get)
		api.PUT("/devices/:id", devices.Update)
		api.DELETE("/devices/:id", devices.Delete)

		interfaces := handler.NewInterfaceHandler()
		api.GET("/interfaces", interfaces.List)
		api.POST("/interfaces", interfaces.Create)
		api.GET("/interfaces/:id", interfaces.Get)
		api.DELETE("/interfaces/:id", interfaces.Delete)
	}

	siteAPI := router.Group("/api/sites")
	{
		sites := handler.NewSiteHandler()
		siteAPI.GET("", sites.List)
		siteAPI.POST("", sites.Create)
		siteAPI.GET("/:id", sites.Get)
		siteAPI.DELETE("/:id", sites.Delete)
	}

	// Start server
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("starting mock-netbox", slog.String("port", port))
	if err := router.Run(":" + port); err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
