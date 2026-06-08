package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rossigee/mock-netbox/internal/middleware"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.StructuredLogging())
	return router
}

func TestHealthHandler(t *testing.T) {
	router := setupTestRouter()
	health := NewHealthHandler()
	router.GET("/health", health.Health)
	router.GET("/ready", health.Ready)

	tests := []struct {
		name   string
		path   string
		status int
	}{
		{"health check", "/health", http.StatusOK},
		{"readiness check", "/ready", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, w.Code)
			}
		})
	}
}

func TestDeviceHandler(t *testing.T) {
	router := setupTestRouter()
	devices := NewDeviceHandler()

	router.GET("/api/dcim/devices", devices.List)
	router.POST("/api/dcim/devices", devices.Create)
	router.GET("/api/dcim/devices/:id", devices.Get)
	router.PUT("/api/dcim/devices/:id", devices.Update)
	router.DELETE("/api/dcim/devices/:id", devices.Delete)

	t.Run("create device", func(t *testing.T) {
		body := map[string]interface{}{
			"name": "test-device",
			"site": 1,
			"type": "server",
		}
		data, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/dcim/devices", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("list devices", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/dcim/devices", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("get device", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/dcim/devices/1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("delete device", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/dcim/devices/1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}

func TestSiteHandler(t *testing.T) {
	router := setupTestRouter()
	sites := NewSiteHandler()

	router.GET("/api/sites", sites.List)
	router.POST("/api/sites", sites.Create)
	router.GET("/api/sites/:id", sites.Get)
	router.DELETE("/api/sites/:id", sites.Delete)

	t.Run("create site", func(t *testing.T) {
		body := map[string]interface{}{
			"name": "test-site",
		}
		data, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/sites", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("list sites", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/sites", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}
