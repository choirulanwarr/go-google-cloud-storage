package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-google-cloud-storage/app/middleware"
	"time"
)

func NewServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	server.Use(middleware.RequestID())
	return server
}
