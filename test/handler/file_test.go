package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-google-cloud-storage/app/config"
	"go-google-cloud-storage/app/constant"
	"go-google-cloud-storage/app/handler"
	"go-google-cloud-storage/app/helper"
	"go-google-cloud-storage/app/middleware"
	"go-google-cloud-storage/app/service"
)

// setupTestRouter creates a Gin engine with all routes registered
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	server := gin.New()
	server.Use(middleware.RequestID())

	// Setup Viper — leave GCS values empty so we test validation behavior
	v := viper.New()
	v.Set("APP_PORT", 4001)

	validator := config.NewValidator()
	fileService := service.NewFileService(v)
	fileHandler := handler.NewFileHandler(fileService, validator)

	api := server.Group("/api/v1")
	api.GET("/list", fileHandler.GetAllFile)
	api.POST("/upload", fileHandler.UploadFile)
	api.GET("/download", fileHandler.DownloadFile)
	api.POST("/bucket/create", fileHandler.CreateBucket)
	api.DELETE("/delete", fileHandler.DeleteFile)
	api.GET("/presigned-url", fileHandler.PresignedURL)

	return server
}

func TestPresignedURL_InvalidRequest(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "missing path parameter",
			query: "",
		},
		{
			name:  "empty path value",
			query: "path=",
		},
		{
			name:  "valid params but no GCS config",
			query: "path=test.jpg&expires=15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/presigned-url?"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Accept 400 (validation error) or 422 (GCS config missing)
			if w.Code != http.StatusBadRequest && w.Code != http.StatusUnprocessableEntity {
				t.Errorf("Expected 400 or 422, got %d: %s", w.Code, w.Body.String())
			}

			var response helper.Response
			json.Unmarshal(w.Body.Bytes(), &response)
			if response.ApiID == "" {
				t.Error("Response missing api_id")
			}
		})
	}
}

func TestPresignedURL_InvalidExpires(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "zero expires",
			query: "path=test.jpg&expires=0",
		},
		{
			name:  "negative expires",
			query: "path=test.jpg&expires=-1",
		},
		{
			name:  "expires exceeds max (10080 minutes = 7 days)",
			query: "path=test.jpg&expires=10081",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/presigned-url?"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// These should fail validation → 400, or fail on GCS config → 422
			if w.Code != http.StatusBadRequest && w.Code != http.StatusUnprocessableEntity {
				t.Errorf("Expected 400 or 422, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestRoutes_ResponseFormat(t *testing.T) {
	router := setupTestRouter()

	endpoints := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/api/v1/list", ""},
		{http.MethodGet, "/api/v1/presigned-url?path=test.jpg&expires=15", ""},
		{http.MethodGet, "/api/v1/download", ""},
		{http.MethodDelete, "/api/v1/delete", `{"path":"test.jpg"}`},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+"_"+ep.path, func(t *testing.T) {
			var req *http.Request
			if ep.body != "" {
				req, _ = http.NewRequest(ep.method, ep.path, nil)
				req.Header.Set("Content-Type", "application/json")
				// For DELETE with body, we need to pass the body differently
				bodyReader := strings.NewReader(ep.body)
				req, _ = http.NewRequest(ep.method, ep.path, bodyReader)
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(ep.method, ep.path, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var response helper.Response
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Response is not valid JSON: %s", w.Body.String())
				return
			}

			if response.ApiID == "" {
				t.Errorf("Response missing api_id field for %s %s", ep.method, ep.path)
			}
			if response.Status != constant.ResponseStatusSuccess &&
				response.Status != constant.ResponseStatusFailed {
				t.Errorf("Invalid status %q for %s %s", response.Status, ep.method, ep.path)
			}
		})
	}
}
